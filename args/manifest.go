package args

type ManifestArg struct {
	usage    string
	manifest *string
}

func NewManifestArg(usage string, manifest *string) *ManifestArg {
	return &ManifestArg{usage, manifest}
}

func (a *ManifestArg) Type() string {
	return "manifest"
}

func (a *ManifestArg) Usage() string {
	return a.usage
}

func (a *ManifestArg) Set(v string) error {
	*a.manifest = v
	return nil
}

func (a *ManifestArg) Autocomplete(input string) []string {
	return nil
}
