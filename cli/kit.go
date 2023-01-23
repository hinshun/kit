package cli
import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hinshun/kit/config"
)

var (
	KitDir         = ".kit"
	ConfigFilename = "config.json"
	ConfigPath     = filepath.Join(KitDir, ConfigFilename)
)

type Kit struct {
	cli *Cli
}

func NewKit() *Kit {
	return &Kit{
		cli: NewCli(NewLoader()),
	}
}

func (k *Kit) Run(ctx context.Context, args []string) error {
	configPath := filepath.Join(os.Getenv("HOME"), ConfigPath)
	cfg, err := config.New(configPath)
	if err != nil {
		return err
	}

	plugin := config.Plugin{
		Plugins: config.Plugins{
			{
				Name:    "kit",
				Usage:   "Composable command-line toolkit.",
				Plugins: cfg.Plugins,
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

	err = k.cli.Parse(args)
	if err != nil {
		return err
	}

	resolved, err := k.cli.GetPlugin(ctx, plugin, k.cli.Args())
	if err != nil {
		return err
	}

	if k.cli.options.Autocomplete != "" {
		// ctx = introspect.WithKit(ctx, k.cli)
		completions := resolved.Autocomplete(ctx, k.cli.options.Autocomplete)
		switch k.cli.options.Autocomplete {
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

	if !k.cli.options.Help {
		err = resolved.Verify(k.cli)
		if err != nil {
			k.cli.UsageError = err
		}
	}

	if k.cli.options.Help || k.cli.UsageError != nil || resolved.Action == nil {
		if resolved.Action == nil {
			k.cli.SetNamespaceUsage(resolved.CommandPath, resolved.Usage)
			return k.cli.printHelp(resolved.Plugins)
		} else {
			namespace := config.Plugin{Plugins: merged}.FindParent(resolved.CommandPath)
			namespaceManifest, err := k.cli.GetManifest(ctx, namespace)
			if err != nil {
				return err
			}

			k.cli.SetNamespaceUsage(resolved.CommandPath[:len(resolved.CommandPath)-1], namespaceManifest.Usage)
			return k.cli.printHelp([]*Plugin{resolved})
		}
	}

	// ctx = introspect.WithKit(ctx, k.cli)
	return resolved.Action(ctx)
}
