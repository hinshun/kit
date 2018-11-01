package cli

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/hinshun/kit/config"
)

var HelpTemplate = `{{header "Usage:"}}
	{{join (commandPath .NamespacePath) " "}}{{if .NamespaceUsage}} - {{.NamespaceUsage}}{{end}}

	{{command "kit"}} {{globalFlag "global options"}} {{join (append (commandPath (offset .NamespacePath 1) ) (command "command")) " "}} {{globalFlag "options"}} {{arg "arguments"}}

{{header "Commands:"}}{{if not .Commands}}
  No commands in {{join (commandPath .NamespacePath) " "}}.{{end}}{{range .Commands}}
  {{join (commandPath .CommandPath) " "}} {{if .Flags}}{{join (flags .Flags) " "}} {{globalFlag "--"}} {{end}}{{join (args .Args) " "}}
		{{if eq .Usage ""}}A plugin namespace.{{else}}{{.Usage}}{{end}}{{range .Flags}}
		{{flag .}}: {{.Usage}}{{end}}{{range .Args}}
		{{arg .Type}}: {{.Usage}}{{end}}
{{end}}{{if .UsageError}}
{{usageError "Usage error:"}}
  {{.UsageError}}
{{end}}`

func (c *Cli) PrintHelp(commands []*Command) error {
	funcs := template.FuncMap{
		"join":        join,
		"offset":      offset,
		"append":      appendStr,
		"header":      c.DecorateHeader,
		"usageError":  c.DecorateUsageError,
		"commandPath": c.DecorateCommandPath,
		"command":     c.DecorateCommand,
		"args":        c.DecorateArgs,
		"arg":         c.DecorateArg,
		"globalFlag":  c.DecorateGlobalFlag,
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

func offset(strs []string, offset int) []string {
	return strs[offset:]
}

func appendStr(strs []string, elems ...string) []string {
	return append(strs, elems...)
}

func (c *Cli) DecorateHeader(header string) string {
	return c.headerColor.Sprint(header)
}

func (c *Cli) DecorateUsageError(header string) string {
	return c.usageErrorColor.Sprint(header)
}

func (c *Cli) DecorateCommandPath(commandPath []string) []string {
	var decorated []string
	for _, command := range commandPath {
		decorated = append(decorated, c.DecorateCommand(command))
	}
	return decorated
}

func (c *Cli) DecorateCommand(command string) string {
	return c.commandColor.Sprint(command)
}

func (c *Cli) DecorateArgs(inputs []config.Arg) []string {
	var args []string
	for _, input := range inputs {
		args = append(args, c.DecorateArg(input.Type))
	}
	return args
}

func (c *Cli) DecorateArg(arg string) string {
	return c.argColor.Sprintf("<%s>", arg)
}

func (c *Cli) DecorateGlobalFlag(flag string) string {
	return fmt.Sprintf("[%s]", c.flagColor.Sprint(flag))
}

func (c *Cli) DecorateFlags(inputs []config.Flag) []string {
	var flags []string
	for _, input := range inputs {
		flags = append(flags, c.DecorateFlag(input))
	}
	return flags
}

func (c *Cli) DecorateFlag(flag config.Flag) string {
	var output string
	if flag.Type == "" {
		output = fmt.Sprintf("--%s", flag.Name)
	} else {
		output = fmt.Sprintf("--%s <%s>", flag.Name, flag.Type)
	}
	return fmt.Sprintf("[%s]", c.flagColor.Sprint(output))
}
