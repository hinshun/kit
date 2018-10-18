package cli

import (
	"html/template"
	"io"
	"text/tabwriter"

	"github.com/hinshun/kit"
)

var HelpTemplate = `NAME:
   {{.Name}}

COMMANDS:{{range .Commands}}
     {{.Name}}{{"\t"}}{{.Usage}}{{end}}
`

func PrintHelp(out io.Writer, command kit.Command) error {
	return HelpPrinter(out, HelpTemplate, nil)
}

func HelpPrinter(out io.Writer, templ string, data interface{}) error {
	w := tabwriter.NewWriter(out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Parse(templ))

	err := t.Execute(w, data)
	if err != nil {
		return err
	}

	return w.Flush()
}
