package modules

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/slimming/config"
)

type DuplicateScanner struct{}

func (s *DuplicateScanner) Name() string {
	return "duplicate"
}

func (s *DuplicateScanner) DefaultRisk() RiskLevel {
	return RiskHigh
}

func (s *DuplicateScanner) Scan(ctx context.Context, cfg *config.Config) ([]FileItem, error) {
	var items []FileItem

	minSize := int64(cfg.Duplicates.MinSizeMB) * 1024 * 1024

	fileHashes := make(map[string][]string)

	drives := []string{"C:\\"}

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
				if dirName == "Windows" || dirName == "Program Files" || dirName == "Program Files (x86)" {
					return filepath.SkipDir
				}
			}

			if !info.IsDir() && info.Size() >= minSize {
				ext := filepath.Ext(path)
				for _, exclude := range cfg.Duplicates.ExcludeExtensions {
					if ext == exclude {
						return nil
					}
				}

				hash, err := hashFile(path)
				if err == nil {
					fileHashes[hash] = append(fileHashes[hash], path)
				}
			}

			return nil
		})

		if err != nil && err != context.Canceled {
			continue
		}
	}

	for _, paths := range fileHashes {
		if len(paths) > 1 {
			for _, path := range paths[1:] {
				info, err := os.Stat(path)
				if err == nil {
					items = append(items, FileItem{
						Path:   path,
						Size:   info.Size(),
						Module: "duplicate",
						Risk:   RiskHigh,
					})
				}
			}
		}
	}

	return items, nil
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
