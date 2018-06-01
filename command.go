package kit

import "context"

type Command interface {
	Name() string
	Usage() string
	Action(ctx context.Context) error
}
