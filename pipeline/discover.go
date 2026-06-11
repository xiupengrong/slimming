package pipeline

import (
	"context"
	"fmt"
	"sync"

	"github.com/slimming/modules"
	"github.com/slimming/utils"
)

func (p *Pipeline) Discover(ctx context.Context) ([]modules.FileItem, error) {
	var (
		mu    sync.Mutex
		items []modules.FileItem
		wg    sync.WaitGroup
	)

	for _, scanner := range p.scanners {
		wg.Add(1)
		go func(s modules.Scanner) {
			defer wg.Done()

			scannerItems, err := s.Scan(ctx, p.config)
			if err != nil {
				return
			}

			var totalSize int64
			for _, item := range scannerItems {
				totalSize += item.Size
			}

			mu.Lock()
			fmt.Printf("  Scanning %s... %d files (%s)\n", s.Name(), len(scannerItems), utils.FormatSize(totalSize))
			items = append(items, scannerItems...)
			mu.Unlock()
		}(scanner)
	}

	wg.Wait()

	return items, nil
}
