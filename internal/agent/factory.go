package agent

import (
	"github.com/raulbattistini/private-empty/internal/collector"
	"github.com/raulbattistini/private-empty/internal/exporter"
)

func BuildCollectors(cfg *Config) []collector.ICollector {
	var collectors []collector.ICollector

	if c, ok := cfg.Collectors["cpu"]; ok && c.Enabled {
		collectors = append(collectors, collector.NewCPUCollector(c.Interval))
	}
	if m, ok := cfg.Collectors["memory"]; ok && m.Enabled {
		collectors = append(collectors, collector.NewMemCollector(m.Interval))
	}
	if n, ok := cfg.Collectors["net"]; ok && n.Enabled {
		collectors = append(collectors, collector.NewNetCollector(n.Interval))
	}

	return collectors
}

func BuildExporters(cfg *Config) []exporter.Exporter {
	var exps []exporter.Exporter

	if e, ok := cfg.Exporters["http"]; ok && e.Enabled {
		exps = append(exps, exporter.NewHTTPExporter(e.Endpoint))
	}
	if e, ok := cfg.Exporters["log"]; ok && e.Enabled {
		exps = append(exps, exporter.NewLogExporter())
	}

	return exps
}
