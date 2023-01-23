package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	InitConfig = Config{}
)

// Config
type Config struct {
	Path    string  `json:"path,omitempty"`
	Plugins Plugins `json:"plugins,omitempty"`
}

func New(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config
	if os.IsNotExist(err) {
		cfg = InitConfig
	} else {
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
