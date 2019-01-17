package cli

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hinshun/kit/config"
	"github.com/hinshun/kitapi/kit"
)

type Cli struct {
	NamespacePath  []string
	NamespaceUsage string
	Commands       []*Command
	UsageError     error

	flagSet      *flag.FlagSet
	help         *bool
	autocomplete *string
	configPath   *string

	loader *Loader
	stdio  kit.Stdio

	theme *Theme
}

func NewCli(loader *Loader) *Cli {
	flagSet := flag.NewFlagSet("cli", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return &Cli{
		flagSet:      flagSet,
		help:         flagSet.Bool("help", false, "display this help text"),
		autocomplete: flagSet.String("autocomplete", "", "print autocomplete word list for a shell"),
		configPath:   flagSet.String("config", filepath.Join(os.Getenv("HOME"), kit.ConfigPath), "path to kit config"),
		loader:       loader,
		stdio: kit.Stdio{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		theme: NewDefaultTheme(),
	}
}

func (c *Cli) GetManifest(ctx context.Context, plugin config.Plugin) (config.Manifest, error) {
	return c.loader.GetManifest(ctx, plugin)
}

func (c *Cli) ConfigPath() string {
	return *c.configPath
}

func (c *Cli) Parse(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *Cli) Args() []string {
	return append([]string{"kit"}, c.flagSet.Args()...)
}

func (c *Cli) GetCommand(ctx context.Context, plugin config.Plugin, args []string) (*Command, error) {
	return c.loader.GetCommand(ctx, plugin, args)
}

func (c *Cli) SetNamespaceUsage(commandPath []string, usage string) {
	if commandPath[0] != "kit" {
		commandPath = append([]string{"kit"}, commandPath...)
	}
	c.NamespacePath = commandPath
	c.NamespaceUsage = usage
}
