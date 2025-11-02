package enum

import (
	"strings"

	"github.com/raulbattistini/private-empty/pkg/server/app_errors"
)

type RuleSource string

const (
	CPU    RuleSource = "cpu"
	Memory RuleSource = "memory"
	Net    RuleSource = "net"
)

var RuleSourceMap = map[string]RuleSource{
	"cpu":    CPU,
	"memory": Memory,
	"net":    Net,
}

func ParseMetrics(s string) (*RuleSource, error) {
	s = strings.ToLower(s)
	if rule, ok := RuleSourceMap[s]; ok {
		return &rule, nil
	}
	return nil, app_errors.InvalidMetricRequest
}

// "cpu", "memory", "net"
