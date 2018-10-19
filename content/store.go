package content

import (
	"context"
)

type Store interface {
	Get(ctx context.Context, manifest string) (string, error)
}
