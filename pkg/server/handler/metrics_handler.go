package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raulbattistini/private-empty/pkg/metrics"
	"github.com/raulbattistini/private-empty/pkg/server/app_errors"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
	"github.com/raulbattistini/private-empty/pkg/server/requests"
	"github.com/raulbattistini/private-empty/pkg/server/responses"
)

// MetricsHandler handles all metrics-related endpoints
type MetricsHandler struct {
	store *metrics.MetricsStore
}

func NewMetricsHandler(store *metrics.MetricsStore) *MetricsHandler {
	return &MetricsHandler{store: store}
}

func (h *MetricsHandler) HandleMetrics(c *gin.Context) {
	source := c.Query("source")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		if err = c.ShouldBindJSON(&app_errors.InvalidLimitNumRequestStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	src, err := enum.ParseMetrics(source)
	invalidMetr := app_errors.ErrorResponse{
		Status:  &requests.BadRequest,
		Message: err.Error(),
	}
	if err != nil {
		if err = c.ShouldBindJSON(&invalidMetr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	history := h.store.GetHistory(*src, limit)

	c.JSON(http.StatusOK, gin.H{
		"count":   len(history),
		"metrics": history,
	})
}

func (h *MetricsHandler) HandleHistory(c *gin.Context) {
	source := c.Query("source")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		if err = c.ShouldBindJSON(&app_errors.InvalidLimitNumRequestStr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	src, err := enum.ParseMetrics(source)
	invalidMetr := app_errors.ErrorResponse{
		Status:  &requests.BadRequest,
		Message: err.Error(),
	}

	if err != nil {
		if err = c.ShouldBindJSON(&invalidMetr); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	history := h.store.GetHistory(*src, limit)

	fullHistory := &responses.HandleMetricsRes{
		Count:   len(history),
		Metrics: history,
	}
	if err := c.ShouldBindJSON(fullHistory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fullHistory)
}

func (h *MetricsHandler) HandleLatest(c *gin.Context) {
	latest := h.store.GetLatest()

	latestRes := &responses.HandleLatestRes{
		Timestamp: time.Now(),
		Metrics:   latest,
	}

	if err := c.ShouldBindJSON(latestRes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, latestRes)
}
