package kit

import (
	"context"
	"io"
	"path/filepath"
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
