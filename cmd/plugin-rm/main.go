package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/cli"
	"github.com/hinshun/kit/config"
)

func main() {
	kit.Serve(&command{})
}

type command struct {
	path string
	pin  bool
}

func (c *command) Usage() string {
	return "Removes a plugin from kit."
}

func (c *command) Args() []kit.Arg {
	return []kit.Arg{
		kit.StringArg(
			"command path",
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
			&c.pin,
		),
	}
}

func (c *command) Run(ctx context.Context) error {
	configPath := filepath.Join(os.Getenv("HOME"), kit.ConfigPath)
	cfg, err := config.New(configPath)
	if err != nil {
		return err
	}

	if c.path == "/" {
		cfg.Path = ""
		cfg.Plugins = nil
	} else {
		c.path = strings.Trim(c.path, "/")
		names := strings.Split(c.path, "/")
		plugin, err := c.removePlugin(ctx, names, config.Plugin{
			Path:    cfg.Path,
			Plugins: cfg.Plugins,
		})
		if err != nil {
			return err
		}

		cfg.Path = plugin.Path
		cfg.Plugins = plugin.Plugins
	}

	cfgJson, err := json.MarshalIndent(&cfg, "", "    ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(configPath), 0o755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, cfgJson, 0o664)
}

func (c *command) removePlugin(ctx context.Context, names []string, plugin config.Plugin) (config.Plugin, error) {
	if len(names) == 0 {
		return plugin, nil
	}

	// Get the manifest's plugins if any.
	manifest, err := cli.GetManifest(ctx, plugin)
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
		if plugin.Path != "" {
			if !c.pin {
				return plugin, fmt.Errorf("requires pinning")
			}
			plugin.Path = ""
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
