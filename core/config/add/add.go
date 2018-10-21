package main

import (
	"context"
	"fmt"

	"github.com/hinshun/kit"
)

type command struct {
	path     string
	manifest string
}

var New kit.Constructor = func() (kit.Command, error) {
	return &command{}, nil
}

func (c command) Usage() string {
	return "Adds a plugin to kit's config."
}

func (c command) Args() []kit.Arg {
	return []kit.Arg{
		kit.CommandPathArg(
			&c.path,
			"The command path to add the plugin.",
		),
		kit.ManifestArg(
			&c.manifest,
			"New plugin's manifest. An empty string will create an empty namespace.",
		),
	}
}

func (c command) Flags() []kit.Flag {
	return nil
}

func (c command) Run(ctx context.Context) error {
	fmt.Printf("Path: %s, Manifest: %s\n", c.path, c.manifest)
	return nil
}
