package risk

import (
	"testing"

	"github.com/slimming/modules"
)

func TestRiskLevel_String(t *testing.T) {
	tests := []struct {
		level    RiskLevel
		expected string
	}{
		{Low, "low"},
		{Medium, "medium"},
		{High, "high"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.expected {
			t.Errorf("RiskLevel.String() = %v, want %v", got, tt.expected)
		}
	}
}

func TestClassifyRisk_DefaultRisk(t *testing.T) {
	engine := NewRiskEngine(nil)

	item := modules.FileItem{
		Path:   "C:\\temp\\file.txt",
		Module: "temp",
		Risk:   modules.RiskLow,
	}

	risk := engine.ClassifyRisk(item)
	if risk != Low {
		t.Errorf("expected Low risk, got %v", risk)
	}
}

func TestClassifyRisk_OverrideRisk(t *testing.T) {
	config := map[string]string{
		"temp": "high",
	}
	engine := NewRiskEngine(config)

	item := modules.FileItem{
		Path:   "C:\\temp\\file.txt",
		Module: "temp",
		Risk:   modules.RiskLow,
	}

	risk := engine.ClassifyRisk(item)
	if risk != High {
		t.Errorf("expected High risk, got %v", risk)
	}
}
