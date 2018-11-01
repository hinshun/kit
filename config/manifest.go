package config

type Manifest struct {
	Revision string       `json:"revision"`
	Usage    string       `json:"usage"`
	Type     ManifestType `json:"type"`

	// Namespace only fields.
	Plugins Plugins `json:"plugins,omitempty"`

	// Command only fields.
	Platforms []Platform `json:"platforms,omitempty"`
	Args      []Arg      `json:"args,omitempty"`
	Flags     []Flag     `json:"flags,omitempty"`
}

type ManifestType string

var (
	CommandManifest   ManifestType = "command"
	NamespaceManifest ManifestType = "namespace"
)

type Platform struct {
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	Digest       string `json:"digest"`
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
