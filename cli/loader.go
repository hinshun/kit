package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"plugin"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/content"
)

type Loader struct {
	cli   *Cli
	store content.Store
}

func NewLoader(c *Cli, store content.Store) *Loader {
	return &Loader{
		cli:   c,
		store: store,
	}
}

func (l *Loader) GetCommand(ctx context.Context, cfg *config.Config, args []string) (*Command, error) {
	plugins := config.Plugins{
		{
			Name:     "kit",
			Manifest: cfg.Manifest,
			Plugins:  cfg.Plugins,
		},
	}

	args = append([]string{"kit"}, args...)
	manifest, depth, err := l.FindManifest(ctx, plugins, args)
	if err != nil {
		return nil, err
	}

	switch manifest.Type {
	case config.NamespaceManifest:
		var commands []*Command
		for _, plugin := range manifest.Plugins {
			submanifest, err := l.GetManifest(ctx, plugin)
			if err != nil {
				return nil, err
			}

			names := make([]string, len(args[1:depth])+1)
			copy(names, args[1:depth])
			names[len(names)-1] = plugin.Name

			commands = append(commands, &Command{
				CommandPath: names,
				Usage:       submanifest.Usage,
				Args:        submanifest.Args,
				Flags:       submanifest.Flags,
			})
		}

		return &Command{
			Action: func(ctx context.Context) error {
				return l.cli.PrintHelp(commands)
			},
		}, nil
	case config.CommandManifest:
		path, err := l.store.Get(ctx, manifest.Hash)
		if err != nil {
			return nil, err
		}

		constructor, err := OpenConstructor(path)
		if err != nil {
			return nil, err
		}

		kitCmd, err := constructor()
		if err != nil {
			return nil, err
		}

		cliCmd := &Command{
			CommandPath: args[1:depth],
			Usage:       manifest.Usage,
			Args:        manifest.Args,
			Flags:       manifest.Flags,
			Action: func(ctx context.Context) error {
				return kitCmd.Run(ctx)
			},
		}

		err = l.cli.Apply(cliCmd, kitCmd, args[depth:])
		if err != nil {
			l.cli.UsageError = err
		}

		return cliCmd, nil
	}

	return nil, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
}

func (l *Loader) FindManifest(ctx context.Context, plugins config.Plugins, args []string) (*config.Manifest, int, error) {
	if len(plugins) == 0 {
		return &config.Manifest{
			Type: config.NamespaceManifest,
		}, 0, nil
	}

	leafDepth := 0
	leaf, err := plugins.Walk(args, func(plugin config.Plugin, depth int) error {
		leafDepth = depth + 1
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	manifest, err := l.GetManifest(ctx, leaf)
	if err != nil {
		return nil, 0, err
	}

	switch manifest.Type {
	case config.NamespaceManifest:
		var depth int
		manifest, depth, err = l.FindManifest(ctx, manifest.Plugins, args[leafDepth:])
		if err != nil {
			return nil, 0, err
		}
		leafDepth += depth
	case config.CommandManifest:
		return manifest, leafDepth, nil
	default:
		return nil, 0, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
	}

	if manifest.Type == config.NamespaceManifest {
		for _, plugin := range leaf.Plugins {
			_, err = l.GetManifest(ctx, plugin)
			if err != nil {
				return nil, 0, err
			}

			manifest.Plugins = append(manifest.Plugins, plugin)
		}
	}

	return manifest, leafDepth, nil
}

func (l *Loader) GetManifest(ctx context.Context, plugin config.Plugin) (*config.Manifest, error) {
	if plugin.Manifest == "" {
		return &config.Manifest{
			Usage: plugin.Usage,
			Type:  config.NamespaceManifest,
		}, nil
	}

	path, err := l.store.Get(ctx, plugin.Manifest)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest config.Manifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
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

	constructor, ok := symbol.(*kit.Constructor)
	if !ok {
		return nil, fmt.Errorf("symbol not a (*kit.Constructor)")
	}

	return *constructor, nil
}
