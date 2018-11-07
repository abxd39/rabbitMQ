package main

import (
	"github.com/gin-gonic/gin"
	"sctek.com/typhoon/th-platform-gateway/controller"
)

func httpRouter(r *gin.Engine) {
	new(controller.ShortMessageService).Router(r)

}
