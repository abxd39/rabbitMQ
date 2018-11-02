package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"sctek.com/typhoon/th-platform-gateway/common"
)

func Logger() gin.HandlerFunc {
	return LoggerWithWriter()
}
func LoggerWithWriter(notlogged ...string) gin.HandlerFunc {
	var skip map[string]struct{}
	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		// Process request
		c.Next()
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)
			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
			if raw != "" {
				path = path + "?" + raw
			}
			deviceSerials := c.GetString("device_serials")
			if deviceSerials == "" {
				deviceSerials = "-"
			}
			common.Log.Infof("[GIN] %s | %3d | %13v | %15s | %s %s\n%s",
				deviceSerials,
				statusCode,
				latency,
				clientIP,
				method,
				path,
				comment,
			)
		}
	}
}
