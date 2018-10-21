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
	plugin := config.Plugin{
		Name:     "kit",
		Manifest: cfg.Manifest,
		Plugins:  cfg.Plugins,
	}

	leaf, depth, err := l.FindPlugin(ctx, plugin, args)
	if err != nil {
		return nil, err
	}

	manifest, err := l.GetManifest(ctx, leaf)
	if err != nil {
		return nil, err
	}

	switch manifest.Type {
	case config.NamespaceManifest:
		var commands []*Command
		for _, plugin := range leaf.Plugins {
			submanifest, err := l.GetManifest(ctx, plugin)
			if err != nil {
				return nil, err
			}

			names := make([]string, len(args[:depth])+1)
			copy(names, args[:depth])
			names[len(names)-1] = plugin.Name

			// Override usage with user-defined usage if available.
			usage := submanifest.Usage
			if plugin.Usage != "" {
				usage = plugin.Usage
			}

			commands = append(commands, &Command{
				CommandPath: names,
				Usage:       usage,
				Args:        submanifest.Args,
				Flags:       submanifest.Flags,
			})
		}

		// For usage errors in the root namespace of kit, we should display "kit"
		// instead of an empty path.
		commandPath := args[:depth]
		if len(commandPath) == 0 {
			commandPath = []string{"kit"}
		}

		cliCmd := &Command{
			CommandPath: commandPath,
			Action: func(ctx context.Context) error {
				return l.cli.PrintHelp(commands)
			},
		}

		err := l.cli.VerifyNamespace(cliCmd, args, depth)
		if err != nil {
			l.cli.UsageError = err
		}

		return cliCmd, nil
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

		// Override usage with user-defined usage if available.
		usage := manifest.Usage
		if leaf.Usage != "" {
			usage = leaf.Usage
		}

		cliCmd := &Command{
			CommandPath: args[:depth],
			Usage:       usage,
			Args:        manifest.Args,
			Flags:       manifest.Flags,
		}

		err = l.cli.VerifyCommand(cliCmd, kitCmd, args[depth:])
		if err != nil {
			l.cli.UsageError = err
		}

		cliCmd.Action = func(ctx context.Context) error {
			if *l.cli.help || l.cli.UsageError != nil {
				return l.cli.PrintHelp([]*Command{cliCmd})
			}
			return kitCmd.Run(ctx)
		}

		return cliCmd, nil
	}

	return nil, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
}

func (l *Loader) FindPlugin(ctx context.Context, plugin config.Plugin, args []string) (config.Plugin, int, error) {
	return l.findPlugin(ctx, plugin, args, 0)
}

func (l *Loader) findPlugin(ctx context.Context, plugin config.Plugin, args []string, depth int) (config.Plugin, int, error) {
	// Load the base manifest and merge with user defined plugins.
	manifest, err := l.GetManifest(ctx, plugin)
	if err != nil {
		return plugin, depth, err
	}

	plugin.Plugins = manifest.Plugins.Merge(plugin.Plugins)

	// If there are no more args, then we found our plugin without extra args.
	if len(args) == depth {
		return plugin, depth, nil
	}

	// Find the plugin matching the next arg.
	index := -1
	for i, p := range plugin.Plugins {
		if args[depth] == p.Name {
			index = i
			break
		}
	}

	// No plugin matched with the next arg, the rest of args are for the command.
	if index == -1 {
		return plugin, depth, nil
	}

	// We matched one level deeper in args.
	depth++

	child := plugin.Plugins[index]
	manifest, err = l.GetManifest(ctx, child)
	if err != nil {
		return child, depth, err
	}

	switch manifest.Type {
	case config.NamespaceManifest:
		plugin = config.Plugin{
			Plugins: manifest.Plugins,
		}

		// If it's a namespace, there are more possible matches.
		return l.findPlugin(ctx, child, args, depth)
	case config.CommandManifest:
		// If it's a command, the rest of args are for the command.
		return child, depth, nil
	default:
		return child, 0, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
	}
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
