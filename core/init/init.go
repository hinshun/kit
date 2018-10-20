package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
)

type command struct{}

var New kit.Constructor = func(c kit.Cli) (kit.Command, error) {
	return &command{}, nil
}

func (c command) Usage() string {
	return "Initializes a kit config."
}

func (c command) Args() []kit.Arg {
	return nil
}

func (c command) Flags() []kit.Flag {
	return nil
}

func (c command) Run(ctx context.Context) error {
	data, err := json.MarshalIndent(&config.BootstrapConfig, "", "    ")
	if err != nil {
		return err
	}

	kitDir := filepath.Join(os.Getenv("HOME"), ".kit")
	err = os.MkdirAll(kitDir, 0775)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(kitDir, "config.json"), data, 0664)
}
