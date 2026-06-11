package modules

import (
	"context"

	"github.com/slimming/config"
)

type RiskLevel int

const (
	RiskLow RiskLevel = iota
	RiskMedium
	RiskHigh
)

type FileItem struct {
	Path      string
	Size      int64
	Module    string
	Risk      RiskLevel
	Confirmed bool
}

type Scanner interface {
	Name() string
	DefaultRisk() RiskLevel
	Scan(ctx context.Context, config *config.Config) ([]FileItem, error)
}
