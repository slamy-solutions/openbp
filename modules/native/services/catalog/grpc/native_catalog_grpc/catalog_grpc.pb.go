// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: catalog.proto

package native_catalog_grpc

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

// CatalogServiceClient is the client API for CatalogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CatalogServiceClient interface {
	// Create new catalog
	Create(ctx context.Context, in *CreateCatalogRequest, opts ...grpc.CallOption) (*CreateCatalogResponse, error)
	// Deletes catalog and all its entries
	Delete(ctx context.Context, in *DeleteCatalogRequest, opts ...grpc.CallOption) (*DeleteCatalogResponse, error)
	// Updates catalog with new information
	Update(ctx context.Context, in *UpdateCatalogRequest, opts ...grpc.CallOption) (*UpdateCatalogReponse, error)
	// Returns catalog by its name
	Get(ctx context.Context, in *GetCatalogRequest, opts ...grpc.CallOption) (*GetCatalogResponse, error)
	// Returns catalog by its name only if provided version differs from the actual. In other case returns NULL. More optimized version, than Get
	GetIfChanged(ctx context.Context, in *GetCatalogIfChangedRequest, opts ...grpc.CallOption) (*GetCatalogIfChangedResponse, error)
	// Streams list of all catalogs
	GetAll(ctx context.Context, in *GetAllCatalogsRequest, opts ...grpc.CallOption) (CatalogService_GetAllClient, error)
	// Lists all indexes in the catalog
	ListIndexes(ctx context.Context, in *ListCatalogIndexesRequest, opts ...grpc.CallOption) (CatalogService_ListIndexesClient, error)
	// Creates or updates index in the catalog
	EnsureIndex(ctx context.Context, in *EnsureCatalogIndexRequest, opts ...grpc.CallOption) (*EnsureCatalogIndexResponse, error)
	// Removes index from the catalog
	RemoveIndex(ctx context.Context, in *RemoveCatalogIndexRequest, opts ...grpc.CallOption) (*RemoveCatalogIndexResponse, error)
}

type catalogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCatalogServiceClient(cc grpc.ClientConnInterface) CatalogServiceClient {
	return &catalogServiceClient{cc}
}

