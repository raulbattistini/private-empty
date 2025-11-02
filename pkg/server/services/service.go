package services

import (
	"context"

	"github.com/raulbattistini/private-empty/pkg/metrics"
)

type MetricsService interface {
	GetHistory(ctx context.Context, source string, limit *int) []metrics.Metric
	GetMetrics(ctx context.Context, source string, limit *int) ([]metrics.Metric, error)
	GetLatest(ctx context.Context, source string) []metrics.Metric
}
