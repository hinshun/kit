package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var (
	InitConfig = Config{
		Plugins: Plugins{{
			Name: "si",
			Usage: "Tools for Software Infra team",
			Plugins: Plugins{{
				Name: "break-glass",
				Usage: "Force approve a pull request and emails change management committee",
				Path: "/n/nix/tech/store/vcddxin45qdxqddl97nwipav56x1jzzg-si-scripts/bin/break-glass",
			}},
		}},
	}
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
