package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raulbattistini/private-empty/internal/util/logger"
	"github.com/raulbattistini/private-empty/pkg/server/app_errors"
	"github.com/raulbattistini/private-empty/pkg/server/models"
)

func Recovery() gin.HandlerFunc {
	logger := logger.Log()

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered: %v", err)
			}

			middlewareErr := models.DefaultMiddlewareError{
				Error: app_errors.ErrInternalServeErrorStr,
			}

			c.JSON(http.StatusInternalServerError, middlewareErr)

		}()
	}
}
