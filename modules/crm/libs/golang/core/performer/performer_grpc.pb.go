// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: performer.proto

package crm_performer_grpc

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

// PerformerServiceClient is the client API for PerformerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PerformerServiceClient interface {
	Create(ctx context.Context, in *CreatePerformerRequest, opts ...grpc.CallOption) (*CreatePerformerResponse, error)
	Get(ctx context.Context, in *GetPerformerRequest, opts ...grpc.CallOption) (*GetPerformerResponse, error)
	Update(ctx context.Context, in *UpdatePerformerRequest, opts ...grpc.CallOption) (*UpdatePerformerResponse, error)
	Delete(ctx context.Context, in *DeletePerformerRequest, opts ...grpc.CallOption) (*DeletePerformerResponse, error)
	List(ctx context.Context, in *ListPerformersRequest, opts ...grpc.CallOption) (*ListPerformersResponse, error)
}

type performerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPerformerServiceClient(cc grpc.ClientConnInterface) PerformerServiceClient {
	return &performerServiceClient{cc}
}

func (c *performerServiceClient) Create(ctx context.Context, in *CreatePerformerRequest, opts ...grpc.CallOption) (*CreatePerformerResponse, error) {
	out := new(CreatePerformerResponse)
	err := c.cc.Invoke(ctx, "/crm_performer.PerformerService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *performerServiceClient) Get(ctx context.Context, in *GetPerformerRequest, opts ...grpc.CallOption) (*GetPerformerResponse, error) {
	out := new(GetPerformerResponse)
	err := c.cc.Invoke(ctx, "/crm_performer.PerformerService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *performerServiceClient) Update(ctx context.Context, in *UpdatePerformerRequest, opts ...grpc.CallOption) (*UpdatePerformerResponse, error) {
	out := new(UpdatePerformerResponse)
	err := c.cc.Invoke(ctx, "/crm_performer.PerformerService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *performerServiceClient) Delete(ctx context.Context, in *DeletePerformerRequest, opts ...grpc.CallOption) (*DeletePerformerResponse, error) {
	out := new(DeletePerformerResponse)
	err := c.cc.Invoke(ctx, "/crm_performer.PerformerService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *performerServiceClient) List(ctx context.Context, in *ListPerformersRequest, opts ...grpc.CallOption) (*ListPerformersResponse, error) {
	out := new(ListPerformersResponse)
	err := c.cc.Invoke(ctx, "/crm_performer.PerformerService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PerformerServiceServer is the server API for PerformerService service.
// All implementations must embed UnimplementedPerformerServiceServer
// for forward compatibility
type PerformerServiceServer interface {
	Create(context.Context, *CreatePerformerRequest) (*CreatePerformerResponse, error)
	Get(context.Context, *GetPerformerRequest) (*GetPerformerResponse, error)
	Update(context.Context, *UpdatePerformerRequest) (*UpdatePerformerResponse, error)
	Delete(context.Context, *DeletePerformerRequest) (*DeletePerformerResponse, error)
	List(context.Context, *ListPerformersRequest) (*ListPerformersResponse, error)
	mustEmbedUnimplementedPerformerServiceServer()
}

// UnimplementedPerformerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPerformerServiceServer struct {
}

func (UnimplementedPerformerServiceServer) Create(context.Context, *CreatePerformerRequest) (*CreatePerformerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedPerformerServiceServer) Get(context.Context, *GetPerformerRequest) (*GetPerformerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedPerformerServiceServer) Update(context.Context, *UpdatePerformerRequest) (*UpdatePerformerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedPerformerServiceServer) Delete(context.Context, *DeletePerformerRequest) (*DeletePerformerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedPerformerServiceServer) List(context.Context, *ListPerformersRequest) (*ListPerformersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedPerformerServiceServer) mustEmbedUnimplementedPerformerServiceServer() {}

// UnsafePerformerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PerformerServiceServer will
// result in compilation errors.
type UnsafePerformerServiceServer interface {
	mustEmbedUnimplementedPerformerServiceServer()
}

func RegisterPerformerServiceServer(s grpc.ServiceRegistrar, srv PerformerServiceServer) {
	s.RegisterService(&PerformerService_ServiceDesc, srv)
}

func _PerformerService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePerformerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerformerServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crm_performer.PerformerService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerformerServiceServer).Create(ctx, req.(*CreatePerformerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerformerService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPerformerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerformerServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crm_performer.PerformerService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerformerServiceServer).Get(ctx, req.(*GetPerformerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerformerService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePerformerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerformerServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crm_performer.PerformerService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerformerServiceServer).Update(ctx, req.(*UpdatePerformerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerformerService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePerformerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerformerServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crm_performer.PerformerService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerformerServiceServer).Delete(ctx, req.(*DeletePerformerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerformerService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPerformersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerformerServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/crm_performer.PerformerService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerformerServiceServer).List(ctx, req.(*ListPerformersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PerformerService_ServiceDesc is the grpc.ServiceDesc for PerformerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PerformerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "crm_performer.PerformerService",
	HandlerType: (*PerformerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _PerformerService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _PerformerService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _PerformerService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _PerformerService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _PerformerService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "performer.proto",
}