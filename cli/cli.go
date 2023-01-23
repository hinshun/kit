package cli

import (
	"context"
	"flag"
	"io/ioutil"
	"os"

	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/introspect"
)

type Cli struct {
	NamespacePath  []string
	NamespaceUsage string
	Plugins       []*Plugin
	UsageError     error

	flagSet *flag.FlagSet

	loader *Loader
	stdio  Stdio

	options *introspect.Options
	theme   *introspect.Theme
}

func NewCli(loader *Loader) *Cli {
	flagSet := flag.NewFlagSet("cli", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return &Cli{
		flagSet: flagSet,
		loader:  loader,
		stdio: Stdio{
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

func (c *Cli) Flags() []Flag {
	return []Flag{
		BoolFlag("help", "Displays help text.", false, &c.options.Help),
		StringFlag("autocomplete", "Prints autocomplete word list for a shell", "", &c.options.Autocomplete),
	}
}

func (c *Cli) GetManifest(ctx context.Context, plugin config.Plugin) (config.Manifest, error) {
	return c.loader.GetManifest(ctx, plugin)
}

func (c *Cli) Parse(args []string) error {
	for _, flag := range c.Flags() {
		flag.Set(c.flagSet)
	}
	return c.flagSet.Parse(args)
}

func (c *Cli) Args() []string {
	return append([]string{"kit"}, c.flagSet.Args()...)
}

func (c *Cli) GetPlugin(ctx context.Context, plugin config.Plugin, args []string) (*Plugin, error) {
	return c.loader.GetPlugin(ctx, plugin, args)
}

func (c *Cli) SetNamespaceUsage(commandPath []string, usage string) {
	if commandPath[0] != "kit" {
		commandPath = append([]string{"kit"}, commandPath...)
	}
	c.NamespacePath = commandPath
	c.NamespaceUsage = usage
}
