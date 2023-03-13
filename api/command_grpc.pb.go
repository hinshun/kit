// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CommandClient is the client API for Command service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommandClient interface {
	Usage(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*UsageResponse, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Autocomplete(ctx context.Context, in *AutocompleteRequest, opts ...grpc.CallOption) (*AutocompleteResponse, error)
	Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (*RunResponse, error)
}

type commandClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandClient(cc grpc.ClientConnInterface) CommandClient {
	return &commandClient{cc}
}

func (c *commandClient) Usage(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*UsageResponse, error) {
	out := new(UsageResponse)
	err := c.cc.Invoke(ctx, "/proto.Command/Usage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commandClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/proto.Command/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commandClient) Autocomplete(ctx context.Context, in *AutocompleteRequest, opts ...grpc.CallOption) (*AutocompleteResponse, error) {
	out := new(AutocompleteResponse)
	err := c.cc.Invoke(ctx, "/proto.Command/Autocomplete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commandClient) Run(ctx context.Context, in *RunRequest, opts ...grpc.CallOption) (*RunResponse, error) {
	out := new(RunResponse)
	err := c.cc.Invoke(ctx, "/proto.Command/Run", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommandServer is the server API for Command service.
// All implementations must embed UnimplementedCommandServer
// for forward compatibility
type CommandServer interface {
	Usage(context.Context, *Empty) (*UsageResponse, error)
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Autocomplete(context.Context, *AutocompleteRequest) (*AutocompleteResponse, error)
	Run(context.Context, *RunRequest) (*RunResponse, error)
	mustEmbedUnimplementedCommandServer()
}

// UnimplementedCommandServer must be embedded to have forward compatible implementations.
type UnimplementedCommandServer struct {
}

func (UnimplementedCommandServer) Usage(context.Context, *Empty) (*UsageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Usage not implemented")
}
func (UnimplementedCommandServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedCommandServer) Autocomplete(context.Context, *AutocompleteRequest) (*AutocompleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Autocomplete not implemented")
}
func (UnimplementedCommandServer) Run(context.Context, *RunRequest) (*RunResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Run not implemented")
}
func (UnimplementedCommandServer) mustEmbedUnimplementedCommandServer() {}

// UnsafeCommandServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommandServer will
// result in compilation errors.
type UnsafeCommandServer interface {
	mustEmbedUnimplementedCommandServer()
}

func RegisterCommandServer(s grpc.ServiceRegistrar, srv CommandServer) {
	s.RegisterService(&Command_ServiceDesc, srv)
}

func _Command_Usage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandServer).Usage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Command/Usage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandServer).Usage(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Command_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Command/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Command_Autocomplete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AutocompleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandServer).Autocomplete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Command/Autocomplete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandServer).Autocomplete(ctx, req.(*AutocompleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Command_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Command/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandServer).Run(ctx, req.(*RunRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Command_ServiceDesc is the grpc.ServiceDesc for Command service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Command_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Command",
	HandlerType: (*CommandServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Usage",
			Handler:    _Command_Usage_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _Command_Set_Handler,
		},
		{
			MethodName: "Autocomplete",
			Handler:    _Command_Autocomplete_Handler,
		},
		{
			MethodName: "Run",
			Handler:    _Command_Run_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "command.proto",
}