func (c *catalogServiceClient) Create(ctx context.Context, in *CreateCatalogRequest, opts ...grpc.CallOption) (*CreateCatalogResponse, error) {
	out := new(CreateCatalogResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) Delete(ctx context.Context, in *DeleteCatalogRequest, opts ...grpc.CallOption) (*DeleteCatalogResponse, error) {
	out := new(DeleteCatalogResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) Update(ctx context.Context, in *UpdateCatalogRequest, opts ...grpc.CallOption) (*UpdateCatalogReponse, error) {
	out := new(UpdateCatalogReponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) Get(ctx context.Context, in *GetCatalogRequest, opts ...grpc.CallOption) (*GetCatalogResponse, error) {
	out := new(GetCatalogResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) GetIfChanged(ctx context.Context, in *GetCatalogIfChangedRequest, opts ...grpc.CallOption) (*GetCatalogIfChangedResponse, error) {
	out := new(GetCatalogIfChangedResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogService/GetIfChanged", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) GetAll(ctx context.Context, in *GetAllCatalogsRequest, opts ...grpc.CallOption) (CatalogService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &CatalogService_ServiceDesc.Streams[0], "/native_catalog.CatalogService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &catalogServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CatalogService_GetAllClient interface {
	Recv() (*GetAllCatalogsResponse, error)
	grpc.ClientStream
}

type catalogServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *catalogServiceGetAllClient) Recv() (*GetAllCatalogsResponse, error) {
	m := new(GetAllCatalogsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *catalogServiceClient) ListIndexes(ctx context.Context, in *ListCatalogIndexesRequest, opts ...grpc.CallOption) (CatalogService_ListIndexesClient, error) {
	stream, err := c.cc.NewStream(ctx, &CatalogService_ServiceDesc.Streams[1], "/native_catalog.CatalogService/ListIndexes", opts...)
	if err != nil {
		return nil, err
	}
	x := &catalogServiceListIndexesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CatalogService_ListIndexesClient interface {
	Recv() (*ListCatalogIndexesResponse, error)
	grpc.ClientStream
}

type catalogServiceListIndexesClient struct {
	grpc.ClientStream
}

func (x *catalogServiceListIndexesClient) Recv() (*ListCatalogIndexesResponse, error) {
	m := new(ListCatalogIndexesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *catalogServiceClient) EnsureIndex(ctx context.Context, in *EnsureCatalogIndexRequest, opts ...grpc.CallOption) (*EnsureCatalogIndexResponse, error) {
	out := new(EnsureCatalogIndexResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogService/EnsureIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) RemoveIndex(ctx context.Context, in *RemoveCatalogIndexRequest, opts ...grpc.CallOption) (*RemoveCatalogIndexResponse, error) {
	out := new(RemoveCatalogIndexResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogService/RemoveIndex", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CatalogServiceServer is the server API for CatalogService service.
// All implementations must embed UnimplementedCatalogServiceServer
// for forward compatibility
type CatalogServiceServer interface {
	// Create new catalog
	Create(context.Context, *CreateCatalogRequest) (*CreateCatalogResponse, error)
	// Deletes catalog and all its entries
	Delete(context.Context, *DeleteCatalogRequest) (*DeleteCatalogResponse, error)
	// Updates catalog with new information
	Update(context.Context, *UpdateCatalogRequest) (*UpdateCatalogReponse, error)
	// Returns catalog by its name
	Get(context.Context, *GetCatalogRequest) (*GetCatalogResponse, error)
	// Returns catalog by its name only if provided version differs from the actual. In other case returns NULL. More optimized version, than Get
	GetIfChanged(context.Context, *GetCatalogIfChangedRequest) (*GetCatalogIfChangedResponse, error)
	// Streams list of all catalogs
	GetAll(*GetAllCatalogsRequest, CatalogService_GetAllServer) error
	// Lists all indexes in the catalog
	ListIndexes(*ListCatalogIndexesRequest, CatalogService_ListIndexesServer) error
	// Creates or updates index in the catalog
	EnsureIndex(context.Context, *EnsureCatalogIndexRequest) (*EnsureCatalogIndexResponse, error)
	// Removes index from the catalog
	RemoveIndex(context.Context, *RemoveCatalogIndexRequest) (*RemoveCatalogIndexResponse, error)
	mustEmbedUnimplementedCatalogServiceServer()
}

// UnimplementedCatalogServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCatalogServiceServer struct {
}

func (UnimplementedCatalogServiceServer) Create(context.Context, *CreateCatalogRequest) (*CreateCatalogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedCatalogServiceServer) Delete(context.Context, *DeleteCatalogRequest) (*DeleteCatalogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCatalogServiceServer) Update(context.Context, *UpdateCatalogRequest) (*UpdateCatalogReponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedCatalogServiceServer) Get(context.Context, *GetCatalogRequest) (*GetCatalogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedCatalogServiceServer) GetIfChanged(context.Context, *GetCatalogIfChangedRequest) (*GetCatalogIfChangedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIfChanged not implemented")
}
func (UnimplementedCatalogServiceServer) GetAll(*GetAllCatalogsRequest, CatalogService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedCatalogServiceServer) ListIndexes(*ListCatalogIndexesRequest, CatalogService_ListIndexesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListIndexes not implemented")
}
func (UnimplementedCatalogServiceServer) EnsureIndex(context.Context, *EnsureCatalogIndexRequest) (*EnsureCatalogIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnsureIndex not implemented")
}
func (UnimplementedCatalogServiceServer) RemoveIndex(context.Context, *RemoveCatalogIndexRequest) (*RemoveCatalogIndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveIndex not implemented")
}
func (UnimplementedCatalogServiceServer) mustEmbedUnimplementedCatalogServiceServer() {}

// UnsafeCatalogServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CatalogServiceServer will
// result in compilation errors.
type UnsafeCatalogServiceServer interface {
	mustEmbedUnimplementedCatalogServiceServer()
}

func RegisterCatalogServiceServer(s grpc.ServiceRegistrar, srv CatalogServiceServer) {
	s.RegisterService(&CatalogService_ServiceDesc, srv)
}

func _CatalogService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCatalogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServiceServer).Create(ctx, req.(*CreateCatalogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCatalogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServiceServer).Delete(ctx, req.(*DeleteCatalogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCatalogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServiceServer).Update(ctx, req.(*UpdateCatalogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCatalogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServiceServer).Get(ctx, req.(*GetCatalogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogService_GetIfChanged_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCatalogIfChangedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServiceServer).GetIfChanged(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogService/GetIfChanged",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServiceServer).GetIfChanged(ctx, req.(*GetCatalogIfChangedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetAllCatalogsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CatalogServiceServer).GetAll(m, &catalogServiceGetAllServer{stream})
}

type CatalogService_GetAllServer interface {
	Send(*GetAllCatalogsResponse) error
	grpc.ServerStream
}

type catalogServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *catalogServiceGetAllServer) Send(m *GetAllCatalogsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _CatalogService_ListIndexes_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListCatalogIndexesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CatalogServiceServer).ListIndexes(m, &catalogServiceListIndexesServer{stream})
}

type CatalogService_ListIndexesServer interface {
	Send(*ListCatalogIndexesResponse) error
	grpc.ServerStream
}

type catalogServiceListIndexesServer struct {
	grpc.ServerStream
}

func (x *catalogServiceListIndexesServer) Send(m *ListCatalogIndexesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _CatalogService_EnsureIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnsureCatalogIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServiceServer).EnsureIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogService/EnsureIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServiceServer).EnsureIndex(ctx, req.(*EnsureCatalogIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogService_RemoveIndex_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveCatalogIndexRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogServiceServer).RemoveIndex(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogService/RemoveIndex",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogServiceServer).RemoveIndex(ctx, req.(*RemoveCatalogIndexRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CatalogService_ServiceDesc is the grpc.ServiceDesc for CatalogService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CatalogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "native_catalog.CatalogService",
	HandlerType: (*CatalogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _CatalogService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _CatalogService_Delete_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _CatalogService_Update_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _CatalogService_Get_Handler,
		},
		{
			MethodName: "GetIfChanged",
			Handler:    _CatalogService_GetIfChanged_Handler,
		},
		{
			MethodName: "EnsureIndex",
			Handler:    _CatalogService_EnsureIndex_Handler,
		},
		{
			MethodName: "RemoveIndex",
			Handler:    _CatalogService_RemoveIndex_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _CatalogService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListIndexes",
			Handler:       _CatalogService_ListIndexes_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "catalog.proto",
}

// CatalogEntryServiceClient is the client API for CatalogEntryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CatalogEntryServiceClient interface {
	// Creates new entry in the specified catalog. Entry will receive uuid after successfull creation. Returns newly created entry
	Create(ctx context.Context, in *CreateCatalogEntryRequest, opts ...grpc.CallOption) (*CreateCatalogEntryResponse, error)
	// Deletes catalog entry
	Delete(ctx context.Context, in *DeleteCatalogEntryRequest, opts ...grpc.CallOption) (*DeleteCatalogEntryResponse, error)
	// Updates catalog entry with new data
	Update(ctx context.Context, in *UpdateCatalogEntryRequest, opts ...grpc.CallOption) (*UpdateCatalogEntryResponse, error)
	// Get catalog entry. Uses cache and works much faster than Query operation
	Get(ctx context.Context, in *GetCatalogEntryRequest, opts ...grpc.CallOption) (*GetCatalogEntryResponse, error)
	// Run custom query on catalog and get all the entries that satisfy parameters
	Query(ctx context.Context, in *QueryCatalogEntriesRequest, opts ...grpc.CallOption) (*QueryCatalogEntriesResponse, error)
}

type catalogEntryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCatalogEntryServiceClient(cc grpc.ClientConnInterface) CatalogEntryServiceClient {
	return &catalogEntryServiceClient{cc}
}

func (c *catalogEntryServiceClient) Create(ctx context.Context, in *CreateCatalogEntryRequest, opts ...grpc.CallOption) (*CreateCatalogEntryResponse, error) {
	out := new(CreateCatalogEntryResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogEntryService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogEntryServiceClient) Delete(ctx context.Context, in *DeleteCatalogEntryRequest, opts ...grpc.CallOption) (*DeleteCatalogEntryResponse, error) {
	out := new(DeleteCatalogEntryResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogEntryService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogEntryServiceClient) Update(ctx context.Context, in *UpdateCatalogEntryRequest, opts ...grpc.CallOption) (*UpdateCatalogEntryResponse, error) {
	out := new(UpdateCatalogEntryResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogEntryService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogEntryServiceClient) Get(ctx context.Context, in *GetCatalogEntryRequest, opts ...grpc.CallOption) (*GetCatalogEntryResponse, error) {
	out := new(GetCatalogEntryResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogEntryService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogEntryServiceClient) Query(ctx context.Context, in *QueryCatalogEntriesRequest, opts ...grpc.CallOption) (*QueryCatalogEntriesResponse, error) {
	out := new(QueryCatalogEntriesResponse)
	err := c.cc.Invoke(ctx, "/native_catalog.CatalogEntryService/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CatalogEntryServiceServer is the server API for CatalogEntryService service.
// All implementations must embed UnimplementedCatalogEntryServiceServer
// for forward compatibility
type CatalogEntryServiceServer interface {
	// Creates new entry in the specified catalog. Entry will receive uuid after successfull creation. Returns newly created entry
	Create(context.Context, *CreateCatalogEntryRequest) (*CreateCatalogEntryResponse, error)
	// Deletes catalog entry
	Delete(context.Context, *DeleteCatalogEntryRequest) (*DeleteCatalogEntryResponse, error)
	// Updates catalog entry with new data
	Update(context.Context, *UpdateCatalogEntryRequest) (*UpdateCatalogEntryResponse, error)
	// Get catalog entry. Uses cache and works much faster than Query operation
	Get(context.Context, *GetCatalogEntryRequest) (*GetCatalogEntryResponse, error)
	// Run custom query on catalog and get all the entries that satisfy parameters
	Query(context.Context, *QueryCatalogEntriesRequest) (*QueryCatalogEntriesResponse, error)
	mustEmbedUnimplementedCatalogEntryServiceServer()
}

// UnimplementedCatalogEntryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCatalogEntryServiceServer struct {
}

func (UnimplementedCatalogEntryServiceServer) Create(context.Context, *CreateCatalogEntryRequest) (*CreateCatalogEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedCatalogEntryServiceServer) Delete(context.Context, *DeleteCatalogEntryRequest) (*DeleteCatalogEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCatalogEntryServiceServer) Update(context.Context, *UpdateCatalogEntryRequest) (*UpdateCatalogEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedCatalogEntryServiceServer) Get(context.Context, *GetCatalogEntryRequest) (*GetCatalogEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedCatalogEntryServiceServer) Query(context.Context, *QueryCatalogEntriesRequest) (*QueryCatalogEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedCatalogEntryServiceServer) mustEmbedUnimplementedCatalogEntryServiceServer() {}

// UnsafeCatalogEntryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CatalogEntryServiceServer will
// result in compilation errors.
type UnsafeCatalogEntryServiceServer interface {
	mustEmbedUnimplementedCatalogEntryServiceServer()
}

func RegisterCatalogEntryServiceServer(s grpc.ServiceRegistrar, srv CatalogEntryServiceServer) {
	s.RegisterService(&CatalogEntryService_ServiceDesc, srv)
}

func _CatalogEntryService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCatalogEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogEntryServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogEntryService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogEntryServiceServer).Create(ctx, req.(*CreateCatalogEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogEntryService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCatalogEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogEntryServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogEntryService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogEntryServiceServer).Delete(ctx, req.(*DeleteCatalogEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogEntryService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCatalogEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogEntryServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogEntryService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogEntryServiceServer).Update(ctx, req.(*UpdateCatalogEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogEntryService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCatalogEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogEntryServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogEntryService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogEntryServiceServer).Get(ctx, req.(*GetCatalogEntryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CatalogEntryService_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryCatalogEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CatalogEntryServiceServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/native_catalog.CatalogEntryService/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CatalogEntryServiceServer).Query(ctx, req.(*QueryCatalogEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CatalogEntryService_ServiceDesc is the grpc.ServiceDesc for CatalogEntryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CatalogEntryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "native_catalog.CatalogEntryService",
	HandlerType: (*CatalogEntryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _CatalogEntryService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _CatalogEntryService_Delete_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _CatalogEntryService_Update_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _CatalogEntryService_Get_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _CatalogEntryService_Query_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "catalog.proto",
}