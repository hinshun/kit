package cli

import (
	"context"

	"github.com/hinshun/kit/config"
)

type Command struct {
	Names  []string
	Usage  string
	Args   []config.Input
	Flags  []config.Input
	Action func(ctx context.Context) error
}
