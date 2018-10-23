package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
)

type command struct {
	path      string
	manifest  string
	usage     string
	pin       bool
	overwrite bool
}

var New kit.Constructor = func() (kit.Command, error) {
	return &command{}, nil
}

func (c *command) Usage() string {
	return "Adds a plugin to kit."
}

func (c *command) Args() []kit.Arg {
	return []kit.Arg{
		kit.CommandPathArg(
			"The command path to add the plugin.",
			&c.path,
		),
		kit.ManifestArg(
			"The content address or resolvable name for a plugin's metadata.",
			&c.manifest,
		),
	}
}

func (c *command) Flags() []kit.Flag {
	return []kit.Flag{
		kit.StringFlag(
			"usage",
			"Specify usage help text for the plugin.",
			"",
			&c.usage,
		),
		kit.BoolFlag(
			"pin",
			"Pins the plugin's parent namespace if adding to an implicit namespace.",
			false,
			&c.pin,
		),
		kit.BoolFlag(
			"overwrite",
			"Overwrites any namespace or command if conflicting at the command path.",
			false,
			&c.overwrite,
		),
	}
}

func (c *command) Run(ctx context.Context) error {
	configPath := kit.Kit(ctx).ConfigPath()
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	var cfg config.Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	if c.path == "/" {
		// The config's manifest is the root's manifest.
		cfg.Manifest = c.manifest
	} else {
		c.path = strings.Trim(c.path, "/")
		names := strings.Split(c.path, "/")
		plugin, err := c.addPlugin(ctx, names, config.Plugin{
			Manifest: cfg.Manifest,
			Plugins:  cfg.Plugins,
		})
		if err != nil {
			return err
		}
		cfg.Plugins = plugin.Plugins
	}

	cfgJson, err := json.MarshalIndent(&cfg, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, cfgJson, 0664)
}

// addPlugin returns a modified plugin with a new plugin with manifest
// c.manifest and usage c.usage added by path in names. If the namespaces
// don't exist along the path, they are created on the fly.
func (c *command) addPlugin(ctx context.Context, names []string, plugin config.Plugin) (config.Plugin, error) {
	if len(names) == 0 {
		return plugin, nil
	}

	// Load the base manifest and merge with config plugins.
	manifest, err := kit.Kit(ctx).GetManifest(ctx, plugin)
	if err != nil {
		return plugin, err
	}
	merged := manifest.Plugins.Merge(plugin.Plugins)

	for i, child := range merged {
		if child.Name == names[0] {
			var replace config.Plugin

			// If names is length 1, there is a conflict. The user can remove the
			// conflicting plugin and re-add, or use "--overwrite" when adding.
			if len(names) == 1 {
				if !c.overwrite {
					return plugin, fmt.Errorf("conflict")
				}

				replace = config.Plugin{
					Name:     child.Name,
					Manifest: c.manifest,
					Usage:    c.usage,
				}
			} else {
				replace, err = c.addPlugin(ctx, names[1:], child)
				if err != nil {
					return plugin, err
				}
			}

			// Determine whether child came from config or a namespace manifest.
			index := -1
			for j, configChild := range plugin.Plugins {
				if configChild.Name == child.Name {
					index = j
					break
				}
			}

			// If no child is found from the config's plugins, then it was merged from a
			// namespace manifest. The user can pin the implicit namespace or use "--pin"
			// when adding.
			if index == -1 {
				if !c.pin {
					return plugin, fmt.Errorf("requires pinning")
				}

				if len(names) > 1 {
					// Lock in the merged plugins as we're not in the immediate parent
					// namespace yet.
					plugin.Plugins = merged
					plugin.Plugins[i] = replace
				} else {
					// Only add the new plugin because we don't want to lock in the sibling
					// plugins too.
					plugin.Plugins = append(plugin.Plugins, replace)
					plugin.Plugins.Sort()
				}

				return plugin, nil
			}

			plugin.Plugins[index] = replace
			return plugin, nil
		}
	}

	// We found no plugin in this level matching the command path, so we can add a
	// new namespace along the way.
	child, err := c.addPlugin(ctx, names[1:], config.Plugin{
		Name: names[0],
	})
	if err != nil {
		return plugin, err
	}

	// Arrived at the leaf node, which is the where we want to add the new plugin.
	if len(child.Plugins) == 0 {
		child.Manifest = c.manifest
		child.Usage = c.usage
	}

	plugin.Plugins = append(plugin.Plugins, child)
	plugin.Plugins.Sort()

	return plugin, nil
}
