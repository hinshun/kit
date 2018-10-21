package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
)

type command struct {
	path     string
	manifest string
	usage    string
}

var New kit.Constructor = func() (kit.Command, error) {
	return &command{}, nil
}

func (c *command) Usage() string {
	return "Adds a plugin to kit's config."
}

func (c *command) Args() []kit.Arg {
	return []kit.Arg{
		kit.CommandPathArg(
			"The command path to add the plugin.",
			&c.path,
		),
		kit.ManifestArg(
			"New plugin's manifest. An empty string will create an empty namespace.",
			&c.manifest,
		),
	}
}

func (c *command) Flags() []kit.Flag {
	return []kit.Flag{
		kit.StringFlag(
			"usage",
			"Usage help text for the plugin.",
			"",
			&c.usage,
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
		cfg.Manifest = c.manifest
	} else {
		c.path = strings.Trim(c.path, "/")
		names := strings.Split(c.path, "/")
		plugin, err := AddPlugin(ctx, names, c.manifest, c.usage, config.Plugin{
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

func AddPlugin(ctx context.Context, names []string, manifest, usage string, plugin config.Plugin) (config.Plugin, error) {
	log.Printf("add plugin '%s' with names '%s'", plugin.Name, names)
	if len(names) == 0 {
		return plugin, nil
	}

	// Load the base manifest and merge with user defined plugins.
	msft, err := kit.Kit(ctx).GetManifest(ctx, plugin)
	if err != nil {
		return plugin, err
	}
	plugin.Plugins = msft.Plugins.Merge(plugin.Plugins)

	for i, child := range plugin.Plugins {
		if child.Name == names[0] {
			if len(names) == 1 {
				return plugin, fmt.Errorf("conflict")
			}

			log.Println("diving in")
			replace, err := AddPlugin(ctx, names[1:], manifest, usage, child)
			if err != nil {
				return plugin, err
			}

			plugin.Plugins[i] = replace
			return plugin, nil
		}
	}

	log.Println("adding new namespace")
	child, err := AddPlugin(ctx, names[1:], manifest, usage, config.Plugin{
		Name: names[0],
	})
	if err != nil {
		return plugin, err
	}

	// Arrived at the leaf node, which is the plugin you are intending to add.
	if len(child.Plugins) == 0 {
		child.Manifest = manifest
		child.Usage = usage
	}

	log.Printf("after new namespace plugin '%s'", plugin.Name)
	plugin.Plugins = append(plugin.Plugins, child)

	// Lexicographically sort plugins by name to produce a deterministic config.
	sort.SliceStable(plugin.Plugins, func(i, j int) bool {
		return plugin.Plugins[i].Name < plugin.Plugins[j].Name
	})

	return plugin, nil
}
