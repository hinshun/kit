package kit

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Commands  []string
	Bootstrap []string
}

func ParseConfig(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config
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
