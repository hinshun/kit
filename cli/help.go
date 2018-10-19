package cli

import (
	"fmt"
	"html/template"
	"strings"
	"text/tabwriter"

	"github.com/hinshun/kit/config"
)

var HelpTemplate = `kit

Commands:{{range .Commands}}
  {{join .Names " "}} {{joinInputs .Args " "}} {{if .Flags}}[{{joinInputs .Flags " "}}]{{end}}
    {{.Usage}}{{range .Args}}
		{{.Name}}: {{.Usage}}{{end}}{{range .Flags}}
		{{.Name}}: {{.Usage}}{{end}}{{end}}
`

func (c *Cli) PrintHelp(commands []*Command) error {
	c.Commands = commands

	funcs := template.FuncMap{
		"join": func(strs []string, separator string) string {
			return strings.Join(strs, separator)
		},
		"joinInputs": func(inputs []config.Input, separator string) string {
			var strs []string
			for _, input := range inputs {
				strs = append(strs, fmt.Sprintf("%s", input))
			}

			return strings.Join(strs, separator)
		},
	}

	w := tabwriter.NewWriter(c.stdio.Out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcs).Parse(HelpTemplate))

	err := t.Execute(w, c)
	if err != nil {
		return err
	}

	return w.Flush()
}
