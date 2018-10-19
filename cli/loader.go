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

func (l *Loader) LoadCommand(ctx context.Context, plugins config.Plugins, args []string) (*Command, error) {
	manifest, depth, err := l.LoadManifest(ctx, plugins, args)
	if err != nil {
		return nil, err
	}

	switch manifest.Type {
	case config.NamespaceManifest:
		var commands []*Command
		for _, plugin := range manifest.Plugins {
			submanifest, err := l.visit(ctx, plugin)
			if err != nil {
				return nil, err
			}

			commands = append(commands, &Command{
				Names: append(args[:depth], plugin.Name),
				Usage: submanifest.Usage,
				Args:  submanifest.Args,
				Flags: submanifest.Flags,
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

		command, err := constructor(l.cli)
		if err != nil {
			return nil, err
		}

		return &Command{
			Names: args[:depth],
			Usage: manifest.Usage,
			Args:  manifest.Args,
			Flags: manifest.Flags,
			Action: func(ctx context.Context) error {
				return command.Run(ctx)
			},
		}, nil
	}

	return nil, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
}

func (l *Loader) LoadManifest(ctx context.Context, plugins config.Plugins, args []string) (*config.Manifest, int, error) {
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

	manifest, err := l.visit(ctx, leaf)
	if err != nil {
		return nil, 0, err
	}

	switch manifest.Type {
	case config.NamespaceManifest:
		var depth int
		manifest, depth, err = l.LoadManifest(ctx, manifest.Plugins, args[leafDepth:])
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
			_, err = l.visit(ctx, plugin)
			if err != nil {
				return nil, 0, err
			}

			manifest.Plugins = append(manifest.Plugins, plugin)
		}
	}

	return manifest, leafDepth, nil
}

func (l *Loader) visit(ctx context.Context, plugin config.Plugin) (*config.Manifest, error) {
	if plugin.Manifest == "" {
		return &config.Manifest{
			Usage: plugin.Usage,
			Type:  config.NamespaceManifest,
		}, nil
	}
	// fmt.Printf("visit %s\n", plugin.Name)

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
