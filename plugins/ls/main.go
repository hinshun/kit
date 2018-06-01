package main

import (
	"os"
	"os/exec"
)

type command struct{}

func (c command) Name() string {
	return "ls"
}

func (c command) Usage() string {
	return "list directory contents"
}

func (c command) Action() error {
	cmd := exec.Command("ls", "-la")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

var Command command
