package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/slimming/config"
	"github.com/slimming/modules"
	"github.com/slimming/risk"
)

type Report struct {
	Timestamp   time.Time         `json:"timestamp"`
	TotalItems  int               `json:"total_items"`
	TotalSize   int64             `json:"total_size"`
	DeletedItems int              `json:"deleted_items"`
	FreedSize   int64             `json:"freed_size"`
	Items       []ReportItem      `json:"items"`
}

type ReportItem struct {
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	Module    string `json:"module"`
	Risk      string `json:"risk"`
	Confirmed bool   `json:"confirmed"`
	Deleted   bool   `json:"deleted"`
}

func GenerateReport(cfg *config.Config, items []modules.FileItem, deletedItems []modules.FileItem, engine *risk.RiskEngine) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %w", err)
	}

	reportDir := filepath.Join(homeDir, ".slimming", "reports")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return fmt.Errorf("创建报告目录失败: %w", err)
	}

	report := Report{
		Timestamp: time.Now(),
	}

	deletedPaths := make(map[string]bool)
	for _, item := range deletedItems {
		deletedPaths[item.Path] = true
	}

	for _, item := range items {
		level := engine.ClassifyRisk(modules.FileItem{
			Path:   item.Path,
			Module: item.Module,
			Risk:   modules.RiskLevel(item.Risk),
		})

		reportItem := ReportItem{
			Path:      item.Path,
			Size:      item.Size,
			Module:    item.Module,
			Risk:      level.String(),
			Confirmed: item.Confirmed,
			Deleted:   deletedPaths[item.Path],
		}

		report.Items = append(report.Items, reportItem)
		report.TotalItems++
		report.TotalSize += item.Size

		if deletedPaths[item.Path] {
			report.DeletedItems++
			report.FreedSize += item.Size
		}
	}

	filename := fmt.Sprintf("report_%s.json", time.Now().Format("2006-01-02_15-04-05"))
	reportPath := filepath.Join(reportDir, filename)

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化报告失败: %w", err)
	}

	if err := os.WriteFile(reportPath, data, 0644); err != nil {
		return fmt.Errorf("保存报告失败: %w", err)
	}

	fmt.Printf("\n📄 报告已保存: %s\n", reportPath)

	return nil
}
