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
	config  *string
	help    *bool
}

func NewKit() *Kit {
	c := NewCli()

	flagSet := flag.NewFlagSet("kit", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return &Kit{
		cli:     c,
		loader:  NewLoader(c, local.NewStore()),
		flagSet: flagSet,
		config:  flagSet.String("config", filepath.Join(os.Getenv("HOME"), kit.ConfigPath), "path to kit config"),
		help:    flagSet.Bool("help", false, "display this help text"),
	}
}

func (k *Kit) Run(ctx context.Context, args []string) error {
	err := k.flagSet.Parse(args[1:])
	if err != nil {
		return err
	}

	k.cli.configPath = *k.config
	cfg, err := config.New(k.cli.ConfigPath())
	if err != nil {
		return err
	}

	command, err := k.loader.GetCommand(ctx, cfg, k.flagSet.Args())
	if err != nil {
		return err
	}

	if *k.help || k.cli.UsageError != nil {
		return k.PrintHelp(ctx, command)
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
