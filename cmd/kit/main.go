package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/hinshun/kit/cli"
	"github.com/hinshun/kit/system/interrupt"
	"github.com/hinshun/kit/system/profile"
)

var (
	EnvEnableProfiling = "KIT_PROF"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "kit: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		enableProfiling = os.Getenv(EnvEnableProfiling) != ""
	)

	ctx, cancel := context.WithCancel(context.Background())

	ih := interrupt.NewInterruptHandler(cancel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	defer ih.Close()

	if enableProfiling {
		p, err := profile.NewProfile()
		if err != nil {
			return err
		}
		defer p.Close()
	}

	return cli.NewKit().Run(ctx, os.Args)
}
