package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type WUpdateScanner struct{}

func (s *WUpdateScanner) Name() string {
	return "wupdate"
}

func (s *WUpdateScanner) DefaultRisk() RiskLevel {
	return RiskMedium
}

func (s *WUpdateScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	wupdateDirs := []string{
		filepath.Join(os.Getenv("windir"), "SoftwareDistribution", "Download"),
		filepath.Join(os.Getenv("windir"), "SoftwareDistribution", "DataStore"),
	}

	for _, dir := range wupdateDirs {
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
					Module: "wupdate",
					Risk:   RiskMedium,
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
