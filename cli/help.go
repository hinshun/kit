package cli

import (
	"fmt"
	"html/template"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/hinshun/kit/config"
)

var HelpTemplate = `kit

{{header "Commands:"}}{{range .Commands}}
  {{join (names .Names) " "}} {{join (args .Args) " "}} {{if .Flags}}[{{join (flags .Flags) " "}}]{{end}}
    {{.Usage}}{{range .Args}}
		{{arg .Name}}: {{.Usage}}{{end}}{{range .Flags}}
		{{flag .Name}}: {{.Usage}}{{end}}{{end}}
`

func (c *Cli) PrintHelp(commands []*Command) error {
	c.Commands = commands

	funcs := template.FuncMap{
		"join":   join,
		"header": decorateHeader,
		"names":  decorateNames,
		"name":   decorateName,
		"args":   decorateArgs,
		"arg":    decorateArg,
		"flags":  decorateFlags,
		"flag":   decorateFlag,
	}

	w := tabwriter.NewWriter(c.stdio.Out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcs).Parse(HelpTemplate))

	err := t.Execute(w, c)
	if err != nil {
		return err
	}

	return w.Flush()
}

func join(strs []string, separator string) string {
	return strings.Join(strs, separator)
}

func decorateHeader(header string) string {
	return color.New(color.Bold, color.Underline).Sprint(header)
}

func decorateNames(names []string) []string {
	for i := 0; i < len(names); i++ {
		names[i] = decorateName(names[i])
	}
	return names
}

func decorateName(name string) string {
	return color.New(color.FgWhite, color.Underline).Sprint(name)
}

func decorateArgs(inputs []config.Input) []string {
	var args []string
	for _, input := range inputs {
		args = append(args, decorateArg(input.String()))
	}
	return args
}

func decorateArg(arg string) string {
	return fmt.Sprintf("<%s>", color.New(color.FgCyan).Sprint(arg))
}

func decorateFlags(inputs []config.Input) []string {
	var flags []string
	for _, input := range inputs {
		flags = append(flags, decorateFlag(input.String()))
	}
	return flags
}

func decorateFlag(flag string) string {
	return color.New(color.FgGreen).Sprint(flag)
}
