package db_models

import (
	"time"

	"github.com/raulbattistini/private-empty/pkg/metrics"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
)

type StorageMetric struct {
	Timestamp time.Time       `json:"timestamp"`
	Source    enum.RuleSource `json:"source"`
	Type      string          `json:"type"`
	Value     float64         `json:"value"`
	Severity  string          `json:"severity"`
	Message   string          `json:"message"`
}

func (sm *StorageMetric) ToMetricStd() metrics.Metric {
	return metrics.Metric{}
}
