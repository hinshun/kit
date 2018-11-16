package plugin

import (
	"context"
	"fmt"
	"regexp"

	"github.com/fatih/color"
	"github.com/hinshun/kit"
)

var CommandPathPattern = regexp.MustCompile(`/(\w+/)*(\w+\?)?`)

type commandPathArg struct {
	usage string
	path  *string
}

func CommandPathArg(usage string, path *string) kit.Arg {
	return &commandPathArg{usage, path}
}

func (a *commandPathArg) Type() string {
	return "command path"
}

func (a *commandPathArg) Usage() string {
	return a.usage
}

func (a *commandPathArg) Set(v string) error {
	if !CommandPathPattern.MatchString(v) {
		regex := color.New(color.FgBlue).Sprintf("%s", CommandPathPattern.String())
		return fmt.Errorf("did not match regex %s", regex)
	}
	*a.path = v
	return nil
}

func (a *commandPathArg) Autocomplete(ctx context.Context, input string) []kit.Completion {
	return nil
}
