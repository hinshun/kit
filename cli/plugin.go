package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
)

type Plugin struct {
	CommandPath  []string
	Usage        string
	Args         []config.Arg
	Flags        []config.Flag
	Verify       VerifyFunc
	Autocomplete AutocompleteFunc
	Action       func(ctx context.Context) error
	Plugins      []*Plugin
}

type VerifyFunc func(ctx context.Context, c *Cli) error

type AutocompleteFunc func(ctx context.Context, input string) ([]kit.Completion, error)

func VerifyNamespace(plugin *Plugin, shellArgs []string, depth int) VerifyFunc {
	return func(ctx context.Context, c *Cli) error {
		if len(shellArgs[depth:]) > 0 {
			return fmt.Errorf(
				"%s does not have command %s",
				strings.Join(c.DecorateCommandPath(plugin.CommandPath), " "),
				strings.Join(c.DecorateCommandPath(shellArgs[depth:]), " "),
			)
		}

		return nil
	}
}

func VerifyCommand(plugin *Plugin, cmd kit.RemoteCommand, shellArgs []string) VerifyFunc {
	return func(ctx context.Context, c *Cli) error {
		cmdFlags, err := cmd.Flags()
		if err != nil {
			return err
		}

		cmdArgs, err := cmd.Args()
		if err != nil {
			return err
		}

		flags := append(c.Flags(), cmdFlags...)
		return Parse(ctx, flags, cmdArgs, shellArgs)

		// parsedArgs := flagSet.Args()
		// if len(plugin.Args) != len(parsedArgs) {
		// 	return fmt.Errorf(
		// 		"%s requires %s args but got %s args.",
		// 		strings.Join(c.DecorateCommandPath(plugin.CommandPath), " "),
		// 		c.theme.Value.Sprintf("%d", len(plugin.Args)),
		// 		c.theme.Value.Sprintf("%d", len(parsedArgs)),
		// 	)
		// }

		// for i, arg := range kitCmd.Args() {
		// 	err = arg.Set(ctx, parsedArgs[i])
		// 	if err != nil {
		// 		return fmt.Errorf(
		// 			"%s is not a valid %s: %s",
		// 			c.theme.Value.Sprint(parsedArgs[i]),
		// 			c.DecorateArg(arg.Type()),
		// 			err,
		// 		)
		// 	}
		// }

		// return nil
	}
}

func AutocompleteNamespace(plugin *Plugin, args []string, depth int) AutocompleteFunc {
	return func(ctx context.Context, input string) ([]kit.Completion, error) {
		var wordlist []string
		for _, plugin := range plugin.Plugins {
			wordlist = append(wordlist, plugin.CommandPath[len(plugin.CommandPath)-1])
		}

		return []kit.Completion{
			{
				Group:    "commands",
				Wordlist: wordlist,
			},
		}, nil
	}
}

func AutocompleteCommand(cmd kit.RemoteCommand, args []string, depth int) AutocompleteFunc {
	return func(ctx context.Context, input string) ([]kit.Completion, error) {
		cmdArgs, err := cmd.Args()
		if err != nil {
			return nil, err
		}

		posArgIndex := len(args[depth:])
		if posArgIndex > len(cmdArgs)-1 {
			return nil, nil
		}
		return cmdArgs[posArgIndex].Autocomplete(ctx, input)
	}
}
