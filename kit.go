package kit

import (
	"context"
	"fmt"
	"io"
	"plugin"
)

type Constructor func() (Command, error)

type Command interface {
	Usage() string
	Args() []Arg
	Flags() []Flag
	Run(ctx context.Context) error
}

type Stdio struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func OpenConstructor(path string) (Constructor, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symbol, err := p.Lookup("New")
	if err != nil {
		return nil, err
	}

	constructor, ok := symbol.(*Constructor)
	if !ok {
		return nil, fmt.Errorf("symbol not a (*kit.Constructor)")
	}

	return *constructor, nil
}
