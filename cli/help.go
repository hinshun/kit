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

	{{command "kit"}} {{globalFlag "global options"}} {{command "command"}} {{globalFlag "options"}} {{arg "arguments"}}

{{header "Commands:"}}{{if not .Plugins}}
  No plugins in {{join (commandPath .NamespacePath) " "}}.{{end}}{{range .Plugins}}
  {{join (commandPath (offset .CommandPath 1)) " "}} {{if .Flags}}{{join (flags .Flags) " "}} {{globalFlag "--"}} {{end}}{{join (args .Args) " "}}
		{{.Usage}}{{range .Flags}}
		{{flag .}}: {{.Usage}}{{end}}{{range .Args}}
		{{arg .Type}}: {{.Usage}}{{end}}
{{end}}{{if .UsageError}}
{{usageError "Usage error:"}}
  {{.UsageError}}
{{end}}`

func (c *Cli) printHelp(plugins []*Plugin) error {
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

	c.Plugins = plugins
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
	return c.theme.Title.Sprint(header)
}

func (c *Cli) DecorateUsageError(header string) string {
	return c.theme.Error.Sprint(header)
}

func (c *Cli) DecorateCommandPath(commandPath []string) []string {
	var decorated []string
	for _, command := range commandPath {
		decorated = append(decorated, c.DecorateCommand(command))
	}
	return decorated
}

func (c *Cli) DecorateCommand(command string) string {
	return c.theme.Keyword.Sprint(command)
}

func (c *Cli) DecorateArgs(inputs []config.Arg) []string {
	var args []string
	for _, input := range inputs {
		args = append(args, c.DecorateArg(input.Type))
	}
	return args
}

func (c *Cli) DecorateArg(arg string) string {
	return c.theme.Value.Sprintf("<%s>", arg)
}

func (c *Cli) DecorateGlobalFlag(flag string) string {
	return fmt.Sprintf("[%s]", c.theme.Option.Sprint(flag))
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
	if flag.Type == "bool" {
		output = fmt.Sprintf("--%s", flag.Name)
	} else {
		output = fmt.Sprintf("--%s <%s>", flag.Name, flag.Type)
	}
	return fmt.Sprintf("[%s]", c.theme.Option.Sprint(output))
}
