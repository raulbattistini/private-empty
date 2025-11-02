package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/raulbattistini/private-empty/pkg/server/requests"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(requests.HeaderAccessControlAllowOrigin, requests.HeaderAccessControlAllowOrigin)
		c.Writer.Header().Set(requests.HeaderAccessControlAllowCredentials, requests.HeaderAccessControlAllowCredentials)
		c.Writer.Header().Set(requests.HeaderAccessControlAllowCredentials, requests.HeaderAccessControlAllowHeaders)
		c.Writer.Header().Set(requests.HeaderAccessControlAllowMethods, requests.HeaderAccessControlAllowMethods)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
