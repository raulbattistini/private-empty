package responses

import (
	"time"

	"github.com/raulbattistini/private-empty/pkg/metrics"
)

type HandleMetricsRes struct {
	Count   int              `json:"count"`
	Metrics []metrics.Metric `json:"metrics"`
}

type HandleLatestRes struct {
	Timestamp time.Time        `json:"timestamp"`
	Metrics   []metrics.Metric `json:"metrics"`
}
