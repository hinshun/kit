package kit

import (
	"context"
)

type Arg interface {
	Type() string
	Usage() string
	Set(ctx context.Context, v string) error
	Autocomplete(ctx context.Context, input string) ([]Completion, error)
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

func (a *stringArg) Set(ctx context.Context, v string) error {
	*a.path = v
	return nil
}

func (a *stringArg) Autocomplete(ctx context.Context, input string) ([]Completion, error) {
	return nil, nil
}
