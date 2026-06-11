package pipeline

import (
	"context"
	"sync"

	"github.com/slimming/modules"
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

			mu.Lock()
			items = append(items, scannerItems...)
			mu.Unlock()
		}(scanner)
	}

	wg.Wait()

	return items, nil
}
