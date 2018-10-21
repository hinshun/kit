package cli

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/hinshun/kit"
)

type Cli struct {
	Commands   []*Command
	UsageError error

	flagSet    *flag.FlagSet
	configPath *string
	help       *bool

	stdio kit.Stdio

	headerColor, usageErrorColor, commandColor, argColor, flagColor *color.Color
}

func NewCli() *Cli {
	flagSet := flag.NewFlagSet("cli", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return &Cli{
		flagSet:    flagSet,
		configPath: flagSet.String("config", filepath.Join(os.Getenv("HOME"), kit.ConfigPath), "path to kit config"),
		help:       flagSet.Bool("help", false, "display this help text"),
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

func (c *Cli) Parse(args []string) error {
	err := c.flagSet.Parse(args)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cli) ConfigPath() string {
	return *c.configPath
}
