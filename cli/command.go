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

type Command struct {
	CommandPath []string
	Usage       string
	Args        []config.Arg
	Flags       []config.Flag
	Action      func(ctx context.Context) error
}

func (c *Cli) VerifyNamespace(cliCmd *Command, args []string, depth int) error {
	if len(args[depth:]) > 0 {
		return fmt.Errorf(
			"%s does not have command %s",
			strings.Join(c.DecorateCommandPath(cliCmd.CommandPath), " "),
			strings.Join(c.DecorateCommandPath(args[depth:]), " "),
		)
	}

	return nil
}

func (c *Cli) VerifyCommand(cliCmd *Command, kitCmd kit.Command, args []string) error {
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
			c.argColor.Sprintf("%d", len(cliCmd.Args)),
			c.argColor.Sprintf("%d", len(parsedArgs)),
		)
	}

	for i, arg := range kitCmd.Args() {
		err = arg.Set(parsedArgs[i])
		if err != nil {
			return fmt.Errorf(
				"%s is not a valid %s: %s",
				c.argColor.Sprint(parsedArgs[i]),
				c.DecorateArg(arg.Type()),
				err,
			)
		}
	}

	return nil
}
