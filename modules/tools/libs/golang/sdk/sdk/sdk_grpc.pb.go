// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: sdk.proto

package tools_sdk_nsdk_grpc

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

// SDKServiceClient is the client API for SDKService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SDKServiceClient interface {
	DownloadTLSResources(ctx context.Context, in *DownloadTLSResourcesRequest, opts ...grpc.CallOption) (*DownloadTLSResourcesResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	RegisterPublicKeyAsUser(ctx context.Context, in *RegisterPublicKeyAsUserRequest, opts ...grpc.CallOption) (*RegisterPublicKeyAsUserResponse, error)
}

type sDKServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSDKServiceClient(cc grpc.ClientConnInterface) SDKServiceClient {
	return &sDKServiceClient{cc}
}

func (c *sDKServiceClient) DownloadTLSResources(ctx context.Context, in *DownloadTLSResourcesRequest, opts ...grpc.CallOption) (*DownloadTLSResourcesResponse, error) {
	out := new(DownloadTLSResourcesResponse)
	err := c.cc.Invoke(ctx, "/tools_sdk_sdk.SDKService/DownloadTLSResources", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sDKServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/tools_sdk_sdk.SDKService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sDKServiceClient) RegisterPublicKeyAsUser(ctx context.Context, in *RegisterPublicKeyAsUserRequest, opts ...grpc.CallOption) (*RegisterPublicKeyAsUserResponse, error) {
	out := new(RegisterPublicKeyAsUserResponse)
	err := c.cc.Invoke(ctx, "/tools_sdk_sdk.SDKService/RegisterPublicKeyAsUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SDKServiceServer is the server API for SDKService service.
// All implementations must embed UnimplementedSDKServiceServer
// for forward compatibility
type SDKServiceServer interface {
	DownloadTLSResources(context.Context, *DownloadTLSResourcesRequest) (*DownloadTLSResourcesResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	RegisterPublicKeyAsUser(context.Context, *RegisterPublicKeyAsUserRequest) (*RegisterPublicKeyAsUserResponse, error)
	mustEmbedUnimplementedSDKServiceServer()
}

// UnimplementedSDKServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSDKServiceServer struct {
}

func (UnimplementedSDKServiceServer) DownloadTLSResources(context.Context, *DownloadTLSResourcesRequest) (*DownloadTLSResourcesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DownloadTLSResources not implemented")
}
func (UnimplementedSDKServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedSDKServiceServer) RegisterPublicKeyAsUser(context.Context, *RegisterPublicKeyAsUserRequest) (*RegisterPublicKeyAsUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterPublicKeyAsUser not implemented")
}
func (UnimplementedSDKServiceServer) mustEmbedUnimplementedSDKServiceServer() {}

// UnsafeSDKServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SDKServiceServer will
// result in compilation errors.
type UnsafeSDKServiceServer interface {
	mustEmbedUnimplementedSDKServiceServer()
}

func RegisterSDKServiceServer(s grpc.ServiceRegistrar, srv SDKServiceServer) {
	s.RegisterService(&SDKService_ServiceDesc, srv)
}

func _SDKService_DownloadTLSResources_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadTLSResourcesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SDKServiceServer).DownloadTLSResources(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tools_sdk_sdk.SDKService/DownloadTLSResources",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SDKServiceServer).DownloadTLSResources(ctx, req.(*DownloadTLSResourcesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SDKService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SDKServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tools_sdk_sdk.SDKService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SDKServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SDKService_RegisterPublicKeyAsUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterPublicKeyAsUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SDKServiceServer).RegisterPublicKeyAsUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tools_sdk_sdk.SDKService/RegisterPublicKeyAsUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SDKServiceServer).RegisterPublicKeyAsUser(ctx, req.(*RegisterPublicKeyAsUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SDKService_ServiceDesc is the grpc.ServiceDesc for SDKService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SDKService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tools_sdk_sdk.SDKService",
	HandlerType: (*SDKServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DownloadTLSResources",
			Handler:    _SDKService_DownloadTLSResources_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _SDKService_Ping_Handler,
		},
		{
			MethodName: "RegisterPublicKeyAsUser",
			Handler:    _SDKService_RegisterPublicKeyAsUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sdk.proto",
}