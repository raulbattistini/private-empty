package metrics

import (
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
)

type SqliteStore struct {
	Db   *sqlx.DB
	Rows sqlx.Rows
}

type MetricsStore struct {
	mu            sync.RWMutex
	latest        map[string]Metric
	history       []Metric
	maxSize       int
	sqliteStorage *SqliteStore
}

type Metric struct {
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
}

type MetricSink interface {
	Receive(m Metric)
	GetLatest() []Metric
	GetHistory(source enum.RuleSource, limit int) []Metric
}

type MultiMetricSink struct {
	Sinks []MetricSink
}

func NewMultiMetricSink(sinks ...MetricSink) *MultiMetricSink {
	return &MultiMetricSink{
		Sinks: sinks,
	}
}

func (m *MultiMetricSink) Receive(metric Metric) {
	for _, sink := range m.Sinks {
		sink.Receive(metric)
	}
}

func (m *MultiMetricSink) GetLatest() []Metric {
	if len(m.Sinks) > 0 {
		return m.Sinks[0].GetLatest()
	}
	return nil
}

func (m *MultiMetricSink) GetHistory(source enum.RuleSource, limit int) []Metric {
	if len(m.Sinks) > 0 {
		return m.Sinks[0].GetHistory(source, limit)
	}
	return nil
}

func NewMetricsStore(maxSize int) *MetricsStore {
	return &MetricsStore{
		latest:  make(map[string]Metric),
		history: make([]Metric, 0, maxSize),
		maxSize: maxSize,
	}
}

func (s *MetricsStore) Receive(m Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.latest[string(m.Source)] = m

	s.history = append(s.history, m)
	if len(s.history) > s.maxSize {
		s.history = s.history[1:]
	}
}

func (s *MetricsStore) GetLatest() []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Metric, 0, len(s.latest))
	for _, m := range s.latest {
		result = append(result, m)
	}
	return result
}

func (s *MetricsStore) GetHistory(source enum.RuleSource, limit int) []Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Metric, 0)
	count := 0

	for i := len(s.history) - 1; i >= 0 && count < limit; i-- {
		if source == "" || enum.RuleSource(s.history[i].Source) == enum.RuleSource(source) {
			result = append(result, s.history[i])
			count++
		}
	}

	return result
}

// started already as a circular buffer, think similar to a linked list as in determining the last node the tail and so on
