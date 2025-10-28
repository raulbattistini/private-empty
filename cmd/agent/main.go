package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/raulbattistini/private-empty/internal/agent"
	"github.com/raulbattistini/private-empty/internal/util/logger"
)

func main() {
	logger := logger.Log()

	if loc, err := time.LoadLocation("America/Sao_Paulo"); err == nil {
		time.Local = loc
	} else {
		logger.Warn("Failed to set timezone: %v", err)
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer cancel()

	cfg, err := agent.LoadConfig("config.yaml")
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		return
	}

	collectors := agent.BuildCollectors(cfg)
	exporters := agent.BuildExporters(cfg)

	if len(collectors) == 0 {
		logger.Warn("No collectors enabled in config")
	}
	if len(exporters) == 0 {
		logger.Warn("No exporters enabled in config")
	}

	a := agent.NewAgent(cfg, collectors, exporters)

	logger.Info("Starting monitoring agent...")
	logger.Info("Collectors: %d, Exporters: %d", len(collectors), len(exporters))

	if err := a.Run(ctx); err != nil {
		logger.Error("Agent error: %v", err)
		return
	}

	logger.Info("Agent stopped gracefully")
}
