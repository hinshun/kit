package cli

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

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
	err := k.flagSet.Parse(args)
	if err != nil {
		return err
	}

	kitDir := filepath.Join(os.Getenv("HOME"), ".kit")
	cfg, err := config.New(filepath.Join(kitDir, "config.json"))
	if err != nil {
		return err
	}

	command, err := k.loader.GetCommand(ctx, cfg, k.flagSet.Args())
	if err != nil {
		return err
	}

	if *k.help && len(command.Names) > 0 {
		return k.cli.PrintHelp([]*Command{command})
	}

	return command.Action(ctx)
}
