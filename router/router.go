package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"sctek.com/typhoon/th-platform-gateway/controller"
)

type HttpService struct {

}

func (h*HttpService)HttpRouter(r *gin.Engine) {
	new(controller.ShortMessageService).Router(r)

}


func (h*HttpService)initService(ctx *context.Context){

	//无法控制结束
	//r:=gin.Default()
	//h.HttpRouter(r)
	//r.Run(common.Config.Listen)

}


