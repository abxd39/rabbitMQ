package common

import (
	"github.com/gin-gonic/gin"
	Log "github.com/sirupsen/logrus"
	"net/http"
)

type HandlerFunc func(*gin.Context) error

func Handle(f HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			Log.Errorln(err)
			RenderJSONWithError(c, err)
		} else if c.Writer.Size() < 0 { //当action没有返回结果，也没有返回错误时，返回给客户端一个默认结构体
			RenderJSON(c, gin.H{})
		}
	}
}

// 处理跨域;
func HandleCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, timestamp, username, signature")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
