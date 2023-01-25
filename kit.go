package kit

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hinshun/kit/api"
)

var (
	KitDir         = ".kit"
	ConfigFilename = "config.json"
	ConfigPath     = filepath.Join(KitDir, ConfigFilename)
)

// Command is a kit tool that has first-class access to kit's features.
// For example, the help text and autocomplete is integrated with kit,
// whereas external plugins is just a syscall.Exec.
//
// Kit interacts with Commands over gRPC, but this is abstracted behind this
// interface. The RPC framework is hashicorp/go-plugin which only allows RPC
// over a local network as the purpose is to allow commands to be built
// out-of-tree and developed in other languages.
type Command interface {
	// Usage is a description of the Command.
	Usage() string

	// Args is a list of positional arguments that may be assigned values.
	Args() []Arg

	// Flags is a list of optional flags that may be assigned values.
	Flags() []Flag

	// Run executes the Command.
	Run(ctx context.Context) error
}

type Completion struct {
	Group    string
	Wordlist []string
}

type Stdio struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// Connect starts a gRPC connection to a kit plugin and returns a Command
// that runs over RPC.
func Connect(path string) (RemoteCommand, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: api.HandshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(path),
		Logger:          hclog.NewNullLogger(),
		SyncStdout:      os.Stdout,
		SyncStderr:      os.Stderr,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense("command")
	if err != nil {
		return nil, err
	}

	cmd, ok := raw.(RemoteCommand)
	if !ok {
		return nil, fmt.Errorf("plugin is not of type Command")
	}

	return &commandCloser{
		RemoteCommand: cmd,
		cleanup: client.Kill,
	}, nil
}

// Serve starts a gRPC service to control this Command over a local network.
//
// Serve doesn't return until the plugin is done being executed. Any fixable errors
func Serve(cmd Command) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: api.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"command": &commandPlugin{
				Command: cmd,
			},
		},
		// A non-nil value here enables gRPC serving for this plugin.
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type commandCloser struct {
	RemoteCommand
	cleanup func()
}

func (c *commandCloser) Close() error {
	c.cleanup()
	return nil
}
