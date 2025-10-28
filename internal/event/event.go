package event

import (
	"time"

	"github.com/raulbattistini/private-empty/internal/util"
)

type Event struct {
	Source    util.CollectorsTypes `json:"source"`
	Type      string               `json:"type"`
	Severity  util.EventSeverity   `json:"severity"`
	Message   string               `json:"message"`
	Timestamp time.Time            `json:"timestamp"`
}

func NewEvent(source util.CollectorsTypes, eventType string, sev util.EventSeverity, msg string) Event {
	return Event{
		Source:    source,
		Type:      eventType,
		Severity:  sev,
		Message:   msg,
		Timestamp: time.Now(),
	}
}
