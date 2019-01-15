package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hinshun/kit/api"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kitapi/kit"
)

type command struct{}

func (c *command) Usage() string {
	return "Initializes a kit config."
}

func (c *command) Args() []kit.Arg {
	return nil
}

func (c *command) Flags() []kit.Flag {
	return nil
}

func (c *command) Run(ctx context.Context) error {
	data, err := json.MarshalIndent(&config.BootstrapConfig, "", "    ")
	if err != nil {
		return err
	}

	configPath := api.Kit(ctx).ConfigPath()
	err = os.MkdirAll(filepath.Dir(configPath), 0775)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, 0664)
}
