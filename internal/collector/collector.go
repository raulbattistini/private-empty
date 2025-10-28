package collector

import (
	"context"

	"github.com/raulbattistini/private-empty/internal/event"
	"github.com/raulbattistini/private-empty/internal/util"
)

type ICollector interface {
	Name() util.CollectorsTypes
	Collect(ctx context.Context, out chan<- event.Event)
}
