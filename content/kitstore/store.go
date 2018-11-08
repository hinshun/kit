package kitstore

import (
	"context"

	"github.com/hinshun/kit/content"
)

var (
	InitDigest          string
	PluginDigest        = "QmetM7PMkuGJtwBS5Lw57cCfMDSYk4UHZRPGNzq9JtGLCP"
	PluginAddDigest     string
	PluginRmDigest      string
	PluginPublishDigest string
)

type store struct {
	content.Store
}

func NewStore(s content.Store) content.Store {
	return &store{
		Store: s,
	}
}
func (s *store) Get(ctx context.Context, digest string) (string, error) {
	switch digest {
	case "/kit/init":
		digest = InitDigest
	case "/kit/plugin":
		digest = PluginDigest
	case "/kit/plugin/add":
		digest = PluginAddDigest
	case "/kit/plugin/rm":
		digest = PluginRmDigest
	case "/kit/plugin/publish":
		digest = PluginPublishDigest
	}

	return s.Store.Get(ctx, digest)
}
