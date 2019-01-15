package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hinshun/kit/api"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/core/plugin"
	"github.com/hinshun/kitapi/kit"
)

type command struct {
	path string
	pin  bool
}

func (c *command) Usage() string {
	return "Removes a plugin from kit."
}

func (c *command) Args() []kit.Arg {
	return []kit.Arg{
		plugin.CommandPathArg(
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
	configPath := api.Kit(ctx).ConfigPath()
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

		cfg.Manifest = plugin.Manifest
		cfg.Plugins = plugin.Plugins
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

	// Get the manifest's plugins if any.
	manifest, err := api.Kit(ctx).GetManifest(ctx, plugin)
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

		// If plugins are from its namespace manifest, the user must pin the
		// namespace or use "--pin" before it can remove the plugin.
		if plugin.Manifest != "" {
			if !c.pin {
				return plugin, fmt.Errorf("requires pinning")
			}
			plugin.Manifest = ""
		}

		if len(names) > 1 {
			// Matched the command path, recurse to find the plugin to remove.
			replace, err := c.removePlugin(ctx, names[1:], child)
			if err != nil {
				return plugin, err
			}

			plugin.Plugins[i] = replace
			return plugin, nil
		} else {
			plugin.Plugins = append(plugin.Plugins[:i], plugin.Plugins[i+1:]...)
			return plugin, nil
		}
	}

	// No plugin found following the command path. Return as an idempotent delete.
	return plugin, nil
}
