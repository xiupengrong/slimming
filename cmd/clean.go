package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/slimming/config"
	"github.com/slimming/modules"
	"github.com/slimming/pipeline"
	"github.com/slimming/report"
	"github.com/slimming/risk"
)

var (
	dryRun    bool
	riskLevel string
	all       bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理C盘空间",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		if dryRun {
			fmt.Println("Dry run mode - no files will be deleted")
			return nil
		}

		engine := risk.NewRiskEngine(cfg.Risk)

		if !all {
			items = p.Confirm(items, engine)
		}

		toDelete := []modules.FileItem{}

		for _, item := range items {
			level := engine.ClassifyRisk(modules.FileItem{
				Path:   item.Path,
				Module: item.Module,
				Risk:   modules.RiskLevel(item.Risk),
			})

			if level == risk.Low {
				toDelete = append(toDelete, item)
			} else if level == risk.Medium && item.Confirmed {
				toDelete = append(toDelete, item)
			} else if level == risk.High && all {
				toDelete = append(toDelete, item)
			}
		}

		if len(toDelete) == 0 {
			fmt.Println("没有需要清理的文件")
			return nil
		}

		fmt.Printf("\n即将清理 %d 个文件...\n", len(toDelete))

		var deleted int
		var freedSize int64

		for _, item := range toDelete {
			if err := os.Remove(item.Path); err != nil {
				fmt.Printf("删除失败: %s - %v\n", item.Path, err)
				continue
			}
			deleted++
			freedSize += item.Size
		}

		fmt.Printf("\n清理完成: 删除 %d 个文件，释放 %s 空间\n", deleted, formatSize(freedSize))

		if err := report.GenerateReport(cfg, items, toDelete, engine); err != nil {
			fmt.Printf("生成报告失败: %v\n", err)
		}

		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "预览模式，不执行删除")
	cleanCmd.Flags().StringVar(&riskLevel, "risk", "", "只清理指定风险等级 (low/medium/high)")
	cleanCmd.Flags().BoolVar(&all, "all", false, "清理所有风险等级")
}
