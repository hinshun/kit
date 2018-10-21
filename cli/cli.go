package cli

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
)

type Cli struct {
	Commands   []*Command
	UsageError error

	flagSet    *flag.FlagSet
	help       *bool
	configPath *string

	loader *Loader
	stdio  kit.Stdio

	headerColor, usageErrorColor, commandColor, argColor, flagColor *color.Color
}

func NewCli(loader *Loader) *Cli {
	flagSet := flag.NewFlagSet("cli", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return &Cli{
		flagSet:    flagSet,
		help:       flagSet.Bool("help", false, "display this help text"),
		configPath: flagSet.String("config", filepath.Join(os.Getenv("HOME"), kit.ConfigPath), "path to kit config"),
		loader:     loader,
		stdio: kit.Stdio{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		headerColor:     color.New(color.Bold, color.Underline),
		usageErrorColor: color.New(color.FgRed, color.Bold, color.Underline),
		commandColor:    color.New(color.FgWhite, color.Underline),
		argColor:        color.New(color.FgYellow),
		flagColor:       color.New(color.FgGreen),
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

func (c *Cli) GetCommand(ctx context.Context, plugin config.Plugin, args []string) (*Command, error) {
	return c.loader.GetCommand(ctx, plugin, args)
}
