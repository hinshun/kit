package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/hinshun/kit/cmd/kit/command"
	"github.com/hinshun/kit/interrupt"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "kit: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	ih, ctx := interrupt.NewInterruptHandler(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer ih.Close()

	app, err := command.App(ctx)
	if err != nil {
		return err
	}

	return app.Run(os.Args)
}
