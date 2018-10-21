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
}

func StringFlag(name, usage, value string, dst *string) Flag {
	return flags.NewStringFlag(name, usage, value, dst)
}
