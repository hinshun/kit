package cli

import (
	"context"
	"fmt"

	"github.com/hinshun/kit"
)

func Parse(ctx context.Context, flags []kit.Flag, args []kit.Arg, shellArgs []string) error {
	parsedArgs, err := parse(ctx, flags, shellArgs)
	if err != nil {
		return err
	}

	if len(args) != len(parsedArgs) {
		return fmt.Errorf(
			"requires %d args but got %d args.",
			len(args),
			len(parsedArgs),
		)
	}

	for i, arg := range args {
		err := arg.Set(ctx, parsedArgs[i])
		if err != nil {
			return fmt.Errorf(
				"%s is not a valid %s: %w",
				parsedArgs[i],
				arg.Type(),
				err,
			)
		}
	}

	return nil
}

func parse(ctx context.Context, flags []kit.Flag, shellArgs []string) (args []string, err error) {
	p := &parser{
		flags:      make(map[string]kit.Flag),
		parsedArgs: shellArgs,
	}

	for _, flag := range flags {
		p.flags[flag.Name()] = flag
	}

	for {
		seen, err := p.parseNext(ctx)
		if seen {
			continue
		}
		if err == nil {
			break
		}
		if err != nil {
			return nil, err
		}
		break
	}

	return p.parsedArgs, nil
}

type parser struct {
	flags      map[string]kit.Flag
	parsedArgs []string
}

func (p *parser) parseNext(ctx context.Context) (bool, error) {
	if len(p.parsedArgs) == 0 {
		return false, nil
	}
	s := p.parsedArgs[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}
	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			p.parsedArgs = p.parsedArgs[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, fmt.Errorf("bad flag syntax: %s", s)
	}

	// It's a flag, does it have an argument?
	p.parsedArgs = p.parsedArgs[1:]
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[:i]
			break
		}
	}

	flag, ok := p.flags[name]
	if !ok {
		return false, fmt.Errorf("flag provided but not defined: -%s", name)
	}

	if flag.Type() == "bool" {
		if hasValue {
			if err := flag.Set(ctx, value); err != nil {
				return false, fmt.Errorf("invalid boolean value %q for -%s: %v", value, name, err)
			}
		} else {
			if err := flag.Set(ctx, "true"); err != nil {
				return false, fmt.Errorf("invalid boolean flag %s: %v", name, err)
			}
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(p.parsedArgs) > 0 {
			// value is the next arg
			hasValue = true
			value, p.parsedArgs = p.parsedArgs[0], p.parsedArgs[1:]
		}
		if !hasValue {
			return false, fmt.Errorf("flag needs an argument: -%s", name)
		}
		if err := flag.Set(ctx, value); err != nil {
			return false, fmt.Errorf("invalid value %q for flag -%s: %v", value, name, err)
		}
	}
	return true, nil
}
