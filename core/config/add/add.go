package main

import (
	"context"
	"os"
	"os/exec"

	"github.com/hinshun/kit"
)

type command struct{}

var New kit.Constructor = func(cli kit.Cli) (kit.Command, error) {
	return &command{}, nil
}

func (c command) Usage() string {
	return "Add a plugin to kit's config"
}

func (c command) Args() []kit.Arg {
	return nil
}

func (c command) Flags() []kit.Flag {
	return nil
}

func (c command) Run(ctx context.Context) error {
	cmd := exec.Command("ls", "-la")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
