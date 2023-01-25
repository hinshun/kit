//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative command.proto

package api

import (
	plugin "github.com/hashicorp/go-plugin"
)

var (
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "KIT_PLUGIN",
		MagicCookieValue: "kit",
	}
)
