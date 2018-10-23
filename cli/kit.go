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
		Usage:    "Composable command-line toolkit.",
		Manifest: cfg.Manifest,
		Plugins:  cfg.Plugins,
	}

	manifest, err := k.cli.GetManifest(ctx, plugin)
	if err != nil {
		return err
	}

	merged := manifest.Plugins.Merge(plugin.Plugins)
	if len(merged) == 0 {
		plugin.Plugins = config.InitConfig.Plugins
	}

	cliArgs := k.cli.flagSet.Args()
	command, err := k.cli.GetCommand(ctx, plugin, cliArgs)
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
		if command.Action == nil {
			k.cli.SetNamespaceUsage(command.CommandPath, command.Usage)
			return k.cli.PrintHelp(command.Commands)
		} else {
			namespace := config.Plugin{Plugins: merged}.FindParent(command.CommandPath)
			namespaceManifest, err := k.cli.GetManifest(ctx, namespace)
			if err != nil {
				return err
			}

			k.cli.SetNamespaceUsage(command.CommandPath[:len(command.CommandPath)-1], namespaceManifest.Usage)
			return k.cli.PrintHelp([]*Command{command})
		}
	}

	ctx = kit.WithKit(ctx, k.cli)
	return command.Action(ctx)
}
