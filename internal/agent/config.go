package agent

import (
	"os"
	"time"

	"github.com/raulbattistini/private-empty/internal/util"
	"github.com/raulbattistini/private-empty/internal/util/logger"
	"gopkg.in/yaml.v3" // Fix: wrong import
)

type CollectorConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
}

type ExporterConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"` // Fix: was time.Duration, should be string
}

type Config struct {
	Collectors   map[string]CollectorConfig `yaml:"collectors"`
	Exporters    map[string]ExporterConfig  `yaml:"exporters"`
	Notification struct {
		Annoying bool `yaml:"annoying"`
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	currentOS := util.ParseRuntimeOS()
	logger.Log().Info("Agent starting on detected OS: %s", currentOS)

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
