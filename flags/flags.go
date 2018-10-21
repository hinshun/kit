package flags

import "flag"

type StringFlag struct {
	name  string
	usage string
	value string
	dst   *string
}

func NewStringFlag(name, usage, value string, dst *string) *StringFlag {
	return &StringFlag{name, usage, value, dst}
}

func (f *StringFlag) Name() string {
	return f.name
}

func (f *StringFlag) Type() string {
	return "string"
}

func (f *StringFlag) Usage() string {
	return f.usage
}

func (f *StringFlag) Set(flagSet *flag.FlagSet) {
	flagSet.StringVar(f.dst, f.name, f.value, f.usage)
}
