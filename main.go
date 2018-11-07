package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sctek.com/typhoon/th-platform-gateway/common"
	"sctek.com/typhoon/th-platform-gateway/rabbitMQ"
	"sctek.com/typhoon/th-platform-gateway/router"
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
	fmt.Println("mq status")
	log.Println("Uri=", common.Config.Mq.Uri)
	log.Println("uName=", common.Config.Mq.QueueName)
	rabbitMQ.NewConsumer(common.Config.Mq.Uri, common.Config.Mq.Exchange, common.Config.Mq.ExchangeType, common.Config.Mq.QueueName, common.Config.Mq.Key, common.Config.Mq.ConsumerTag)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 30 seconds.
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
