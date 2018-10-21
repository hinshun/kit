package cli

import (
	"context"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/content/local"
)

type Kit struct {
	cli *Cli
}

func NewKit() *Kit {
	return &Kit{
		cli: NewCli(NewLoader(local.NewStore())),
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

	plugin := config.Plugin{
		Name:     "kit",
		Manifest: cfg.Manifest,
		Plugins:  cfg.Plugins,
	}
	command, err := k.cli.GetCommand(ctx, plugin, k.cli.flagSet.Args())
	if err != nil {
		return err
	}

	if !*k.cli.help {
		err = command.Verify(k.cli)
		if err != nil {
			k.cli.UsageError = err
		}
	}

	if *k.cli.help || k.cli.UsageError != nil || command.Action == nil {
		commands := command.Commands
		if len(commands) == 0 {
			commands = []*Command{command}
		}
		return k.cli.PrintHelp(commands)
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
