package kit

import (
	"github.com/hinshun/kit/args"
)

type Arg interface {
	Type() string
	Usage() string
	Set(v string) error
}

func StringArg(name, usage string, path *string) Arg {
	return args.NewStringArg(name, usage, path)
}

func CommandPathArg(usage string, path *string) Arg {
	return args.NewCommandPathArg(usage, path)
}

func ManifestArg(usage string, manifest *string) Arg {
	return args.NewManifestArg(usage, manifest)
}
