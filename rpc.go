package kit

import (
	"context"
	"io"
	"fmt"
	"strconv"
	"strings"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hinshun/kit/api"
	"google.golang.org/grpc"
)

var (
	pluginMap = map[string]plugin.Plugin{
		"command": &commandPlugin{},
	}
)

type RemoteCommand interface {
	io.Closer

	Usage() (string, error)

	Args() ([]Arg, error)

	Flags() ([]Flag, error)

	Run(ctx context.Context) error
}

type remoteCommand struct {
	Command
}

func (rc *remoteCommand) Usage() (string, error) {
	return rc.Command.Usage(), nil
}

func (rc *remoteCommand) Args() ([]Arg, error) {
	return rc.Command.Args(), nil
}

func (rc *remoteCommand) Flags() ([]Flag, error) {
	return rc.Command.Flags(), nil
}

func (rc *remoteCommand) Close() error {
	return nil
}

type commandPlugin struct {
	plugin.Plugin
	Command
}

func (cp *commandPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	api.RegisterCommandServer(s, &grpcServer{RemoteCommand: &remoteCommand{cp.Command}})
	return nil
}

func (cp *commandPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{client: api.NewCommandClient(c)}, nil
}

type grpcClient struct {
	client api.CommandClient

	usageResp bool
	usage     string
	flags     []*rpcFlag
	args      []*rpcArg
}

func (gc *grpcClient) Usage() (string, error) {
	if gc.usageResp {
		return gc.usage, nil
	}

	resp, err := gc.client.Usage(context.Background(), &api.Empty{})
	if err != nil {
		return "", err
	}

	gc.usage = resp.Usage

	for _, flag := range resp.Flags {
		gc.flags = append(gc.flags, &rpcFlag{
			client:   gc.client,
			id:       flag.Id,
			name:     flag.Name,
			flagType: flag.Type,
			usage:    flag.Usage,
		})
	}

	for _, arg := range resp.Args {
		gc.args = append(gc.args, &rpcArg{
			client:  gc.client,
			id:      arg.Id,
			argType: arg.Type,
			usage:   arg.Usage,
		})
	}

	gc.usageResp = true
	return gc.usage, nil
}

func (gc *grpcClient) Args() ([]Arg, error) {
	if !gc.usageResp {
		_, err := gc.Usage()
		if err != nil {
			return nil, err
		}
	}

	var args []Arg
	for _, arg := range gc.args {
		args = append(args, arg)
	}
	return args, nil
}

func (gc *grpcClient) Flags() ([]Flag, error) {
	if !gc.usageResp {
		_, err := gc.Usage()
		if err != nil {
			return nil, err
		}
	}

	var flags []Flag
	for _, flag := range gc.flags {
		flags = append(flags, flag)
	}
	return flags, nil
}

func (gc *grpcClient) Run(ctx context.Context) error {
	_, err := gc.client.Run(ctx, &api.RunRequest{})
	return err
}

func (gc *grpcClient) Close() error {
	return nil
}

type rpcArg struct {
	client api.CommandClient

	id      string
	argType string
	usage   string
}

func (a *rpcArg) Type() string {
	return a.argType
}

func (a *rpcArg) Usage() string {
	return a.usage
}

func (a *rpcArg) Set(ctx context.Context, v string) error {
	_, err := a.client.Set(ctx, &api.SetRequest{
		Id:    a.id,
		Value: v,
	})
	return err
}

func (a *rpcArg) Autocomplete(ctx context.Context, input string) ([]Completion, error) {
	resp, err := a.client.Autocomplete(ctx, &api.AutocompleteRequest{
		Id:    a.id,
		Input: input,
	})
	if err != nil {
		return nil, err
	}

	var completions []Completion
	for _, c := range resp.Completions {
		completions = append(completions, Completion{
			Group:    c.Group,
			Wordlist: c.Wordlist,
		})
	}
	return completions, nil
}

type rpcFlag struct {
	client api.CommandClient

	id       string
	name     string
	flagType string
	usage    string
}

func (f *rpcFlag) Name() string {
	return f.name
}

func (f *rpcFlag) Type() string {
	return f.flagType
}

