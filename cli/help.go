package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/hinshun/kit/config"
)

var HelpTemplate = `{{header "Usage:"}}
  kit - Composable command-line toolkit.

	kit {{flag "global options"}} {{command "command"}} {{flag "command options"}}

{{header "Commands:"}}{{range .Commands}}
  {{join (commandPath .CommandPath) " "}} {{join (args .Args) " "}} {{if .Flags}}{{join (flags .Flags) " "}}{{end}}
    {{.Usage}}{{range .Args}}
		{{arg .Type}}: {{.Usage}}{{end}}{{range .Flags}}
		{{flag .Type}}: {{.Usage}}{{end}}{{end}}{{if .UsageError}}

{{usageError "Usage error:"}}
  {{.UsageError}}{{end}}
`

func (c *Cli) PrintHelp(commands []*Command) error {
	funcs := template.FuncMap{
		"join":        join,
		"header":      c.DecorateHeader,
		"usageError":  c.DecorateUsageError,
		"commandPath": c.DecorateCommandPath,
		"command":     c.DecorateCommand,
		"args":        c.DecorateArgs,
		"arg":         c.DecorateArg,
		"flags":       c.DecorateFlags,
		"flag":        c.DecorateFlag,
	}

	w := tabwriter.NewWriter(c.stdio.Out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcs).Parse(HelpTemplate))

	c.Commands = commands
	err := t.Execute(w, c)
	if err != nil {
		return err
	}

	return w.Flush()
}

func join(strs []string, separator string) string {
	return strings.Join(strs, separator)
}

func (c *Cli) DecorateHeader(header string) string {
	return c.headerColor.Sprint(header)
}

func (c *Cli) DecorateUsageError(header string) string {
	return c.usageErrorColor.Sprint(header)
}

func (c *Cli) DecorateCommandPath(commandPath []string) []string {
	for i := 0; i < len(commandPath); i++ {
		commandPath[i] = c.DecorateCommand(commandPath[i])
	}
	return commandPath
}

func (c *Cli) DecorateCommand(command string) string {
	return c.commandColor.Sprint(command)
}

func (c *Cli) DecorateArgs(inputs []config.Input) []string {
	var args []string
	for _, input := range inputs {
		args = append(args, c.DecorateArg(input.String()))
	}
	return args
}

func (c *Cli) DecorateArg(arg string) string {
	return c.argColor.Sprintf("<%s>", arg)
}

func (c *Cli) DecorateFlags(inputs []config.Input) []string {
	var flags []string
	for _, input := range inputs {
		flags = append(flags, c.DecorateFlag(input.String()))
	}
	return flags
}

func (c *Cli) DecorateFlag(flag string) string {
	return fmt.Sprintf("[%s]", c.flagColor.Sprint(flag))
}
