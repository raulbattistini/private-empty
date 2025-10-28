package agent

import (
	"context"
	"sync"

	"github.com/raulbattistini/private-empty/internal/collector"
	"github.com/raulbattistini/private-empty/internal/event"
	"github.com/raulbattistini/private-empty/internal/exporter"
)

type Agent struct {
	cfg        *Config
	collectors []collector.ICollector
	exporters  []exporter.Exporter
	eventCh    chan event.Event
	wg         sync.WaitGroup
}

func NewAgent(cfg *Config, collectors []collector.ICollector, exporters []exporter.Exporter) *Agent {
	return &Agent{
		cfg:        cfg,
		collectors: collectors,
		exporters:  exporters,
		eventCh:    make(chan event.Event, 100),
	}
}

func (a *Agent) Run(ctx context.Context) error {
	for _, c := range a.collectors {
		a.wg.Add(1)
		go func(col collector.ICollector) {
			defer a.wg.Done()
			col.Collect(ctx, a.eventCh)
		}(c)
	}

	defer a.wg.Wait()
	go func() {
		defer a.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-a.eventCh:
				for _, ex := range a.exporters {
					go ex.Export(ctx, e)
				}
			}
		}
	}()

	<-ctx.Done()
	close(a.eventCh)
	a.wg.Wait()
	return nil
}
