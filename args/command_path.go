package args

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

var CommandPathPattern = regexp.MustCompile(`/(\w+/)*(\w+\?)?`)

type CommandPathArg struct {
	usage string
	path  *string
}

func NewCommandPathArg(usage string, path *string) *CommandPathArg {
	return &CommandPathArg{usage, path}
}

func (a *CommandPathArg) Type() string {
	return "command path"
}

func (a *CommandPathArg) Usage() string {
	return a.usage
}

func (a *CommandPathArg) Set(v string) error {
	if !CommandPathPattern.MatchString(v) {
		regex := color.New(color.FgBlue).Sprintf("%s", CommandPathPattern.String())
		return fmt.Errorf("did not match regex %s", regex)
	}
	*a.path = v
	return nil
}
