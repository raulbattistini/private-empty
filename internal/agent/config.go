package agent

import (
	"monitoring-agent/internal/collector"
	"monitoring-agent/internal/exporter"
	"os"
	"time"

	"github.com/stretchr/testify/assert/yaml"
)

type CollectorConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
}

type ExporterConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Endpoint time.Duration `yaml:"endpoint"`
}
type Config struct {
	Collectors map[string]CollectorConfig `yaml:"collectors"`
	Exporters  map[string]ExporterConfig  `yaml:"exporters"`
}

// avoiding 3rd parties as simple as the project is

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// helper factory funcs
func BuildCollectors(cfg *Config) []collector.Collector {
	var cols []collector.Collector
	if c, ok := cfg.Collectors["cpu"]; ok && c.Enabled {
		cols = append(cols, collector.NewCPUCollector(c.Interval))
	}
	if m, ok := cfg.Collectors["memory"]; ok && m.Enabled {
		cols = append(cols, collector.NewMemCollector(m.Interval))
	}
	if n, ok := cfg.Collectors["net"]; ok && n.Enabled {
		cols = append(cols, collector.NewNetCollector(n.Interval))
	}

	return cols
}

func BuildExporters(cfg *Config) []exporter.Exporter {
	var exps []exporter.Exporter
	if h, ok := cfg.Collectors["http"]; ok && h.Enabled {
		exps = append(exps, exporter.NewHttpExporter(h.Endpoint))
	}
	if e, ok := cfg.Exporters["log"]; ok && e.Enabled {
		exps = append(exps, exporter.NewLogExporter())
	}
	return exps
}
