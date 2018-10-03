package command

import (
	"context"
	"fmt"
	"time"

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
	cfg, err := kit.ParseConfig("kit.json")
	if err != nil {
		return nil, err
	}

	before := time.Now()
	n, err := ipfs.NewNode(ctx, cfg.Bootstrap)
	if err != nil {
		return nil, err
	}
	fmt.Printf("took %s to create bootstrapped ipfs node\n", time.Now().Sub(before))

	api := coreapi.NewCoreAPI(n)

	// f, err := os.Open("./bin/ls")
	// if err != nil {
	// 	return nil, err
	// }

	// p, err := api.Block().Put(ctx, f)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("IPFS block put: %s", p.Cid())
	// return nil, nil

	before = time.Now()
	paths, err := ipfs.SyncCommands(ctx, api, cfg.Commands)
	if err != nil {
		return nil, err
	}
	fmt.Printf("took %s to sync commands\n", time.Now().Sub(before))

	before = time.Now()
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
	fmt.Printf("took %s to link commands\n", time.Now().Sub(before))

	return commands, nil
}
