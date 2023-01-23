package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
)

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) GetCommand(ctx context.Context, plugin config.Plugin, args []string) (*Command, error) {
	leaf, depth, err := l.FindPlugin(ctx, plugin, args)
	if err != nil {
		return nil, err
	}

	manifest, err := l.GetManifest(ctx, leaf)
	if err != nil {
		return nil, err
	}

	switch manifest.Type {
	case config.ManifestExternal:
		path, err := manifest.MatchPlatform(runtime.GOOS, runtime.GOARCH)
		if err != nil {
			return nil, err
		}
		return nil, syscall.Exec(path, args[depth:], os.Environ())
	case config.ManifestCommand:
		path, err := manifest.MatchPlatform(runtime.GOOS, runtime.GOARCH)
		if err != nil {
			return nil, err
		}

		constructor, err := kit.OpenConstructor(path)
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
			Action: func(ctx context.Context) error {
				return kitCmd.Run(ctx)
			},
		}

		cliCmd.Verify = VerifyCommand(cliCmd, kitCmd, args[depth:])
		cliCmd.Autocomplete = AutocompleteCommand(kitCmd, args, depth)
		return cliCmd, nil
	case config.ManifestNamespace:
		var commands []*Command
		for _, subplugin := range leaf.Plugins {
			submanifest, err := l.GetManifest(ctx, subplugin)
			if err != nil {
				return nil, err
			}

			names := make([]string, len(args[:depth])+1)
			copy(names, args[:depth])
			names[len(names)-1] = subplugin.Name

			commands = append(commands, &Command{
				CommandPath: names,
				Usage:       namespaceUsage(subplugin, submanifest),
				Args:        submanifest.Args,
				Flags:       submanifest.Flags,
			})
		}

		// For usage errors in the root namespace, set the root plugin's name.
		commandPath := args[:depth]
		if len(commandPath) == 0 {
			commandPath = []string{plugin.Name}
		}

		cliCmd := &Command{
			CommandPath: commandPath,
			Usage:       namespaceUsage(plugin, manifest),
			Commands:    commands,
		}

		cliCmd.Verify = VerifyNamespace(cliCmd, args, depth)
		cliCmd.Autocomplete = AutocompleteNamespace(cliCmd, args, depth)
		return cliCmd, nil
	default:
		return nil, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
	}
}

func (l *Loader) FindPlugin(ctx context.Context, plugin config.Plugin, args []string) (config.Plugin, int, error) {
	return l.findPlugin(ctx, plugin, args, 0)
}

func (l *Loader) findPlugin(ctx context.Context, plugin config.Plugin, args []string, depth int) (config.Plugin, int, error) {
	// Get the manifest's plugins if any.
	manifest, err := l.GetManifest(ctx, plugin)
	if err != nil {
		return plugin, depth, err
	}
	plugin.Plugins = manifest.Plugins

	// If there are no more args, then we found our plugin without extra args.
	if len(args) == depth {
		return plugin, depth, nil
	}

	// Find the plugin with its name matching the next arg.
	index := -1
	for i, p := range plugin.Plugins {
		if args[depth] == p.Name {
			index = i
			break
		}
	}

	// No plugin matched with the next arg, then the rest of args are for the
	// command.
	if index == -1 {
		return plugin, depth, nil
	}

	// We matched one level deeper in args.
	depth++

	// Check the type of the plugin to decide whether to continue finding.
	child := plugin.Plugins[index]
	manifest, err = l.GetManifest(ctx, child)
	if err != nil {
		return child, depth, err
	}

	switch manifest.Type {
	case config.ManifestExternal, config.ManifestCommand:
		// If it's a command, the rest of args are for the command.
		return child, depth, nil
	case config.ManifestNamespace:
		// If it's a namespace, there are more possible matches.
		return l.findPlugin(ctx, child, args, depth)
	default:
		return child, 0, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
	}
}

func (l *Loader) GetManifest(ctx context.Context, plugin config.Plugin) (config.Manifest, error) {
	if plugin.Path == "" {
		return config.Manifest{
			Usage:   plugin.Usage,
			Type:    config.ManifestNamespace,
			Plugins: plugin.Plugins,
		}, nil
	}

	f, err := os.Open(plugin.Path)
	if err != nil {
		return config.Manifest{}, err
	}
	defer f.Close()

	var manifest config.Manifest
	err = json.NewDecoder(f).Decode(&manifest)
	if err != nil {
		if errors.Is(err, &json.SyntaxError{}) {
			return config.Manifest{}, err
		}
		manifest.Type = config.ManifestExternal
		manifest.Platforms = append(manifest.Platforms, config.Platform{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
			Path: plugin.Path,
		})
	}
	return manifest, nil
}

func namespaceUsage(plugin config.Plugin, manifest config.Manifest) string {
	usage := manifest.Usage
	if plugin.Usage != "" {
		usage = plugin.Usage
	}
	if usage == "" {
		usage = "A plugin namespace."
	}
	return usage
}
