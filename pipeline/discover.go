package pipeline

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/slimming/modules"
	"github.com/slimming/utils"
)

func (p *Pipeline) Discover(ctx context.Context) ([]modules.FileItem, error) {
	var (
		mu        sync.Mutex
		items     []modules.FileItem
		wg        sync.WaitGroup
		completed int32
		total     = int32(len(p.scanners))
	)

	for _, scanner := range p.scanners {
		wg.Add(1)
		go func(s modules.Scanner) {
			defer wg.Done()

			scannerItems, err := s.Scan(ctx, p.config)
			if err != nil {
				atomic.AddInt32(&completed, 1)
				return
			}

			var totalSize int64
			for _, item := range scannerItems {
				totalSize += item.Size
			}

			current := atomic.AddInt32(&completed, 1)
			percent := int(float64(current) / float64(total) * 100)
			bar := strings.Repeat("█", percent/5) + strings.Repeat("░", 20-percent/5)

			mu.Lock()
			fmt.Printf("  [%s] %3d%% Scanning %-12s %d files (%s)\n", bar, percent, s.Name()+"...", len(scannerItems), utils.FormatSize(totalSize))
			items = append(items, scannerItems...)
			mu.Unlock()
		}(scanner)
	}

	wg.Wait()

	return items, nil
}
