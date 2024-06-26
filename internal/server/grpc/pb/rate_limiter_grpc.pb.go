// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: rate_limiter.proto

package pb

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

const (
	RateLimiter_Allow_FullMethodName      = "/limiter.RateLimiter/Allow"
	RateLimiter_AddNet_FullMethodName     = "/limiter.RateLimiter/AddNet"
	RateLimiter_DeleteNet_FullMethodName  = "/limiter.RateLimiter/DeleteNet"
	RateLimiter_DropBucket_FullMethodName = "/limiter.RateLimiter/DropBucket"
)

// RateLimiterClient is the client API for RateLimiter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RateLimiterClient interface {
	Allow(ctx context.Context, in *AllowRequest, opts ...grpc.CallOption) (*AllowResponse, error)
	AddNet(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*Response, error)
	DeleteNet(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Response, error)
	DropBucket(ctx context.Context, in *DropRequest, opts ...grpc.CallOption) (*Response, error)
}

type rateLimiterClient struct {
	cc grpc.ClientConnInterface
}

func NewRateLimiterClient(cc grpc.ClientConnInterface) RateLimiterClient {
	return &rateLimiterClient{cc}
}

func (c *rateLimiterClient) Allow(ctx context.Context, in *AllowRequest, opts ...grpc.CallOption) (*AllowResponse, error) {
	out := new(AllowResponse)
	err := c.cc.Invoke(ctx, RateLimiter_Allow_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rateLimiterClient) AddNet(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, RateLimiter_AddNet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rateLimiterClient) DeleteNet(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, RateLimiter_DeleteNet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rateLimiterClient) DropBucket(ctx context.Context, in *DropRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, RateLimiter_DropBucket_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RateLimiterServer is the server API for RateLimiter service.
// All implementations must embed UnimplementedRateLimiterServer
// for forward compatibility
type RateLimiterServer interface {
	Allow(context.Context, *AllowRequest) (*AllowResponse, error)
	AddNet(context.Context, *AddRequest) (*Response, error)
	DeleteNet(context.Context, *DeleteRequest) (*Response, error)
	DropBucket(context.Context, *DropRequest) (*Response, error)
	mustEmbedUnimplementedRateLimiterServer()
}

// UnimplementedRateLimiterServer must be embedded to have forward compatible implementations.
type UnimplementedRateLimiterServer struct {
}

func (UnimplementedRateLimiterServer) Allow(context.Context, *AllowRequest) (*AllowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Allow not implemented")
}
func (UnimplementedRateLimiterServer) AddNet(context.Context, *AddRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNet not implemented")
}
func (UnimplementedRateLimiterServer) DeleteNet(context.Context, *DeleteRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNet not implemented")
}
func (UnimplementedRateLimiterServer) DropBucket(context.Context, *DropRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DropBucket not implemented")
}
func (UnimplementedRateLimiterServer) mustEmbedUnimplementedRateLimiterServer() {}

// UnsafeRateLimiterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RateLimiterServer will
// result in compilation errors.
type UnsafeRateLimiterServer interface {
	mustEmbedUnimplementedRateLimiterServer()
}

func RegisterRateLimiterServer(s grpc.ServiceRegistrar, srv RateLimiterServer) {
	s.RegisterService(&RateLimiter_ServiceDesc, srv)
}

func _RateLimiter_Allow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllowRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServer).Allow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiter_Allow_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServer).Allow(ctx, req.(*AllowRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RateLimiter_AddNet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServer).AddNet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiter_AddNet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServer).AddNet(ctx, req.(*AddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RateLimiter_DeleteNet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServer).DeleteNet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiter_DeleteNet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServer).DeleteNet(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RateLimiter_DropBucket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimiterServer).DropBucket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimiter_DropBucket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimiterServer).DropBucket(ctx, req.(*DropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RateLimiter_ServiceDesc is the grpc.ServiceDesc for RateLimiter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RateLimiter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "limiter.RateLimiter",
	HandlerType: (*RateLimiterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Allow",
			Handler:    _RateLimiter_Allow_Handler,
		},
		{
			MethodName: "AddNet",
			Handler:    _RateLimiter_AddNet_Handler,
		},
		{
			MethodName: "DeleteNet",
			Handler:    _RateLimiter_DeleteNet_Handler,
		},
		{
			MethodName: "DropBucket",
			Handler:    _RateLimiter_DropBucket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rate_limiter.proto",
}
