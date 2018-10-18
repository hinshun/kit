package cli

import (
	"context"
	"flag"
	"io/ioutil"
	"os"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/content/ipfs"
	"github.com/hinshun/kit/loader"
)

type Kit struct {
	loader  *loader.Loader
	flagSet *flag.FlagSet
	help    *bool
}

func NewKit() kit.Kit {
	flagSet := flag.NewFlagSet("kit", flag.ContinueOnError)
	flagSet.SetOutput(ioutil.Discard)

	return &Kit{
		loader:  loader.New(ipfs.NewStore()),
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

	command, err := k.loader.Load(ctx, cfg.Plugins, k.flagSet.Args())
	if err != nil {
		return err
	}

	if *k.help {
		return PrintHelp(os.Stdout, command)
	}

	return command.Run(ctx)
}
