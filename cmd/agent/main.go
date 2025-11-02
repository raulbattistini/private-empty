package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/raulbattistini/private-empty/internal/agent"
	"github.com/raulbattistini/private-empty/internal/exporter"
	"github.com/raulbattistini/private-empty/internal/util/logger"
	"github.com/raulbattistini/private-empty/pkg/alert"
	"github.com/raulbattistini/private-empty/pkg/metrics"
	"github.com/raulbattistini/private-empty/pkg/server"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
	"github.com/raulbattistini/private-empty/pkg/storage"
)

var (
	// could be dismissed
	audioOnce sync.Once

	GlobalSampleRate = beep.SampleRate(44100)
)

func InitSpeaker() {
	log := logger.Log()

	audioOnce.Do(func() {
		if err := speaker.Init(GlobalSampleRate, GlobalSampleRate.N(time.Millisecond*500)); err != nil {
			log.Error("Failed to initialize global audio speaker. Sound alerts disabled: %v", err)
			return
		}

		log.Info("Audio speaker initialized successfully.")
	})
}

func main() {
	log := logger.Log()

	InitSpeaker()
	// Load configs
	cfg, err := agent.LoadConfig("config.yaml")
	if err != nil {
		log.Error("Failed to load config: %v", err)
		return
	}

	srvCfg, err := server.LoadConfig("server.yaml")
	if err != nil {
		log.Error("Failed to load server config: %v", err)
		return
	}

	// Create stores
	metricsStore := metrics.NewMetricsStore(srvCfg.MaxStoreSize)

	sqliteStore, err := storage.NewSQLiteStore("metrics.db")
	if err != nil {
		log.Error("Failed to create SQLite store: %v", err)
		return
	}
	defer sqliteStore.Close()

	defer speaker.Close()
	multiSink := metrics.NewMultiMetricSink(metricsStore, sqliteStore)

	// Setup alerter

	alerter := alert.NewAlerter()
	alerter.AddRule(alert.Rule{
		Source:    "cpu",
		Condition: enum.Above,
		Threshold: 0.05,
		Duration:  1 * time.Second,
		Action: func(m metrics.Metric) {
			log.Warn("🚨 ALERT: CPU above 5%% for 6s: %.2f%%", m.Value)
			alert.DesktopNotifier(m)
		},
	})

	alerter.AddRule(alert.Rule{
		Source:    "memory",
		Condition: enum.Above,
		Threshold: 256.0, // 4GB
		Duration:  30 * time.Second,
		Action: func(m metrics.Metric) {
			log.Warn("🚨 ALERT: Memory above 250GB for 30s: %.2f MB", m.Value)
			if cfg.Notification.Annoying {
				alert.AnnoyingDesktopNotifier(m)
			} else {
				alert.DesktopNotifier(m) // OS notification!
			}
		},
	})

	alertingSink := &alertingMetricSink{
		wrapped:  multiSink,
		alerter:  alerter,
		annoying: cfg.Notification.Annoying,
	}

	// Build agent components
	collectors := agent.BuildCollectors(cfg)
	exporters := agent.BuildExporters(cfg)
	exporters = append(exporters, exporter.NewMetricsExporter(alertingSink))

	a := agent.NewAgent(cfg, collectors, exporters)

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	// Start HTTP server in background
	srv := server.NewServer(srvCfg, metricsStore)
	go func() {
		log.Info("Starting HTTP server on %s", srvCfg.Addr())
		if err := srv.Start(); err != nil {
			log.Error("HTTP server error: %v", err)
		}
	}()

	// Start agent in background
	go func() {
		log.Info("Starting monitoring agent...")
		log.Info("Collectors: %d, Exporters: %d", len(collectors), len(exporters))
		log.Info("HTTP API: http://localhost:8080/api/v1/metrics/latest")

		if err := a.Run(ctx); err != nil {
			log.Error("Agent error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info("Shutting down gracefully...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error: %v", err)
	}

	log.Info("Server stopped")
}

// Helper types

type alertingMetricSink struct {
	wrapped  metrics.MetricSink
	alerter  *alert.Alerter
	annoying bool
}

func (a *alertingMetricSink) Receive(m metrics.Metric) {
	a.alerter.Check(m)
	a.wrapped.Receive(m)
}

func (a *alertingMetricSink) GetLatest() []metrics.Metric {
	return a.wrapped.GetLatest()
}

func (a *alertingMetricSink) GetHistory(source enum.RuleSource, limit int) []metrics.Metric {
	return a.wrapped.GetHistory(source, limit)
}
