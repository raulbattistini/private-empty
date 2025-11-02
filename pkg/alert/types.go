package alert

import (
	"sync"
	"time"

	"github.com/raulbattistini/private-empty/pkg/metrics"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
)

type Rule struct {
	Source    enum.RuleSource
	Condition enum.Condition
	Threshold float64
	Duration  time.Duration
	Action    func(metrics.Metric)
}

type Alerter struct {
	mu     sync.RWMutex
	rules  []Rule
	states map[string]*AlertState // source -> state
}

type AlertState struct {
	triggered bool
	firstSeen time.Time
	lastValue float64
}
