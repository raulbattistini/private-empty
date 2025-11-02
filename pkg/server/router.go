package server

import (
	"github.com/gin-gonic/gin"
	"github.com/raulbattistini/private-empty/pkg/server/handler"
	"github.com/raulbattistini/private-empty/pkg/server/middleware"
)

func (s *Server) setupRouter() {
	// Global middleware
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recovery())
	s.router.Use(middleware.CORS())

	// Create handlers
	metricsHandler := handler.NewMetricsHandler(s.store)
	healthHandler := handler.NewHealthHandler()

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Metrics endpoints
		metrics := v1.Group("/metrics")
		{
			metrics.GET("", metricsHandler.HandleMetrics)
			metrics.GET("/latest", metricsHandler.HandleLatest)
			metrics.GET("/history/:source", metricsHandler.HandleHistory)
		}

		// Health endpoint
		v1.GET("/health", healthHandler.HandleHealth)
	}

	// Root redirect
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Metrics Monitoring API",
			"version": "1.0.0",
			"docs":    "/api/v1/metrics/latest",
		})
	})
}
