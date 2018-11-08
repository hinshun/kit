package main

import (
	"fmt"
	"os"

	"github.com/hinshun/kit/publish"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "publish: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	return publish.Publish(os.Args[1:])
}
