package exporter

import (
	"context"

	"github.com/raulbattistini/private-empty/internal/event"
)

type Exporter interface {
	Export(ctx context.Context, e event.Event) error
}
