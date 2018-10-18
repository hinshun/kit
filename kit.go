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

type Constructor func(cli Cli) Command

type Command interface {
	Name() string
	Usage() string
	Args() []Arg
	Flags() []Flag
	Run(ctx context.Context) error
}

type Arg interface {
}

type Flag interface {
}

// func loadCommands(ctx context.Context) ([]*cli.Command, error) {
// 	cfg, err := ParseConfig(".kit/config.json")
// 	if err != nil {
// 		return nil, err
// 	}

// 	sh := shell.NewLocalShell()
// 	refs, err := ipfs.SyncCommands(ctx, sh, cfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	commandByRef := make(map[string]Command)
// 	for _, ref := range refs {
// 		if _, ok := commandByRef[ref]; ok {
// 			continue
// 		}

// 		command, err := linker.LinkCommand(ref)
// 		if err != nil {
// 			return nil, err
// 		}

// 		commandByRef[ref] = command
// 	}

// 	commands, err := buildCommands(cfg.Plugins, commandByRef)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return commands, nil
// }

// func buildCommands(plugins Plugins, commandByRef map[string]Command) ([]*cli.Command, error) {
// 	var commands []*cli.Command
// 	for _, plugin := range plugins {
// 		switch plugin.Type {
// 		case PluginCommand:
// 			kitCommand, ok := commandByRef[plugin.Ref]
// 			if !ok {
// 				return nil, fmt.Errorf("failed to find command for ref '%s'", plugin.Ref)
// 			}

// 			subCommands, err := buildCommands(plugin.Plugins, commandByRef)
// 			if err != nil {
// 				return nil, err
// 			}

// 			commands = append(commands, &cli.Command{
// 				Name:  plugin.Name,
// 				Usage: kitCommand.Usage(),
// 				Action: func(c *cli.Context) error {
// 					return kitCommand.Action(context.TODO())
// 				},
// 				Commands: subCommands,
// 			})
// 		case PluginManifest:
// 			subCommands, err := buildCommands(plugin.Plugins, commandByRef)
// 			if err != nil {
// 				return nil, err
// 			}

// 			commands = append(commands, &cli.Command{
// 				Name: plugin.Name,
// 				Action: func(c *cli.Context) error {
// 					// args := c.Args()
// 					// if args.Present() {
// 					// 	cli.ShowCommandHelp(c, args.First())
// 					// } else {
// 					// 	cli.ShowAppHelp(c)
// 					// }
// 					return nil
// 				},
// 				Commands: subCommands,
// 			})
// 		default:
// 			return nil, fmt.Errorf("plugin type '%s' for '%s' not implemented", plugin.Type, plugin.Ref)
// 		}
// 	}
// 	return commands, nil
// }
