package modules

import (
	"context"
	"testing"

	"github.com/slimming/config"
)

func TestTempScanner_Name(t *testing.T) {
	scanner := &TempScanner{}
	if scanner.Name() != "temp" {
		t.Errorf("expected name 'temp', got '%s'", scanner.Name())
	}
}

func TestTempScanner_DefaultRisk(t *testing.T) {
	scanner := &TempScanner{}
	if scanner.DefaultRisk() != RiskLow {
		t.Errorf("expected RiskLow, got %v", scanner.DefaultRisk())
	}
}

func TestTempScanner_Scan(t *testing.T) {
	scanner := &TempScanner{}
	cfg := config.DefaultConfig()

	items, err := scanner.Scan(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	for _, item := range items {
		if item.Module != "temp" {
			t.Errorf("expected module 'temp', got '%s'", item.Module)
		}
	}
}
