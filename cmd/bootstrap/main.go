package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/linker"
	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := &kit.Config{}

	sh := shell.NewLocalShell()
	err := filepath.Walk("./bin", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		hash, err := sh.Add(f)
		if err != nil {
			return err
		}

		command, err := linker.LinkCommand(path)
		if err != nil {
			return err
		}

		fmt.Printf("published plugin '%s' to ipfs '%s'\n", command.Name(), hash)
		cfg.Plugins = append(cfg.Plugins, kit.Plugin{
			Name: command.Name(),
			Type: kit.PluginCommand,
			Ref:  hash,
		})
		return nil
	})
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(".kit/config.json", data, 0644)
}
