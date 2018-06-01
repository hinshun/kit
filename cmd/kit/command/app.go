package command

import (
	"context"
	"gx/ipfs/QmcKwjeebv5SX3VFUGDFa4BNMYhy14RRaCzQP7JN3UQDpB/go-ipfs/core/coreapi"
	"time"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/ipfs"
	"github.com/hinshun/kit/linker"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v2"
)

func App() (*cli.App, error) {
	ctx := context.Background()
	app := &cli.App{
		Name: "kit",
	}
	logrus.SetLevel(logrus.DebugLevel)

	commands, err := loadCommands(ctx)
	if err != nil {
		return app, err
	}

	app.Commands = commands
	return app, nil
}

func loadCommands(ctx context.Context) ([]*cli.Command, error) {
	cfg, err := kit.ParseConfig(".kit.toml")
	if err != nil {
		return nil, err
	}

	before := time.Now()
	n, err := ipfs.NewNode(ctx, cfg.Bootstrap)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("took %s to create bootstrapped ipfs node", time.Now().Sub(before))

	api := coreapi.NewCoreAPI(n)

	paths, err := ipfs.SyncCommands(ctx, api, cfg.Commands)
	if err != nil {
		return nil, err
	}

	var commands []*cli.Command
	for _, path := range paths {
		command, err := linker.LinkCommand(path)
		if err != nil {
			return nil, err
		}

		commands = append(commands, &cli.Command{
			Name:  command.Name(),
			Usage: command.Usage(),
			Action: func(c *cli.Context) error {
				return command.Action(ctx)
			},
		})
	}

	return commands, nil
}
