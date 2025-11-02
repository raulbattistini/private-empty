package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/raulbattistini/private-empty/pkg/metrics"
)

type Server struct {
	config Config
	store  *metrics.MetricsStore
	router *gin.Engine
	srv    *http.Server
}

func NewServer(cfg Config, store *metrics.MetricsStore) *Server {
	if cfg.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		config: cfg,
		store:  store,
		router: gin.New(),
	}

	s.setupRouter()

	return s
}

func (s *Server) Start() error {
	addr := s.config.Addr()

	s.srv = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s.srv.ListenAndServe()
}

// graceful shutdown is real important to adopt from early stages
func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv != nil {
		return s.srv.Shutdown(ctx)
	}
	return nil
}
