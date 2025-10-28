package agent

import (
	"monitoring-agent/internal/collector"
	"monitoring-agent/internal/event"
	"monitoring-agent/internal/exporter"
	"sync"
)

type Agent struct {
	cfg        *Config
	collectors []collector.Collector
	exporters  []exporter.Exporter
	eventCh    chan event.Event
	wg         sync.WaitGroup
}

func (a *Agent) NewAgent(cfg *Config, collectors []collector.Collector, exporters []exporter.Exporter) *Agent {
	return &Agent{
		cfg:        cfg,
		collectors: collectors,
		exporters:  exporters,
		eventCh:    make(chan event.Event, 100),
	}
}

func (a *Agent) Run(ctx context.context) error {
	for _, c := range a.collectors {
		a.wg.Add(1)
		go func(col collector.Collector) {
			defer a.wg.Done()
			col.Collect(ctx, a.eventCh)
		}(c)
	}

	a.wg.Add(1)
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
