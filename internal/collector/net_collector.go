package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/raulbattistini/private-empty/internal/event"
	"github.com/raulbattistini/private-empty/internal/util"
)

type NetCollector struct {
	interval time.Duration
}

func NewNetCollector(interval time.Duration) *NetCollector {
	return &NetCollector{
		interval: interval,
	}
}

func (n *NetCollector) Name() util.CollectorsTypes {
	return util.Net
}

func (net *NetCollector) Collect(ctx context.Context, out chan<- event.Event) {
	ticker := time.NewTicker(net.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			usage := util.GetNetworkUsagePercent()
			msg := fmt.Sprintf("Net usage: %.2f%%", usage)

			sev := util.Info
			if usage > 85 {
				sev = util.Warn
			}
			out <- event.NewEvent(net.Name(), "net_usage", sev, msg)
		}
	}
}
