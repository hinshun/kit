package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	InitConfig = Config{
		Manifest: "/kit/initial",
	}

	BootstrapConfig = Config{
		Manifest: "/kit/bootstrap",
	}
)

type Config struct {
	Manifest string  `json:"manifest,omitempty"`
	Plugins  Plugins `json:"plugins"`
}

func New(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	cfg := InitConfig
	if !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &cfg)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}
