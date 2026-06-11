package pipeline

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/slimming/modules"
	"github.com/slimming/risk"
)

func (p *Pipeline) Confirm(items []modules.FileItem, engine *risk.RiskEngine) []modules.FileItem {
	var confirmed []modules.FileItem

	mediumItems := []modules.FileItem{}
	for _, item := range items {
		level := engine.ClassifyRisk(modules.FileItem{
			Path:   item.Path,
			Module: item.Module,
			Risk:   modules.RiskLevel(item.Risk),
		})
		if level == risk.Medium {
			mediumItems = append(mediumItems, item)
		}
	}

	if len(mediumItems) == 0 {
		return items
	}

	fmt.Printf("\n⚠️  发现 %d 个中等风险项需要确认:\n", len(mediumItems))

	reader := bufio.NewReader(os.Stdin)

	for i, item := range mediumItems {
		fmt.Printf("\n[%d/%d] %s (%.2f MB)\n", i+1, len(mediumItems), item.Path, float64(item.Size)/1024/1024)
		fmt.Print("删除此文件? [Y]es / [N]o / [A]ll remaining / [Q]uit: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToUpper(input))

		switch input {
		case "Y", "YES":
			item.Confirmed = true
			confirmed = append(confirmed, item)
		case "A", "ALL":
			item.Confirmed = true
			confirmed = append(confirmed, item)
			for _, remaining := range mediumItems[i+1:] {
				remaining.Confirmed = true
				confirmed = append(confirmed, remaining)
			}
			return confirmed
		case "Q", "QUIT":
			return confirmed
		default:
			fmt.Println("跳过此文件")
		}
	}

	lowItems := []modules.FileItem{}
	for _, item := range items {
		level := engine.ClassifyRisk(modules.FileItem{
			Path:   item.Path,
			Module: item.Module,
			Risk:   modules.RiskLevel(item.Risk),
		})
		if level == risk.Low {
			lowItems = append(lowItems, item)
		}
	}

	return append(lowItems, confirmed...)
}
