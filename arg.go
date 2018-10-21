package kit

import (
	"github.com/hinshun/kit/args"
)

type Arg interface {
	Type() string
	Usage() string
	Set(v string) error
}

func CommandPathArg(path *string, usage string) Arg {
	return args.NewCommandPathArg(path, usage)
}

func ManifestArg(manifest *string, usage string) Arg {
	return args.NewManifestArg(manifest, usage)
}
