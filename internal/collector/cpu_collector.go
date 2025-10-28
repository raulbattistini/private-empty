package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/raulbattistini/private-empty/internal/event"
	"github.com/raulbattistini/private-empty/internal/util"
)

type CPUCollector struct {
	interval time.Duration
}

func NewCPUCollector(interval time.Duration) *CPUCollector {
	return &CPUCollector{
		interval: interval,
	}
}

func (c *CPUCollector) Name() util.CollectorsTypes {
	return util.CPU
}

func (c *CPUCollector) Collect(ctx context.Context, out chan<- event.Event) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			usage := util.GetCPUUsage()
			msg := fmt.Sprintf("CPU usage: %.2f%%", usage)

			sev := util.Info
			if usage > 85 {
				sev = util.Warn
			}
			out <- event.NewEvent(c.Name(), "cpu_usage", sev, msg)
		}
	}
}
