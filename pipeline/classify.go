package pipeline

import (
	"github.com/slimming/modules"
	"github.com/slimming/risk"
)

func (p *Pipeline) Classify(items []modules.FileItem) []modules.FileItem {
	engine := risk.NewRiskEngine(p.config.Risk)

	for i := range items {
		items[i].Risk = modules.RiskLevel(engine.ClassifyRisk(items[i]))
	}

	return items
}
