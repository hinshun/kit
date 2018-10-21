package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		cfg.Plugins, err = AddPlugin(names, cfg.Plugins, c.manifest, c.usage)
		if err != nil {
			return err
		}
	}

	cfgJson, err := json.MarshalIndent(&cfg, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, cfgJson, 0664)
}

func AddPlugin(names []string, plugins config.Plugins, manifest, usage string) (config.Plugins, error) {
	if len(names) == 0 {
		return nil, nil
	}

	var err error
	for i, plugin := range plugins {
		if plugin.Name == names[0] {
			if len(names) == 1 {
				return nil, fmt.Errorf("conflict")
			}

			plugins[i].Plugins, err = AddPlugin(names[1:], plugin.Plugins, manifest, usage)
			if err != nil {
				return nil, err
			}
			return plugins, nil
		}
	}

	plugin := config.Plugin{
		Name: names[0],
	}

	if len(names) == 1 {
		plugin.Manifest = manifest
		plugin.Usage = usage
	} else {
		plugin.Plugins, err = AddPlugin(names[1:], nil, manifest, usage)
		if err != nil {
			return nil, err
		}
	}

	plugins = append(plugins, plugin)

	// Lexicographically sort plugins by name to produce a deterministic config.
	sort.SliceStable(plugins, func(i, j int) bool {
		return plugins[i].Name < plugins[j].Name
	})

	return plugins, nil
}
