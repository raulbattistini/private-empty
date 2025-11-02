package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
	"github.com/raulbattistini/private-empty/pkg/server/responses"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HandleHealth(c *gin.Context) {
	// mock
	healthStatus := &responses.HandleHealthRes{
		Status: enum.Health,
		Time:   time.Now().Format(time.RFC3339),
	}
	if err := c.ShouldBindJSON(healthStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
