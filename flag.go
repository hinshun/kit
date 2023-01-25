package kit

import (
	"context"
	"strconv"
)

type Flag interface {
	Name() string
	Type() string
	Usage() string
	Set(ctx context.Context, v string) error
	Autocomplete(ctx context.Context, input string) ([]Completion, error)
}

type stringFlag struct {
	name  string
	usage string
	value string
	dst   *string
}

func StringFlag(name, usage, value string, dst *string) Flag {
	return &stringFlag{name, usage, value, dst}
}

func (f *stringFlag) Name() string {
	return f.name
}

func (f *stringFlag) Type() string {
	return "string"
}

func (f *stringFlag) Usage() string {
	return f.usage
}

func (f *stringFlag) Set(ctx context.Context, v string) error {
	*f.dst = v
	return nil
}

func (f *stringFlag) Autocomplete(ctx context.Context, input string) ([]Completion, error) {
	return nil, nil
}

type boolFlag struct {
	name  string
	usage string
	value bool
	dst   *bool
}

func BoolFlag(name, usage string, value bool, dst *bool) Flag {
	return &boolFlag{name, usage, value, dst}
}

func (f *boolFlag) Name() string {
	return f.name
}

func (f *boolFlag) Type() string {
	return "bool"
}

func (f *boolFlag) Usage() string {
	return f.usage
}

func (f *boolFlag) Autocomplete(ctx context.Context, input string) ([]Completion, error) {
	return nil, nil
}

func (f *boolFlag) Set(ctx context.Context, v string) error {
	b, err := strconv.ParseBool(v)
	if err != nil {
		return err
	}
	*f.dst = b
	return nil
}
