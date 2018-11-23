package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/router"
	"sctek.com/typhoon/th-platform-gateway/service"
	"time"
)

func main() {
	common.CheckErr(common.LoadConfig())
	// common.CheckErr(common.OpenRedis())
	common.CheckErr(common.OpenDb())
	//common.CheckErr(common.SetupLogger())
	defer common.DB.Close()
	common.CheckErr(
		common.InitLogger(
			common.Config.Log.LogFileDir,
			common.Config.Log.LogFileName,
			""))
	//启动定时器
	cron := cron.New()
	cron.Start()
	defer cron.Stop()
	if len(common.Config.PollingTime) == 0 {
		panic("配置文件定时器的的轮询时间读取失败！！！")
	}
	//60秒执行一次
	cron.AddFunc(common.Config.PollingTime, service.CronSelect)

	//mq 初始化
	common.CheckErr(service.Init())
	defer service.Fini()
	common.CheckErr(service.Receive())

	//工作池初始化
	service.InitPool()
	defer service.ClosePool()

	r := gin.New()
	//r.Use(middleware.Logger(), gin.Recovery())
	// 路由
	new(router.HttpService).HttpRouter(r)
	srv := &http.Server{
		Addr:    common.Config.Listen,
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
	fmt.Println("listen:", common.Config.Listen)
	//Wait for interrupt signal to gracefully shutdown the server with
	//a timeout of 30 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
	//stop http listen
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalln("Server Shutdown:", err)
	}
	log.Println("Server exiting")
	//stop http listen

}
