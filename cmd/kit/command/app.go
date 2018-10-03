package command

import (
	"context"

	"github.com/ipfs/go-ipfs/core/coreapi"

	"github.com/codegangsta/cli"
	"github.com/hinshun/kit"
	"github.com/hinshun/kit/ipfs"
	"github.com/hinshun/kit/linker"
)

func App(ctx context.Context) (*cli.App, error) {
	commands, err := loadCommands(ctx)
	if err != nil {
		return nil, err
	}

	app := cli.NewApp()
	app.Commands = commands
	return app, nil
}

func loadCommands(ctx context.Context) ([]cli.Command, error) {
	cfg, err := kit.ParseConfig(".kit/config.json")
	if err != nil {
		return nil, err
	}

	n, err := ipfs.NewNode(ctx, cfg)
	if err != nil {
		return nil, err
	}

	api := coreapi.NewCoreAPI(n)
	paths, err := ipfs.SyncCommands(ctx, cfg, api, cfg.Commands)
	if err != nil {
		return nil, err
	}

	var commands []cli.Command
	for _, path := range paths {
		command, err := linker.LinkCommand(path)
		if err != nil {
			return nil, err
		}

		commands = append(commands, cli.Command{
			Name:  command.Name(),
			Usage: command.Usage(),
			Action: func(c *cli.Context) error {
				return command.Action(context.TODO())
			},
		})
	}

	return commands, nil
}
