package modules

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sync"

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

	// 第一阶段：按文件大小分组
	sizeGroups := make(map[int64][]string)

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
				if dirName == "Windows" || dirName == "Program Files" || dirName == "Program Files (x86)" ||
					dirName == "ProgramData" || dirName == "$Recycle.Bin" || dirName == "System Volume Information" {
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
				sizeGroups[info.Size()] = append(sizeGroups[info.Size()], path)
			}

			return nil
		})

		if err != nil && err != context.Canceled {
			continue
		}
	}

	// 第二阶段：只对相同大小的文件计算哈希（并发）
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // 限制并发数

	for _, paths := range sizeGroups {
		if len(paths) < 2 {
			continue
		}

		fileHashes := make(map[string][]string)

		for _, path := range paths {
			wg.Add(1)
			sem <- struct{}{}
			go func(p string) {
				defer wg.Done()
				defer func() { <-sem }()

				hash, err := hashFile(p)
				if err == nil {
					mu.Lock()
					fileHashes[hash] = append(fileHashes[hash], p)
					mu.Unlock()
				}
			}(path)
		}

		wg.Wait()

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
