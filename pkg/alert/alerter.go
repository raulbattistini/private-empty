package alert

import (
	"time"

	"github.com/raulbattistini/private-empty/pkg/metrics"
	"github.com/raulbattistini/private-empty/pkg/server/models/enum"
)

func NewAlerter() *Alerter {
	return &Alerter{
		rules:  make([]Rule, 0),
		states: make(map[string]*AlertState),
	}
}

func (a *Alerter) AddRule(rule Rule) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rules = append(a.rules, rule)
}

func (a *Alerter) Check(m metrics.Metric) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, rule := range a.rules {
		if rule.Source != enum.RuleSource(m.Source) {
			continue
		}

		state, exists := a.states[string(m.Source)]
		if !exists {
			state = &AlertState{}
			a.states[string(m.Source)] = state
		}

		violated := a.checkCondition(rule, m.Value)

		if violated {
			if !state.triggered {
				state.firstSeen = time.Now()
				state.triggered = true
			}

			// Check if duration threshold met
			if time.Since(state.firstSeen) >= rule.Duration {
				go rule.Action(m)
			}
		} else {
			// Reset state
			state.triggered = false
		}

		state.lastValue = m.Value
	}
}

func (a *Alerter) checkCondition(rule Rule, value float64) bool {
	switch rule.Condition {
	case enum.Above:
		return value > rule.Threshold
	case enum.Below:
		return value < rule.Threshold
	default:
		return false
	}
}

type AlertingMetricSink struct {
	wrapped metrics.MetricSink
	alerter *Alerter
}

func (a *AlertingMetricSink) Receive(m metrics.Metric) {
	a.alerter.Check(m)
	a.wrapped.Receive(m)
}

func (a *AlertingMetricSink) GetLatest() []metrics.Metric {
	return a.wrapped.GetLatest()
}

func (a *AlertingMetricSink) GetHistory(source enum.RuleSource, limit int) []metrics.Metric {
	return a.wrapped.GetHistory(source, limit)
}
