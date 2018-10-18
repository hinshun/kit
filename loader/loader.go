package loader

import (
	"context"
	"fmt"
	"plugin"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/content"
)

type Loader struct {
	ContentStore content.Store
}

func New(store content.Store) *Loader {
	return &Loader{
		ContentStore: store,
	}
}

func (l *Loader) Load(ctx context.Context, plugins config.Plugins, args []string) (kit.Command, error) {
	var err error
	args, err = plugins.Walk(args, func(plugin config.Plugin) error {
		fmt.Printf("visited plugin '%s'\n", plugin.Name)
		return nil
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("args '%s'\n", args)

	return nil, fmt.Errorf("end")
}

func OpenConstructor(path string) (kit.Constructor, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symbol, err := p.Lookup("New")
	if err != nil {
		return nil, err
	}

	constructor, ok := symbol.(kit.Constructor)
	if !ok {
		return nil, fmt.Errorf("symbol not a kit.Constructor")
	}

	return constructor, nil
}
