package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/raulbattistini/private-empty/internal/event"
	"github.com/raulbattistini/private-empty/internal/util"
	"github.com/raulbattistini/private-empty/internal/util/metrics"
)

type MemoryCollector struct {
	interval time.Duration
}

func NewMemCollector(interval time.Duration) *MemoryCollector {
	return &MemoryCollector{
		interval: interval,
	}
}

func (m *MemoryCollector) Name() util.CollectorsTypes {
	return util.Memory
}

func (mem *MemoryCollector) Collect(ctx context.Context, out chan<- event.Event) {
	ticker := time.NewTicker(mem.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			usage := metrics.GetMemoryUsagePercent()
			msg := fmt.Sprintf("Mmemory usage: %.2f%%", usage)

			sev := util.Info
			if usage > 85 {
				sev = util.Warn
			}
			out <- event.NewEvent(mem.Name(), "mem_usage", sev, msg)
		}
	}
}
