package linker

import (
	"fmt"
	"plugin"

	"github.com/hinshun/kit"
)

func LinkCommand(path string) (kit.Command, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symbol, err := p.Lookup("Command")
	if err != nil {
		return nil, err
	}

	command, ok := symbol.(kit.Command)
	if !ok {
		return nil, fmt.Errorf("symbol not a Command")
	}

	return command, nil
}
