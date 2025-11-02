package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raulbattistini/private-empty/internal/util/logger"
)

func Logger() gin.HandlerFunc {
	logger := logger.DefaultLogger
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)

		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("[%s] %d | %13v | %15s | %-7s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
