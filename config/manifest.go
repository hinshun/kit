package config

type Manifest struct {
	Revision string       `json:"revision"`
	Usage    string       `json:"usage"`
	Type     ManifestType `json:"type"`

	// Namespace only fields.
	Plugins Plugins `json:"plugins,omitempty"`

	// Command only fields.
	Hash  string  `json:"hash,omitempty"`
	Args  []Input `json:"args,omitempty"`
	Flags []Input `json:"flags,omitempty"`
}

type ManifestType string

var (
	CommandManifest   ManifestType = "command"
	NamespaceManifest ManifestType = "namespace"
)

type Input struct {
	Type  string `json:"type"`
	Usage string `json:"usage"`
}

func (i Input) String() string {
	return i.Type
}
