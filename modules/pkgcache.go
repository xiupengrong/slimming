package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type PkgCacheScanner struct{}

func (s *PkgCacheScanner) Name() string {
	return "pkgcache"
}

func (s *PkgCacheScanner) DefaultRisk() RiskLevel {
	return RiskMedium
}

func (s *PkgCacheScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return items, nil
	}

	pkgDirs := []string{
		filepath.Join(homeDir, "AppData", "Local", "npm-cache"),
		filepath.Join(homeDir, "AppData", "Local", "pip", "Cache"),
		filepath.Join(homeDir, "go", "pkg", "mod"),
		filepath.Join(homeDir, ".cache"),
	}

	for _, dir := range pkgDirs {
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
					Module: "pkgcache",
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
