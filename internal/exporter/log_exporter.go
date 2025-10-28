package exporter

import (
	"context"

	"github.com/raulbattistini/private-empty/internal/event"
	"github.com/raulbattistini/private-empty/internal/util/logger"
)

type LogExporter struct{}

func NewLogExporter() *LogExporter {
	return &LogExporter{}
}

func (e *LogExporter) Export(ctx context.Context, ev event.Event) error {
	logger.Log().Info(
		"Event: Source=%s, Severity=%s, Message='%s', Type=%v",
		ev.Source,
		ev.Severity,
		ev.Message,
		ev.Type,
	)

	return nil
}
