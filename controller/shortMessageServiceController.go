package controller

import (
	"github.com/gin-gonic/gin"
)

type ShortMessageService struct {
	Controller
}

func (s *ShortMessageService) Router(r *gin.Engine)  {
	sms := r.Group("/sms")
	{
		sms.GET("/hello", s.Hello)
	}

}

func (s*ShortMessageService)Hello(c*gin.Context){
	//c.JSON(http.StatusOK,gin.H{
	//	"hello":"hi world!",
	//})
	//name:="wangYinWen"
	//c.String(http.StatusOK,"Hello %s", name)
	s.Put(c,"Hello","Hello 世界")
	s.RespOK(c)
	return
}


