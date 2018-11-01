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
		// The root manifest is the config's manifest.
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
		cfg.Manifest = plugin.Manifest
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

	// Get the manifest's plugins if any.
	manifest, err := kit.Kit(ctx).GetManifest(ctx, plugin)
	if err != nil {
		return plugin, err
	}
	plugin.Plugins = manifest.Plugins
	if plugin.Usage == "" {
		plugin.Usage = manifest.Usage
	}

	for i, child := range plugin.Plugins {
		if child.Name != names[0] {
			continue
		}

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

		if plugin.Manifest != "" {
			// If plugins are from its namespace manifest, the user must pin the
			// namespace or use "--pin" before it can add the plugin.
			if !c.pin {
				return plugin, fmt.Errorf("requires pinning")
			}
		}

		plugin.Manifest = ""
		plugin.Plugins[i] = replace
		return plugin, nil
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

	// If plugins are from its namespace manifest, the user must pin the
	// namespace or use "--pin" before it can add the plugin.
	if plugin.Manifest != "" {
		if !c.pin {
			return plugin, fmt.Errorf("requires pinning")
		}
	}

	plugin.Manifest = ""
	plugin.Plugins = append(plugin.Plugins, child)
	plugin.Plugins.Sort()

	return plugin, nil
}
