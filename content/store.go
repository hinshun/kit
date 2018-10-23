package content

import (
	"context"
)

type Store interface {
	Get(ctx context.Context, digest string) (string, error)
}
