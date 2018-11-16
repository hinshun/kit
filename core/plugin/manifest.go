package plugin

import (
	"context"

	"github.com/hinshun/kit"
)

type manifestArg struct {
	usage    string
	manifest *string
}

func ManifestArg(usage string, manifest *string) kit.Arg {
	return &manifestArg{usage, manifest}
}

func (a *manifestArg) Type() string {
	return "manifest"
}

func (a *manifestArg) Usage() string {
	return a.usage
}

func (a *manifestArg) Set(v string) error {
	*a.manifest = v
	return nil
}

func (a *manifestArg) Autocomplete(ctx context.Context, input string) []kit.Completion {
	return []kit.Completion{
		{
			Group: "manifest",
			Wordlist: []string{
				"/kit/init",
				"/kit/plugin/add",
				"/kit/plugin/rm",
			},
		},
	}
}
