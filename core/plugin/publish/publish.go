package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hinshun/kit/publish"
	"github.com/hinshun/kitapi/kit"
	shell "github.com/ipfs/go-ipfs-api"
)

type command struct {
	paths string
	host  string
}

func (c *command) Usage() string {
	return "Publishes a plugin to IPFS."
}

func (c *command) Args() []kit.Arg {
	return []kit.Arg{
		kit.StringArg(
			"paths",
			"The comma delimited paths to compiled plugins of form name-GOOS-GOARCH.",
			&c.paths,
		),
	}
}

func (c *command) Flags() []kit.Flag {
	return []kit.Flag{
		kit.StringFlag(
			"host",
			"The host address of the ipfs daemon to publish to.",
			"",
			&c.host,
		),
	}
}

func (c *command) Run(ctx context.Context) error {
	var sh *shell.Shell
	if c.host == "" {
		sh = shell.NewLocalShell()
	} else {
		sh = shell.NewShell(c.host)
	}

	pluginPaths := strings.Split(c.paths, ",")
	digest, err := publish.Publish(sh, pluginPaths)
	if err != nil {
		return err
	}

	fmt.Println(digest)
	return nil
}
