package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/hinshun/kit"
)

type Cli struct {
	Commands   []*Command
	UsageError error

	stdio      kit.Stdio
	configPath string

	headerColor, usageErrorColor, commandColor, argColor, flagColor *color.Color
}

func NewCli() *Cli {
	return &Cli{
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

func (c *Cli) ConfigPath() string {
	return c.configPath
}
