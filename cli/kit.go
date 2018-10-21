package cli

import (
	"context"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/content/local"
)

type Kit struct {
	cli    *Cli
	loader *Loader
}

func NewKit() *Kit {
	c := NewCli()

	return &Kit{
		cli:    c,
		loader: NewLoader(c, local.NewStore()),
	}
}

func (k *Kit) Run(ctx context.Context, args []string) error {
	err := k.cli.Parse(args[1:])
	if err != nil {
		return err
	}

	cfg, err := config.New(k.cli.ConfigPath())
	if err != nil {
		return err
	}

	command, err := k.loader.GetCommand(ctx, cfg, k.cli.flagSet.Args())
	if err != nil {
		return err
	}

	ctx = kit.WithKit(ctx, k.cli)
	return command.Action(ctx)
}

func (k *Kit) PrintHelp(ctx context.Context, command *Command) error {
	if len(command.CommandPath) > 0 {
		return k.cli.PrintHelp([]*Command{command})
	}
	return command.Action(ctx)
}
