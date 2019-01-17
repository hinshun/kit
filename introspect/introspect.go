package introspect

import (
	"context"

	"github.com/hinshun/kit/config"
)

type KitAPI interface {
	GetManifest(ctx context.Context, plugin config.Plugin) (config.Manifest, error)
	ConfigPath() string
}

type kitKey struct{}

func WithKit(ctx context.Context, api KitAPI) context.Context {
	return context.WithValue(ctx, kitKey{}, api)
}

func Kit(ctx context.Context) KitAPI {
	api, ok := ctx.Value(kitKey{}).(KitAPI)
	if !ok {
		return nil
	}
	return api
}
