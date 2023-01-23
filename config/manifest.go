package config

import "fmt"

type Manifest struct {
	// Revision optionally provides the Version control revision hash.
	Revision string `json:"revision,omitempty"`

	// Description of the plugin.
	Usage string `json:"usage"`

	// Type of the plugin.
	Type ManifestType `json:"type"`

	// ---------------------
	// Namespace only fields

	// Plugins specify the plugins within the namespace.
	Plugins Plugins `json:"plugins,omitempty"`

	// -------------------
	// Command only fields

	// Platforms specify the path to the command.
	Platforms []Platform `json:"platforms,omitempty"`

	Args  []Arg  `json:"args,omitempty"`
	Flags []Flag `json:"flags,omitempty"`
}

func (m Manifest) MatchPlatform(os, arch string) (path string, err error) {
	for _, platform := range m.Platforms {
		if platform.OS == os && platform.Arch == arch {
			return platform.Path, nil
		}
	}
	return "", fmt.Errorf("unable to find path for platform %s %s", os, arch)
}

type ManifestType string

var (
	ManifestExternal  ManifestType = "external"
	ManifestCommand   ManifestType = "command"
	ManifestNamespace ManifestType = "namespace"
)

type Platform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
	Path string `json:"path"`
}

type Arg struct {
	Type  string `json:"type"`
	Usage string `json:"usage"`
}

type Flag struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Usage string `json:"usage"`
}
