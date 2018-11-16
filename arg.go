package kit

import (
	"context"
)

type Arg interface {
	Type() string
	Usage() string
	Set(v string) error
	Autocomplete(ctx context.Context, input string) []Completion
}

type Completion struct {
	Group    string
	Wordlist []string
}

type stringArg struct {
	name  string
	usage string
	path  *string
}

func StringArg(name, usage string, path *string) Arg {
	return &stringArg{name, usage, path}
}

func (a *stringArg) Type() string {
	return a.name
}

func (a *stringArg) Usage() string {
	return a.usage
}

func (a *stringArg) Set(v string) error {
	*a.path = v
	return nil
}

func (a *stringArg) Autocomplete(ctx context.Context, input string) []Completion {
	return nil
}
