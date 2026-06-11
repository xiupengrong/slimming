package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type BrowserScanner struct{}

func (s *BrowserScanner) Name() string {
	return "browser"
}

func (s *BrowserScanner) DefaultRisk() RiskLevel {
	return RiskMedium
}

func (s *BrowserScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	browserCaches := []string{
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Cache"),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Cache"),
		filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Mozilla", "Firefox", "Profiles"),
	}

	for _, cacheDir := range browserCaches {
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
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
					Module: "browser",
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
