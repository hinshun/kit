package kit

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"plugin"
)

var (
	KitDir         = ".kit"
	ConfigPath     = filepath.Join(KitDir, ConfigFilename)
	ConfigFilename = "config.json"
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
