package kitstore

import (
	"context"

	"github.com/hinshun/kit/content"
)

var (
	InitDigest    string
	PluginDigest  string
	AddDigest     string
	RmDigest      string
	PublishDigest string
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
		digest = AddDigest
	case "/kit/plugin/rm":
		digest = RmDigest
	case "/kit/plugin/publish":
		digest = PublishDigest
	}

	return s.Store.Get(ctx, digest)
}
