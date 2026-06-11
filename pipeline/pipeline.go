package pipeline

import (
	"context"

	"github.com/slimming/config"
	"github.com/slimming/modules"
)

type Pipeline struct {
	scanners []modules.Scanner
	config   *config.Config
}

func NewPipeline(config *config.Config) *Pipeline {
	return &Pipeline{
		scanners: []modules.Scanner{
			&modules.TempScanner{},
			&modules.RecycleScanner{},
			&modules.ThumbnailScanner{},
			&modules.LogScanner{},
			&modules.BrowserScanner{},
			&modules.WUpdateScanner{},
			&modules.PkgCacheScanner{},
			&modules.WinOldScanner{},
			&modules.LargeFileScanner{},
			&modules.DuplicateScanner{},
		},
		config: config,
	}
}

func (p *Pipeline) Run(ctx context.Context) ([]modules.FileItem, error) {
	items, err := p.Discover(ctx)
	if err != nil {
		return nil, err
	}

	items = p.Classify(items)

	return items, nil
}
