package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type ThumbnailScanner struct{}

func (s *ThumbnailScanner) Name() string {
	return "thumbnail"
}

func (s *ThumbnailScanner) DefaultRisk() RiskLevel {
	return RiskLow
}

func (s *ThumbnailScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	thumbDirs := []string{
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "Windows", "Explorer"),
		filepath.Join(os.Getenv("windir"), "CSC"),
	}

	for _, dir := range thumbDirs {
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
					Module: "thumbnail",
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
