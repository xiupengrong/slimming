package modules

import (
	"context"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type LargeFileScanner struct{}

func (s *LargeFileScanner) Name() string {
	return "largefile"
}

func (s *LargeFileScanner) DefaultRisk() RiskLevel {
	return RiskHigh
}

func (s *LargeFileScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	minSize := int64(cfg.LargeFiles.MinSizeMB) * 1024 * 1024

	drives := []string{"C:\\"}

	// 需要跳过的目录
	skipDirs := map[string]bool{
		"Windows":                  true,
		"Program Files":            true,
		"Program Files (x86)":      true,
		"ProgramData":              true,
		"$Recycle.Bin":             true,
		"System Volume Information": true,
		"Recovery":                 true,
		"Config.Msi":               true,
		"MSOCache":                 true,
	}

	for _, drive := range drives {
		err := filepath.Walk(drive, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if ctx.Err() != nil {
				return ctx.Err()
			}

			if info.IsDir() {
				dirName := info.Name()
				if skipDirs[dirName] {
					return filepath.SkipDir
				}

				for _, exclude := range cfg.LargeFiles.ExcludePaths {
					if matched, _ := filepath.Match(exclude, path); matched {
						return filepath.SkipDir
					}
				}
			}

			if !info.IsDir() && info.Size() >= minSize {
				items = append(items, FileItem{
					Path:   path,
					Size:   info.Size(),
					Module: "largefile",
					Risk:   RiskHigh,
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
