package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/introspect"
)

type Cli struct {
	NamespacePath  []string
	NamespaceUsage string
	Plugins        []*Plugin
	UsageError     error

	parsedArgs []string

	stdio kit.Stdio

	options *introspect.Options
	theme   *introspect.Theme
}

func New() *Cli {
	return &Cli{
		stdio: kit.Stdio{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		options: &introspect.Options{},
		theme:   introspect.NewDefaultTheme(),
	}
}

func (c *Cli) Options() introspect.Options {
	return *c.options
}

func (c *Cli) Theme() introspect.Theme {
	return *c.theme
}

func (c *Cli) Flags() []kit.Flag {
	return []kit.Flag{
		kit.BoolFlag("help", "Displays help text.", &c.options.Help),
		kit.StringFlag("autocomplete", "Prints autocomplete word list for a shell.", "", &c.options.Autocomplete),
	}
}

func (c *Cli) Parse(ctx context.Context, shellArgs []string) error {
	var err error
	c.parsedArgs, err = parse(ctx, c.Flags(), shellArgs)
	return err
}

func (c *Cli) Args() []string {
	return append([]string{"kit"}, c.parsedArgs...)
}

func (c *Cli) SetNamespaceUsage(commandPath []string, usage string) {
	if commandPath[0] != "kit" {
		commandPath = append([]string{"kit"}, commandPath...)
	}
	c.NamespacePath = commandPath
	c.NamespaceUsage = usage
}

func (c *Cli) Run(ctx context.Context, shellArgs []string) error {
	configPath := filepath.Join(os.Getenv("HOME"), kit.ConfigPath)
	cfg, err := config.New(configPath)
	if err != nil {
		return err
	}

	plugin := config.Plugin{
		Plugins: config.Plugins{
			{
				Name:    "kit",
				Usage:   "Composable command-line toolkit.",
				Plugins: cfg.Plugins,
			},
		},
	}

	manifest, err := GetManifest(ctx, plugin)
	if err != nil {
		return err
	}

	merged := manifest.Plugins.Merge(plugin.Plugins)
	if len(merged) == 0 {
		plugin.Plugins = config.InitConfig.Plugins
	}

	err = c.Parse(ctx, shellArgs)
	if err != nil {
		return err
	}

	resolved, err := GetPlugin(ctx, plugin, c.Args())
	if err != nil {
		return err
	}

	if c.options.Autocomplete != "" {
		// ctx = introspect.WithKit(ctx, c)
		completions, err := resolved.Autocomplete(ctx, c.options.Autocomplete)
		if err != nil {
			return err
		}

		switch c.options.Autocomplete {
		case "bash", "fish":
			var wordlist []string
			for _, completion := range completions {
				wordlist = append(wordlist, completion.Wordlist...)
			}
			fmt.Printf("%s", strings.Join(wordlist, " "))
		case "zsh":
			var shellCmds []string
			for _, completion := range completions {
				shellCmds = append(
					shellCmds,
					fmt.Sprintf("local -a %s", completion.Group),
					fmt.Sprintf("%s=(%s)", completion.Group, strings.Join(completion.Wordlist, " ")),
					fmt.Sprintf("_describe %s %s", completion.Group, completion.Group),
				)
			}
			fmt.Printf("%s", strings.Join(shellCmds, ";"))
		}
		return nil
	}

	if !c.options.Help {
		err = resolved.Verify(ctx, c)
		if err != nil {
			c.UsageError = err
		}
	}

	if c.options.Help || c.UsageError != nil || resolved.Action == nil {
		if resolved.Action == nil {
			c.SetNamespaceUsage(resolved.CommandPath, resolved.Usage)
			return c.printHelp(resolved.Plugins)
		} else {
			namespace := config.Plugin{Plugins: merged}.FindParent(resolved.CommandPath)
			namespaceManifest, err := GetManifest(ctx, namespace)
			if err != nil {
				return err
			}

			c.SetNamespaceUsage(resolved.CommandPath[:len(resolved.CommandPath)-1], namespaceManifest.Usage)
			return c.printHelp([]*Plugin{resolved})
		}
	}

	// ctx = introspect.WithKit(ctx, c)
	return resolved.Action(ctx)
}
