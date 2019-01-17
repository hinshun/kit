package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/hinshun/kit/config"
	"github.com/hinshun/kit/content/ipfsstore"
	"github.com/hinshun/kit/content/kitstore"
	"github.com/hinshun/kit/content/localstore"
	"github.com/hinshun/kit/introspect"
)

type Kit struct {
	cli *Cli
}

func NewKit() *Kit {
	s := kitstore.NewStore(
		localstore.NewStore(
			ipfsstore.NewStore(),
		),
	)

	return &Kit{
		cli: NewCli(
			NewLoader(s),
		),
	}
}

func (k *Kit) Run(ctx context.Context, args []string) error {
	err := k.cli.Parse(args)
	if err != nil {
		return err
	}

	cfg, err := config.New(k.cli.ConfigPath())
	if err != nil {
		return err
	}

	plugin := config.Plugin{
		Plugins: config.Plugins{
			{
				Name:     "kit",
				Usage:    "Composable command-line toolkit.",
				Manifest: cfg.Manifest,
				Plugins:  cfg.Plugins,
			},
		},
	}

	manifest, err := k.cli.GetManifest(ctx, plugin)
	if err != nil {
		return err
	}

	merged := manifest.Plugins.Merge(plugin.Plugins)
	if len(merged) == 0 {
		plugin.Plugins = config.InitConfig.Plugins
	}

	cliArgs := k.cli.Args()
	command, err := k.cli.GetCommand(ctx, plugin, cliArgs)
	if err != nil {
		return err
	}

	if *k.cli.autocomplete != "" {
		ctx = introspect.WithKit(ctx, k.cli)
		completions := command.Autocomplete(ctx, *k.cli.autocomplete)
		switch *k.cli.autocomplete {
		case "bash", "fish":
			var wordlist []string
			for _, completion := range completions {
				wordlist = append(wordlist, completion.Wordlist...)
			}
			fmt.Printf("%s", strings.Join(wordlist, " "))
		case "zsh":
			var shellCmds []string
			for _, completion := range completions {
				shellCmds = append(
					shellCmds,
					fmt.Sprintf("local -a %s", completion.Group),
					fmt.Sprintf("%s=(%s)", completion.Group, strings.Join(completion.Wordlist, " ")),
					fmt.Sprintf("_describe %s %s", completion.Group, completion.Group),
				)
			}
			fmt.Printf("%s", strings.Join(shellCmds, ";"))
		}
		return nil
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

	ctx = introspect.WithKit(ctx, k.cli)
	return command.Action(ctx)
}
