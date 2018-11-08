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

func (f *StringFlag) Autocomplete(input string) []string {
	return nil
}

type BoolFlag struct {
	name  string
	usage string
	value bool
	dst   *bool
}

func NewBoolFlag(name, usage string, value bool, dst *bool) *BoolFlag {
	return &BoolFlag{name, usage, value, dst}
}

func (f *BoolFlag) Name() string {
	return f.name
}

func (f *BoolFlag) Type() string {
	return ""
}

func (f *BoolFlag) Usage() string {
	return f.usage
}

func (f *BoolFlag) Autocomplete(input string) []string {
	return nil
}

func (f *BoolFlag) Set(flagSet *flag.FlagSet) {
	flagSet.BoolVar(f.dst, f.name, f.value, f.usage)
}
