package router

import (
	"github.com/gin-gonic/gin"
	"sctek.com/typhoon/th-platform-gateway/controller"
)

func HttpRouter(r *gin.Engine) {
	new(controller.ShortMessageService).Router(r)

}
