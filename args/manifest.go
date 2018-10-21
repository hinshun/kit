package args

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

var ManifestPattern = regexp.MustCompile(`/kit|ipfs/(\w+/)*(\w+\?)?`)

type ManifestArg struct {
	manifest *string
	usage    string
}

func NewManifestArg(manifest *string, usage string) *ManifestArg {
	return &ManifestArg{manifest, usage}
}

func (a *ManifestArg) Type() string {
	return "manifest"
}

func (a *ManifestArg) Usage() string {
	return a.usage
}

func (a *ManifestArg) Set(v string) error {
	if !ManifestPattern.MatchString(v) {
		regex := color.New(color.FgBlue).Sprintf("%s", ManifestPattern.String())
		return fmt.Errorf("did not match regex %s", regex)
	}
	*a.manifest = v
	return nil
}
