package models

import (
	"time"

	"github.com/raulbattistini/private-empty/pkg/metrics"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
)

type LatestMetric struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int       `json:"latest"`
	Metrics   *metrics.MetricsStore
}

type GetMetrics struct {
	Count   *int `json:"count"`
	Metrics *[]metrics.Metric
}

type HealthStatus struct {
	Status enum.ApiHealthStatus `json:"apiStatus"`
	Time   time.Time            `json:"time"`
}
