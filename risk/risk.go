package risk

import (
	"github.com/slimming/modules"
)

type RiskLevel int

const (
	Low RiskLevel = iota
	Medium
	High
)

func (r RiskLevel) String() string {
	switch r {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	default:
		return "unknown"
	}
}

type RiskEngine struct {
	overrides map[string]string
}

func NewRiskEngine(overrides map[string]string) *RiskEngine {
	if overrides == nil {
		overrides = make(map[string]string)
	}
	return &RiskEngine{overrides: overrides}
}

func (e *RiskEngine) ClassifyRisk(item modules.FileItem) RiskLevel {
	if override, ok := e.overrides[item.Module]; ok {
		switch override {
		case "low":
			return Low
		case "medium":
			return Medium
		case "high":
			return High
		}
	}

	switch item.Risk {
	case modules.RiskLow:
		return Low
	case modules.RiskMedium:
		return Medium
	case modules.RiskHigh:
		return High
	default:
		return Low
	}
}
