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

func GetPlugin(ctx context.Context, plugin config.Plugin, args []string) (*Plugin, error) {
	leaf, depth, err := FindPlugin(ctx, plugin, args)
	if err != nil {
		return nil, err
	}

	manifest, err := GetManifest(ctx, leaf)
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

		cmd, err := kit.Connect(path)
		if err != nil {
			return nil, err
		}

		// Override usage with user-defined usage if available.
		usage := manifest.Usage
		if leaf.Usage != "" {
			usage = leaf.Usage
		}

		resolved := &Plugin{
			CommandPath: args[:depth],
			Usage:       usage,
			Args:        manifest.Args,
			Flags:       manifest.Flags,
			Action: func(ctx context.Context) error {
				defer cmd.Close()
				return cmd.Run(ctx)
			},
		}

		resolved.Verify = VerifyCommand(resolved, cmd, args[depth:])
		resolved.Autocomplete = AutocompleteCommand(cmd, args, depth)
		return resolved, nil
	case config.ManifestNamespace:
		var plugins []*Plugin
		for _, subplugin := range leaf.Plugins {
			submanifest, err := GetManifest(ctx, subplugin)
			if err != nil {
				return nil, err
			}

			names := make([]string, len(args[:depth])+1)
			copy(names, args[:depth])
			names[len(names)-1] = subplugin.Name

			plugins = append(plugins, &Plugin{
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

		resolved := &Plugin{
			CommandPath: commandPath,
			Usage:       namespaceUsage(plugin, manifest),
			Plugins:     plugins,
		}

		resolved.Verify = VerifyNamespace(resolved, args, depth)
		resolved.Autocomplete = AutocompleteNamespace(resolved, args, depth)
		return resolved, nil
	default:
		return nil, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
	}
}

func FindPlugin(ctx context.Context, plugin config.Plugin, args []string) (config.Plugin, int, error) {
	return findPlugin(ctx, plugin, args, 0)
}

func findPlugin(ctx context.Context, plugin config.Plugin, args []string, depth int) (config.Plugin, int, error) {
	// Get the manifest's plugins if any.
	manifest, err := GetManifest(ctx, plugin)
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
	manifest, err = GetManifest(ctx, child)
	if err != nil {
		return child, depth, err
	}

	switch manifest.Type {
	case config.ManifestExternal, config.ManifestCommand:
		// If it's a command, the rest of args are for the command.
		return child, depth, nil
	case config.ManifestNamespace:
		// If it's a namespace, there are more possible matches.
		return findPlugin(ctx, child, args, depth)
	default:
		return child, 0, fmt.Errorf("unrecognized manifest type '%s'", manifest.Type)
	}
}

func GetManifest(ctx context.Context, plugin config.Plugin) (config.Manifest, error) {
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
		panic(err)
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
