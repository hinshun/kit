package config

type Manifest struct {
	Revision string       `json:"revision"`
	Usage    string       `json:"usage"`
	Type     ManifestType `json:"type"`

	// Namespace only fields.
	Plugins Plugins `json:"plugins,omitempty"`

	// Command only fields.
	Hash  string `json:"hash,omitempty"`
	Args  []Arg  `json:"args,omitempty"`
	Flags []Flag `json:"flags,omitempty"`
}

type ManifestType string

var (
	CommandManifest   ManifestType = "command"
	NamespaceManifest ManifestType = "namespace"
)

type Arg struct {
	Type  string `json:"type"`
	Usage string `json:"usage"`
}

type Flag struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Usage string `json:"usage"`
}
