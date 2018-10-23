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
	path string
	pin  bool
}

var New kit.Constructor = func() (kit.Command, error) {
	return &command{}, nil
}

func (c *command) Usage() string {
	return "Removes a plugin from kit."
}

func (c *command) Args() []kit.Arg {
	return []kit.Arg{
		kit.CommandPathArg(
			"The command path to remove the plugin.",
			&c.path,
		),
	}
}

func (c *command) Flags() []kit.Flag {
	return []kit.Flag{
		kit.BoolFlag(
			"pin",
			"Pins the plugin's parent namespace if removing from an implicit namespace.",
			false,
			&c.pin,
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
		cfg.Manifest = ""
		cfg.Plugins = nil
	} else {
		c.path = strings.Trim(c.path, "/")
		names := strings.Split(c.path, "/")
		plugin, err := c.removePlugin(ctx, names, config.Plugin{
			Manifest: cfg.Manifest,
			Plugins:  cfg.Plugins,
		})
		if err != nil {
			return err
		}

		cfg.Plugins = plugin.Plugins

		// If removing a plugin from the root namespace and the namespace manifest
		// was unpinned then unpin the config manifest too.
		if len(names) == 1 {
			cfg.Manifest = plugin.Manifest
		}
	}

	cfgJson, err := json.MarshalIndent(&cfg, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, cfgJson, 0664)
}

func (c *command) removePlugin(ctx context.Context, names []string, plugin config.Plugin) (config.Plugin, error) {
	if len(names) == 0 {
		return plugin, nil
	}

	// Load the base manifest and merge with user defined plugins.
	manifest, err := kit.Kit(ctx).GetManifest(ctx, plugin)
	if err != nil {
		return plugin, err
	}
	merged := manifest.Plugins.Merge(plugin.Plugins)

	for i, child := range merged {
		if child.Name == names[0] {
			// Determine whether child came from config or a namespace manifest.
			configIndex := -1
			for j, configChild := range plugin.Plugins {
				if configChild.Name == child.Name {
					configIndex = j
					break
				}
			}

			if len(names) > 1 {
				// Matched the command path, recurse to find the plugin to remove.
				replace, err := c.removePlugin(ctx, names[1:], child)
				if err != nil {
					return plugin, err
				}

				// If no child is found from the config's plugins, then it was merged from a
				// namespace manifest. The user can pin the implict namespace or use "--pin"
				// when removing.
				if configIndex == -1 {
					plugin.Plugins = merged
					plugin.Plugins[i] = replace
				} else {
					plugin.Plugins[configIndex] = replace
				}

				return plugin, nil
			} else {
				// If no child is found from the config's plugins, then it was merged from a
				// namespace manifest. The user can pin the plugin's siblings or use "--pin"
				// when removing.
				if configIndex == -1 {
					if !c.pin {
						return plugin, fmt.Errorf("requires pinning")
					}

					plugin.Manifest = ""
					plugin.Usage = manifest.Usage
					plugin.Plugins = merged
					plugin.Plugins = append(plugin.Plugins[:i], plugin.Plugins[i+1:]...)
					return plugin, nil
				}

				manifestIndex := -1
				for j, manifestChild := range manifest.Plugins {
					if manifestChild.Name == child.Name {
						manifestIndex = j
						break
					}
				}

				if manifestIndex != -1 {
					plugin.Manifest = ""
					plugin.Usage = manifest.Usage
				}

				plugin.Plugins = append(plugin.Plugins[:configIndex], plugin.Plugins[configIndex+1:]...)
				return plugin, nil
			}
		}
	}

	// No plugin found following the command path. Return as an idempotent delete.
	return plugin, nil
}
