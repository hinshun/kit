package main

import (
	"fmt"
	"os"

	"github.com/hinshun/kit/cmd/kit/command"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "kit: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	app, err := command.App()
	if err != nil {
		return err
	}

	return app.Run(os.Args)
}
