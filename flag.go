package kit

import (
	"flag"

	"github.com/hinshun/kit/flags"
)

type Flag interface {
	Name() string
	Type() string
	Usage() string
	Set(*flag.FlagSet)
	Autocomplete(input string) []string
}

func StringFlag(name, usage, value string, dst *string) Flag {
	return flags.NewStringFlag(name, usage, value, dst)
}

func BoolFlag(name, usage string, value bool, dst *bool) Flag {
	return flags.NewBoolFlag(name, usage, value, dst)
}
