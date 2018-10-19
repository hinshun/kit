package kit

import (
	"context"
	"io"
)

type Kit interface {
	Run(ctx context.Context, args []string) error
}

type Cli interface {
	Stdio() Stdio
	Args() Args
	Flags() Flags
}

type Stdio struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

type Args interface {
}

type Flags interface {
}

type Constructor func(cli Cli) (Command, error)

type Command interface {
	Usage() string
	Args() []Arg
	Flags() []Flag
	Run(ctx context.Context) error
}

type Arg interface {
}

type Flag interface {
}
