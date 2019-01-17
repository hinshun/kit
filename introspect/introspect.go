package introspect

import (
	"context"

	"github.com/hinshun/kit/config"
)

type API interface {
	GetManifest(ctx context.Context, plugin config.Plugin) (config.Manifest, error)
	Options() Options
	Theme() Theme
}

type kitKey struct{}

func WithKit(ctx context.Context, api API) context.Context {
	return context.WithValue(ctx, kitKey{}, api)
}

func Kit(ctx context.Context) API {
	api, ok := ctx.Value(kitKey{}).(API)
	if !ok {
		return nil
	}
	return api
}
