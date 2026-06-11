package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type TempScanner struct{}

func (s *TempScanner) Name() string {
	return "temp"
}

func (s *TempScanner) DefaultRisk() RiskLevel {
	return RiskLow
}

func (s *TempScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	tempDirs := []string{
		os.TempDir(),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Temp"),
		filepath.Join(os.Getenv("windir"), "Temp"),
	}

	for _, dir := range tempDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
					Module: "temp",
					Risk:   RiskLow,
				})
			}

			return nil
		})

		if err != nil && err != context.Canceled {
			continue
		}
	}

	return items, nil
}
