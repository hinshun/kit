package cli

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hinshun/kit/config"
)

type VerifyFunc func(c *Cli) error

type AutocompleteFunc func(ctx context.Context, input string) []Completion

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

func VerifyNamespace(plugin *Plugin, args []string, depth int) VerifyFunc {
	return func(c *Cli) error {
		if len(args[depth:]) > 0 {
			return fmt.Errorf(
				"%s does not have command %s",
				strings.Join(c.DecorateCommandPath(plugin.CommandPath), " "),
				strings.Join(c.DecorateCommandPath(args[depth:]), " "),
			)
		}

		return nil
	}
}

func VerifyCommand(plugin *Plugin, kitCmd Command, args []string) VerifyFunc {
	return func(c *Cli) error {
		name := plugin.CommandPath[len(plugin.CommandPath)-1]
		flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
		flagSet.SetOutput(ioutil.Discard)

		for _, flag := range kitCmd.Flags() {
			flag.Set(flagSet)
		}

		err := flagSet.Parse(args)
		if err != nil {
			return err
		}

		parsedArgs := flagSet.Args()
		if len(plugin.Args) != len(parsedArgs) {
			return fmt.Errorf(
				"%s requires %s args but got %s args.",
				strings.Join(c.DecorateCommandPath(plugin.CommandPath), " "),
				c.theme.Value.Sprintf("%d", len(plugin.Args)),
				c.theme.Value.Sprintf("%d", len(parsedArgs)),
			)
		}

		for i, arg := range kitCmd.Args() {
			err = arg.Set(parsedArgs[i])
			if err != nil {
				return fmt.Errorf(
					"%s is not a valid %s: %s",
					c.theme.Value.Sprint(parsedArgs[i]),
					c.DecorateArg(arg.Type()),
					err,
				)
			}
		}

		return nil
	}
}

func AutocompleteNamespace(plugin *Plugin, args []string, depth int) AutocompleteFunc {
	return func(ctx context.Context, input string) []Completion {
		var wordlist []string
		for _, plugin := range plugin.Plugins {
			wordlist = append(wordlist, plugin.CommandPath[len(plugin.CommandPath)-1])
		}

		return []Completion{
			{
				Group:    "commands",
				Wordlist: wordlist,
			},
		}
	}
}

func AutocompleteCommand(cmd Command, args []string, depth int) AutocompleteFunc {
	return func(ctx context.Context, input string) []Completion {
		posArgIndex := len(args[depth:])
		if posArgIndex > len(cmd.Args())-1 {
			return nil
		}
		return cmd.Args()[posArgIndex].Autocomplete(ctx, input)
	}
}
