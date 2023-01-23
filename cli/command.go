package cli

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hinshun/kit"
	"github.com/hinshun/kit/config"
)

type VerifyFunc func(c *Cli) error

type AutocompleteFunc func(ctx context.Context, input string) []kit.Completion

type Command struct {
	CommandPath  []string
	Usage        string
	Args         []config.Arg
	Flags        []config.Flag
	Verify       VerifyFunc
	Autocomplete AutocompleteFunc
	Action       func(ctx context.Context) error
	Commands     []*Command
}

func VerifyNamespace(cliCmd *Command, args []string, depth int) VerifyFunc {
	return func(c *Cli) error {
		if len(args[depth:]) > 0 {
			return fmt.Errorf(
				"%s does not have command %s",
				strings.Join(c.DecorateCommandPath(cliCmd.CommandPath), " "),
				strings.Join(c.DecorateCommandPath(args[depth:]), " "),
			)
		}

		return nil
	}
}

func VerifyCommand(cliCmd *Command, kitCmd kit.Command, args []string) VerifyFunc {
	return func(c *Cli) error {
		name := cliCmd.CommandPath[len(cliCmd.CommandPath)-1]

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
		if len(cliCmd.Args) != len(parsedArgs) {
			return fmt.Errorf(
				"%s requires %s args but got %s args.",
				strings.Join(c.DecorateCommandPath(cliCmd.CommandPath), " "),
				c.theme.Value.Sprintf("%d", len(cliCmd.Args)),
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

func AutocompleteNamespace(cliCmd *Command, args []string, depth int) AutocompleteFunc {
	return func(ctx context.Context, input string) []kit.Completion {
		var wordlist []string
		for _, command := range cliCmd.Commands {
			wordlist = append(wordlist, command.CommandPath[len(command.CommandPath)-1])
		}

		return []kit.Completion{
			{
				Group:    "commands",
				Wordlist: wordlist,
			},
		}
	}
}

func AutocompleteCommand(kitCmd kit.Command, args []string, depth int) AutocompleteFunc {
	return func(ctx context.Context, input string) []kit.Completion {
		posArgIndex := len(args[depth:])

		kitArgs := kitCmd.Args()
		if posArgIndex > len(kitArgs)-1 {
			return nil
		}

		return kitArgs[posArgIndex].Autocomplete(ctx, input)
	}
}
