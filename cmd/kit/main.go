package main

import (
	"context"
	"fmt"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/blockservice"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core/coreapi"
	coreiface "gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core/coreapi/interface"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core/coreunix"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/merkledag"
	"gx/ipfs/QmceUdzxkimdYsgtX733uNgzf1DLHyBKN6ehGSp85ayppM/go-ipfs-cmdkit/files"
	"io/ioutil"
	"os"
	"plugin"

	"github.com/BurntSushi/toml"
	"github.com/hinshun/kit/ephemeral"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

type Config struct {
	Commands  []string
	Bootstrap []string
}

func main() {
	ctx := context.Background()

	app := &cli.App{
		Name: "kit",
	}
	logrus.SetLevel(logrus.DebugLevel)

	commands, err := loadCommands(ctx)
	if err != nil {
		fmt.Printf("fatal: %s\n", err)
		os.Exit(1)
	}
	app.Commands = commands

	app.Run(os.Args)
}

func loadConfig() (*Config, error) {
	kitToml := ".kit.toml"
	_, err := os.Stat(kitToml)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config
	if !os.IsNotExist(err) {
		kitData, err := ioutil.ReadFile(kitToml)
		if err != nil {
			return nil, err
		}

		_, err = toml.Decode(string(kitData), &cfg)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func loadCommands(ctx context.Context) ([]*cli.Command, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	n, err := ephemeral.NewNode(ctx, cfg.Bootstrap)
	if err != nil {
		return nil, err
	}

	api := coreapi.NewCoreAPI(n)

	var paths []string
	for _, name := range cfg.Commands {
		path, err := api.Name().Resolve(ctx, name)
		if err != nil {
			return nil, err
		}
		logrus.Debugf("resolved %s to %s", name, path)

		// path, err := coreapi.ParsePath(name)
		// if err != nil {
		// 	return nil, err
		// }

		localPath := fmt.Sprintf("%s/.kit%s", os.Getenv("HOME"), path)
		stat, err := os.Stat(localPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		if os.IsNotExist(err) {
			logrus.Debugf("downloading %s from ipfs...", path)
			err = Get(ctx, api, path, localPath)
			if err != nil {
				return nil, err
			}
		} else {
			logrus.Debugf("found %s, verifying hash...", localPath)
			bserv := blockservice.New(n.Blockstore, n.Exchange)
			dserv := merkledag.NewDAGService(bserv)
			fileAdder, err := coreunix.NewAdder(ctx, n.Pinning, n.Blockstore, dserv)
			if err != nil {
				return nil, fmt.Errorf("1: %s", err)
			}

			outChan := make(chan interface{}, 1)
			fileAdder.Out = outChan

			r, err := os.Open(localPath)
			if err != nil {
				return nil, err
			}

			f := files.NewReaderFile(stat.Name(), localPath, r, stat)
			err = fileAdder.AddFile(f)
			if err != nil {
				return nil, err
			}

			out := <-outChan
			output := out.(*coreunix.AddedObject)

			if output.Hash != stat.Name() {
				logrus.Debugf("local hash mismatch '%s', downloading %s from ipfs...", output.Hash, path)
				err = Get(ctx, api, path, localPath)
				if err != nil {
					return nil, err
				}
			}
		}
		paths = append(paths, localPath)
	}

	var commands []*cli.Command
	for _, path := range paths {
		command, err := loadCommand(path)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}

	return commands, nil
}

func Get(ctx context.Context, api coreiface.CoreAPI, path coreiface.Path, localPath string) error {
	readCloser, err := api.Unixfs().Cat(ctx, path)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	data, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(localPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

type Command interface {
	Name() string
	Usage() string
	Action() error
}

func loadCommand(path string) (*cli.Command, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symbol, err := p.Lookup("Command")
	if err != nil {
		return nil, err
	}

	command, ok := symbol.(Command)
	if !ok {
		return nil, fmt.Errorf("symbol not a Command")
	}

	return &cli.Command{
		Name:  command.Name(),
		Usage: command.Usage(),
		Action: func(c *cli.Context) error {
			return command.Action()
		},
	}, nil
}
