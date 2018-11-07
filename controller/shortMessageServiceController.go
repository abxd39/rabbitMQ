package controller

import (
	"github.com/gin-gonic/gin"
)

type ShortMessageService struct {

}

func (s *ShortMessageService) Router(r *gin.Engine)  {
	sms := r.Group("/sms")
	{
		//sms.GET("/Test/hello", action.Hello)
		_=sms
	}

}


