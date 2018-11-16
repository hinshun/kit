package kit

import (
	"context"
	"flag"
)

type Flag interface {
	Name() string
	Type() string
	Usage() string
	Set(*flag.FlagSet)
	Autocomplete(ctx context.Context, input string) []string
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

func (f *stringFlag) Set(flagSet *flag.FlagSet) {
	flagSet.StringVar(f.dst, f.name, f.value, f.usage)
}

func (f *stringFlag) Autocomplete(ctx context.Context, input string) []string {
	return nil
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
	return ""
}

func (f *boolFlag) Usage() string {
	return f.usage
}

func (f *boolFlag) Autocomplete(ctx context.Context, input string) []string {
	return nil
}

func (f *boolFlag) Set(flagSet *flag.FlagSet) {
	flagSet.BoolVar(f.dst, f.name, f.value, f.usage)
}
