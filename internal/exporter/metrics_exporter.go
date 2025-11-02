package exporter

import (
	"context"
	"fmt"

	"github.com/raulbattistini/private-empty/internal/event"
	"github.com/raulbattistini/private-empty/pkg/metrics"
)

type MetricsExporter struct {
	sink metrics.MetricSink
}

func NewMetricsExporter(sink metrics.MetricSink) *MetricsExporter {
	return &MetricsExporter{sink: sink}
}

func (e *MetricsExporter) Export(ctx context.Context, ev event.Event) error {
	m := metrics.Metric{
		Timestamp: ev.Timestamp,
		Source:    string(ev.Source),
		Type:      ev.Type,
		Value:     extractValue(ev.Message),
		Severity:  string(ev.Severity),
		Message:   ev.Message,
	}

	e.sink.Receive(m)
	return nil
}

func extractValue(msg string) float64 {
	var value float64
	fmt.Sscanf(msg, "%*s %*s %f", &value)
	return value
}
