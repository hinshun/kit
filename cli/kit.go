package cli

import (
	"context"
	"flag"
	"io/ioutil"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/content/local"
)

type Kit struct {
	cli     *Cli
	loader  *Loader
	flagSet *flag.FlagSet
	help    *bool
}

func NewKit() kit.Kit {
	c := NewCli()

	flagSet := flag.NewFlagSet("kit", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return &Kit{
		cli:     c,
		loader:  NewLoader(c, local.NewStore()),
		flagSet: flagSet,
		help:    flagSet.Bool("help", false, "display this help text"),
	}
}

func (k *Kit) Run(ctx context.Context, args []string) error {
	err := k.flagSet.Parse(args[1:])
	if err != nil {
		return err
	}

	cfg, err := config.New(".kit/config.json")
	if err != nil {
		return err
	}

	command, err := k.loader.LoadCommand(ctx, cfg.Plugins, k.flagSet.Args())
	if err != nil {
		return err
	}

	if *k.help {
		return k.cli.PrintHelp([]*Command{command})
	}

	return command.Action(ctx)
}
