package server

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host         string
	Port         int
	MaxStoreSize int
	Production   bool
}

func DefaultConfig() Config {
	port := 8080
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err != nil {
			port = p
		}
	}

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}
	production := os.Getenv("ENV") == "production"

	return Config{
		Host:         host,
		Port:         port,
		MaxStoreSize: 1000,
		Production:   production,
	}
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, nil
	}
	var yamlCfg YamlCfg
	if err := yaml.Unmarshal(data, &yamlCfg); err != nil {
		return cfg, err
	}

	if yamlCfg.Server.Host != "" {
		cfg.Host = yamlCfg.Server.Host
	}

	if yamlCfg.Server.Port != 0 {
		cfg.Port = yamlCfg.Server.Port
	}

	if yamlCfg.Server.MaxStoreSize != 0 {
		cfg.MaxStoreSize = yamlCfg.Server.MaxStoreSize
	}

	cfg.Production = yamlCfg.Server.Production

	if envHost := os.Getenv("SERVER_HOST"); envHost != "" {
		cfg.Host = envHost
	}

	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			cfg.Port = p
		}
	}

	return cfg, nil
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
