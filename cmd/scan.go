package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/slimming/config"
	"github.com/slimming/modules"
	"github.com/slimming/pipeline"
	"github.com/slimming/risk"
	"github.com/slimming/utils"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "扫描可清理内容",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\n🔍 Scanning C:\\ drive...\n")

		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("加载配置失败: %w", err)
		}

		p := pipeline.NewPipeline(cfg)
		items, err := p.Discover(context.Background())
		if err != nil {
			return fmt.Errorf("扫描失败: %w", err)
		}

		items = p.Classify(items)

		engine := risk.NewRiskEngine(cfg.Risk)
		lowCount, mediumCount, highCount := 0, 0, 0
		lowSize, mediumSize, highSize := int64(0), int64(0), int64(0)

		for _, item := range items {
			level := engine.ClassifyRisk(modules.FileItem{
				Path:   item.Path,
				Module: item.Module,
				Risk:   modules.RiskLevel(item.Risk),
			})

			switch level {
			case risk.Low:
				lowCount++
				lowSize += item.Size
			case risk.Medium:
				mediumCount++
				mediumSize += item.Size
			case risk.High:
				highCount++
				highSize += item.Size
			}
		}

		fmt.Printf("\n📊 Scan Summary\n\n")
		fmt.Printf("LOW risk     (auto-clean)         %d items  %s\n", lowCount, utils.FormatSize(lowSize))
		fmt.Printf("MEDIUM risk  (needs confirm)       %d items  %s\n", mediumCount, utils.FormatSize(mediumSize))
		fmt.Printf("HIGH risk    (skipped by default)  %d items  %s\n", highCount, utils.FormatSize(highSize))
		fmt.Printf("\nTotal: %d items, %s potential space\n", len(items), utils.FormatSize(lowSize+mediumSize+highSize))

		return nil
	},
}
