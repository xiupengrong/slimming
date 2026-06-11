package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type WinOldScanner struct{}

func (s *WinOldScanner) Name() string {
	return "winold"
}

func (s *WinOldScanner) DefaultRisk() RiskLevel {
	return RiskHigh
}

func (s *WinOldScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	winOldDir := filepath.Join(os.Getenv("windir"), "Windows.old")
	if _, err := os.Stat(winOldDir); os.IsNotExist(err) {
		return items, nil
	}

	err := filepath.Walk(winOldDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !info.IsDir() {
			items = append(items, FileItem{
				Path:   path,
				Size:   info.Size(),
				Module: "winold",
				Risk:   RiskHigh,
			})
		}

		return nil
	})

	if err != nil && err != context.Canceled {
		return items, nil
	}

	return items, nil
}
