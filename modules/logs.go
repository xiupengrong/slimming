package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type LogScanner struct{}

func (s *LogScanner) Name() string {
	return "logs"
}

func (s *LogScanner) DefaultRisk() RiskLevel {
	return RiskLow
}

func (s *LogScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	logDirs := []string{
		filepath.Join(os.Getenv("windir"), "Logs"),
		filepath.Join(os.Getenv("windir"), "System32", "LogFiles"),
		filepath.Join(os.Getenv("windir"), "CBS"),
	}

	for _, dir := range logDirs {
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
					Module: "logs",
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
