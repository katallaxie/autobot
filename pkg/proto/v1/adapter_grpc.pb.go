// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: adapter.proto

package proto

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

// AdapterClient is the client API for Adapter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdapterClient interface {
	// Subscribe ...
	Subscribe(ctx context.Context, in *Subscribe_Request, opts ...grpc.CallOption) (Adapter_SubscribeClient, error)
	// Publish ...
	Publish(ctx context.Context, in *Publish_Request, opts ...grpc.CallOption) (*Publish_Response, error)
}

type adapterClient struct {
	cc grpc.ClientConnInterface
}

func NewAdapterClient(cc grpc.ClientConnInterface) AdapterClient {
	return &adapterClient{cc}
}

func (c *adapterClient) Subscribe(ctx context.Context, in *Subscribe_Request, opts ...grpc.CallOption) (Adapter_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &Adapter_ServiceDesc.Streams[0], "/autobot.adapter.v1.Adapter/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &adapterSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Adapter_SubscribeClient interface {
	Recv() (*Subscribe_Response, error)
	grpc.ClientStream
}

type adapterSubscribeClient struct {
	grpc.ClientStream
}

func (x *adapterSubscribeClient) Recv() (*Subscribe_Response, error) {
	m := new(Subscribe_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *adapterClient) Publish(ctx context.Context, in *Publish_Request, opts ...grpc.CallOption) (*Publish_Response, error) {
	out := new(Publish_Response)
	err := c.cc.Invoke(ctx, "/autobot.adapter.v1.Adapter/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdapterServer is the server API for Adapter service.
// All implementations must embed UnimplementedAdapterServer
// for forward compatibility
type AdapterServer interface {
	// Subscribe ...
	Subscribe(*Subscribe_Request, Adapter_SubscribeServer) error
	// Publish ...
	Publish(context.Context, *Publish_Request) (*Publish_Response, error)
	mustEmbedUnimplementedAdapterServer()
}

// UnimplementedAdapterServer must be embedded to have forward compatible implementations.
type UnimplementedAdapterServer struct {
}

func (UnimplementedAdapterServer) Subscribe(*Subscribe_Request, Adapter_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedAdapterServer) Publish(context.Context, *Publish_Request) (*Publish_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedAdapterServer) mustEmbedUnimplementedAdapterServer() {}

// UnsafeAdapterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdapterServer will
// result in compilation errors.
type UnsafeAdapterServer interface {
	mustEmbedUnimplementedAdapterServer()
}

func RegisterAdapterServer(s grpc.ServiceRegistrar, srv AdapterServer) {
	s.RegisterService(&Adapter_ServiceDesc, srv)
}

func _Adapter_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Subscribe_Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AdapterServer).Subscribe(m, &adapterSubscribeServer{stream})
}

type Adapter_SubscribeServer interface {
	Send(*Subscribe_Response) error
	grpc.ServerStream
}

type adapterSubscribeServer struct {
	grpc.ServerStream
}

func (x *adapterSubscribeServer) Send(m *Subscribe_Response) error {
	return x.ServerStream.SendMsg(m)
}

func _Adapter_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Publish_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdapterServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/autobot.adapter.v1.Adapter/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdapterServer).Publish(ctx, req.(*Publish_Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Adapter_ServiceDesc is the grpc.ServiceDesc for Adapter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Adapter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "autobot.adapter.v1.Adapter",
	HandlerType: (*AdapterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _Adapter_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Adapter_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "adapter.proto",
}
