package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/ipfs"
	"github.com/hinshun/kit/linker"
	"github.com/ipfs/go-ipfs/core/coreapi"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "publish: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	n, err := ipfs.NewNode(ctx, nil)
	if err != nil {
		return err
	}
	api := coreapi.NewCoreAPI(n)

	var commands []string
	err = filepath.Walk("./bin", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		command, err := linker.LinkCommand(path)
		if err != nil {
			return err
		}

		p, err := api.Block().Put(ctx, f)
		if err != nil {
			return err
		}

		fmt.Printf("published plugin '%s' to ipfs '%s'\n", command.Name(), p.Cid())
		commands = append(commands, p.Cid().String())
		return nil
	})
	if err != nil {
		return err
	}

	config := kit.Config{
		Commands: commands,
	}
	data, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("kit.json", data, 0644)
}
