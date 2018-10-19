package cli

import (
	"os"

	"github.com/hinshun/kit"
)

type Cli struct {
	Commands []*Command

	stdio kit.Stdio
}

func NewCli() *Cli {
	return &Cli{
		stdio: kit.Stdio{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
	}
}

func (c *Cli) Stdio() kit.Stdio {
	return c.stdio
}

func (c *Cli) Args() kit.Args {
	return nil
}

func (c *Cli) Flags() kit.Flags {
	return nil
}
