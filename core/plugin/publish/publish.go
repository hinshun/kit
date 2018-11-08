package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/publish"
)

type command struct {
	paths string
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
	return nil
}

func (c *command) Run(ctx context.Context) error {
	pluginPaths := strings.Split(c.paths, ",")
	digest, err := publish.Publish(pluginPaths)
	if err != nil {
		return err
	}

	fmt.Println(digest)
	return nil
}
