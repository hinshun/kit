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
	commandPath string
	pluginPath  string
	usage       string
	pin         bool
	overwrite   bool
}

func (c *command) Usage() string {
	return "Adds a plugin to kit."
}

func (c *command) Args() []kit.Arg {
	return []kit.Arg{
		kit.StringArg(
			"command path",
			"The command path to add the plugin.",
			&c.commandPath,
		),
		kit.StringArg(
			"plugin path",
			"The path on disk for an executable or a plugin manifest.",
			&c.pluginPath,
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
			&c.pin,
		),
		kit.BoolFlag(
			"overwrite",
			"Overwrites any namespace or command if conflicting at the command path.",
			&c.overwrite,
		),
	}
}

func (c *command) Run(ctx context.Context) error {
	configPath := filepath.Join(os.Getenv("HOME"), kit.ConfigPath)
	cfg, err := config.New(configPath)
	if err != nil {
		return err
	}

	if c.commandPath == "/" {
		// The root manifest is the config's manifest.
		cfg.Path = c.pluginPath
	} else {
		c.commandPath = strings.Trim(c.commandPath, "/")
		names := strings.Split(c.commandPath, "/")
		plugin, err := c.addPlugin(ctx, names, config.Plugin{
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

// addPlugin returns a modified plugin with a new plugin with manifest
// c.pluginPath and usage c.usage added by path in names. If the namespaces
// don't exist along the path, they are created on the fly.
func (c *command) addPlugin(ctx context.Context, names []string, plugin config.Plugin) (config.Plugin, error) {
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

		var replace config.Plugin

		// If names is length 1, there is a conflict. The user can remove the
		// conflicting plugin and re-add, or use "--overwrite" when adding.
		if len(names) == 1 {
			if !c.overwrite {
				return plugin, fmt.Errorf("conflict")
			}

			replace = config.Plugin{
				Name:  child.Name,
				Path:  c.pluginPath,
				Usage: c.usage,
			}
		} else {
			replace, err = c.addPlugin(ctx, names[1:], child)
			if err != nil {
				return plugin, err
			}
		}

		if plugin.Path != "" {
			// If plugins are from its namespace manifest, the user must pin the
			// namespace or use "--pin" before it can add the plugin.
			if !c.pin {
				return plugin, fmt.Errorf("requires pinning")
			}
		}

		plugin.Path = ""
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
		child.Path = c.pluginPath
		child.Usage = c.usage
	}

	// If plugins are from its namespace manifest, the user must pin the
	// namespace or use "--pin" before it can add the plugin.
	if plugin.Path != "" {
		if !c.pin {
			return plugin, fmt.Errorf("requires pinning")
		}
	}

	plugin.Path = ""
	plugin.Plugins = append(plugin.Plugins, child)
	plugin.Plugins.Sort()

	return plugin, nil
}
