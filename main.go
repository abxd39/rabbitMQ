package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"log"
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
	common.CheckErr(common.SetupLogger())
	defer common.DB.Close()

	r := gin.New()
	//r.Use(middleware.Logger(), gin.Recovery())
	// 路由
	router.HttpRouter(r)
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
	//启动定时器
	cron := cron.New()
	cron.Start()
	defer cron.Stop()
	//5秒执行一次
	cron.AddFunc("*/60 * * * * *", service.CronSelect)
	//mq 初始化
	common.CheckErr(service.Init())
	//工作池初始化
	service.InitPool()
	defer service.ClosePool()
	defer service.Fini()
	common.CheckErr(service.Receive())
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

}
