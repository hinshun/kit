package cli

import "github.com/hinshun/kit"

type cli struct {
}

func NewCli() kit.Cli {
	return &cli{}
}

func (c *cli) Stdio() kit.Stdio {
	return kit.Stdio{}
}

func (c *cli) Args() kit.Args {
	return nil
}

func (c *cli) Flags() kit.Flags {
	return nil
}
