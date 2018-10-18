package main

import (
	"context"
	"os"
	"os/exec"

	"github.com/hinshun/kit"
)

type command struct{}

var Constructor kit.Constructor = func(cli kit.Cli) (kit.Command, error) {
	return &command{}
}

func (c command) Name() string {
	return "ls"
}

func (c command) Usage() string {
	return "list directory contents"
}

func (c command) Args() []kit.Args {
	return nil
}

func (c command) Flags() []kit.Flags {
	return nil
}

func (c command) Run(ctx context.Context) error {
	cmd := exec.Command("ls", "-la")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
