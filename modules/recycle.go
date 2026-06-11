package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type RecycleScanner struct{}

func (s *RecycleScanner) Name() string {
	return "recycle"
}

func (s *RecycleScanner) DefaultRisk() RiskLevel {
	return RiskLow
}

func (s *RecycleScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	recycleDir := filepath.Join(os.Getenv("windir"), "System32", "config", "systemprofile", "Recycle.Bin")
	if _, err := os.Stat(recycleDir); os.IsNotExist(err) {
		return items, nil
	}

	err := filepath.Walk(recycleDir, func(path string, info os.FileInfo, err error) error {
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
				Module: "recycle",
				Risk:   RiskLow,
			})
		}

		return nil
	})

	if err != nil && err != context.Canceled {
		return items, nil
	}

	return items, nil
}