func (f *rpcFlag) Usage() string {
	return f.usage
}

func (f *rpcFlag) Set(ctx context.Context, v string) error {
	_, err := f.client.Set(ctx, &api.SetRequest{
		Id:    f.id,
		Value: v,
	})
	return err
}

func (f *rpcFlag) Autocomplete(ctx context.Context, input string) ([]Completion, error) {
	resp, err := f.client.Autocomplete(ctx, &api.AutocompleteRequest{
		Id:    f.id,
		Input: input,
	})
	if err != nil {
		return nil, err
	}

	var completions []Completion
	for _, c := range resp.Completions {
		completions = append(completions, Completion{
			Group:    c.Group,
			Wordlist: c.Wordlist,
		})
	}
	return completions, nil
}

type grpcServer struct {
	*api.UnimplementedCommandServer
	RemoteCommand
}

func (gs *grpcServer) Usage(ctx context.Context, in *api.Empty) (*api.UsageResponse, error) {
	usage, err := gs.RemoteCommand.Usage()
	if err != nil {
		return nil, err
	}

	cmdFlags, err := gs.RemoteCommand.Flags()
	if err != nil {
		return nil, err
	}

	var flags []*api.Flag
	for i, flag := range cmdFlags {
		flags = append(flags, &api.Flag{
			Id:    fmt.Sprintf("flag/%d", i),
			Name:  flag.Name(),
			Type:  flag.Type(),
			Usage: flag.Usage(),
		})
	}

	cmdArgs, err := gs.RemoteCommand.Args()
	if err != nil {
		return nil, err
	}

	var args []*api.Arg
	for i, arg := range cmdArgs {
		args = append(args, &api.Arg{
			Id:    fmt.Sprintf("arg/%d", i),
			Type:  arg.Type(),
			Usage: arg.Usage(),
		})
	}

	return &api.UsageResponse{
		Usage: usage,
		Flags: flags,
		Args:  args,
	}, nil
}

func (gs *grpcServer) Set(ctx context.Context, in *api.SetRequest) (*api.SetResponse, error) {
	flags, err := gs.RemoteCommand.Flags()
	if err != nil {
		return nil, err
	}

	args, err := gs.RemoteCommand.Args()
	if err != nil {
		return nil, err
	}

	setType, i, err := splitId(in.Id)
	if err != nil {
		return nil, err
	}

	switch setType {
	case "flag":
		err = flags[i].Set(ctx, in.Value)
	case "arg":
		err = args[i].Set(ctx, in.Value)
	default:
		return nil, fmt.Errorf("unrecognized id %q", in.Id)
	}
	if err != nil {
		return nil, err
	}

	return &api.SetResponse{}, nil
}

func (gs *grpcServer) Autocomplete(ctx context.Context, in *api.AutocompleteRequest) (*api.AutocompleteResponse, error) {
	flags, err := gs.RemoteCommand.Flags()
	if err != nil {
		return nil, err
	}

	args, err := gs.RemoteCommand.Args()
	if err != nil {
		return nil, err
	}

	setType, i, err := splitId(in.Id)
	if err != nil {
		return nil, err
	}

	var completions []Completion
	switch setType {
	case "flag":
		completions, err = flags[i].Autocomplete(ctx, in.Input)
	case "arg":
		completions, err = args[i].Autocomplete(ctx, in.Input)
	default:
		return nil, fmt.Errorf("unrecognized id %q", in.Id)
	}
	if err != nil {
		return nil, err
	}

	resp := &api.AutocompleteResponse{}
	for _, c := range completions {
		resp.Completions = append(resp.Completions, &api.Completion{
			Group:    c.Group,
			Wordlist: c.Wordlist,
		})
	}
	return resp, nil
}

func (gs *grpcServer) Run(ctx context.Context, in *api.RunRequest) (*api.RunResponse, error) {
	return &api.RunResponse{}, gs.RemoteCommand.Run(ctx)
}

func splitId(id string) (setType string, index int, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("expected id to have 2 parts but got %d in %q", len(parts), id)
	}
	index, err = strconv.Atoi(parts[1])
	return parts[0], index, err
}
