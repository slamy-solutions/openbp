// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.12.4
// source: cache.proto

package native_cache_grpc

import (
	_struct "github.com/golang/protobuf/ptypes/struct"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CacheEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique key
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Data to cache
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// Remove data from the cache after specified period of time in miliseconds. 0 to disable
	Timeout uint64 `protobuf:"varint,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *CacheEntry) Reset() {
	*x = CacheEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CacheEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CacheEntry) ProtoMessage() {}

func (x *CacheEntry) ProtoReflect() protoreflect.Message {
	mi := &file_cache_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CacheEntry.ProtoReflect.Descriptor instead.
func (*CacheEntry) Descriptor() ([]byte, []int) {
	return file_cache_proto_rawDescGZIP(), []int{0}
}

func (x *CacheEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *CacheEntry) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *CacheEntry) GetTimeout() uint64 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

type CacheSetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Cache entries to set
	Entries []*CacheEntry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *CacheSetRequest) Reset() {
	*x = CacheSetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CacheSetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CacheSetRequest) ProtoMessage() {}

func (x *CacheSetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CacheSetRequest.ProtoReflect.Descriptor instead.
func (*CacheSetRequest) Descriptor() ([]byte, []int) {
	return file_cache_proto_rawDescGZIP(), []int{1}
}

func (x *CacheSetRequest) GetEntries() []*CacheEntry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type CacheSetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CacheSetResponse) Reset() {
	*x = CacheSetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CacheSetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CacheSetResponse) ProtoMessage() {}

func (x *CacheSetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CacheSetResponse.ProtoReflect.Descriptor instead.
func (*CacheSetResponse) Descriptor() ([]byte, []int) {
	return file_cache_proto_rawDescGZIP(), []int{2}
}

type CacheRemoveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Keys to remove
	Keys []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *CacheRemoveRequest) Reset() {
	*x = CacheRemoveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CacheRemoveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CacheRemoveRequest) ProtoMessage() {}

func (x *CacheRemoveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CacheRemoveRequest.ProtoReflect.Descriptor instead.
func (*CacheRemoveRequest) Descriptor() ([]byte, []int) {
	return file_cache_proto_rawDescGZIP(), []int{3}
}

func (x *CacheRemoveRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type CacheRemoveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CacheRemoveResponse) Reset() {
	*x = CacheRemoveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CacheRemoveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CacheRemoveResponse) ProtoMessage() {}

func (x *CacheRemoveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CacheRemoveResponse.ProtoReflect.Descriptor instead.
func (*CacheRemoveResponse) Descriptor() ([]byte, []int) {
	return file_cache_proto_rawDescGZIP(), []int{4}
}

type CacheGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Keys to get
	Keys []string `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
}

func (x *CacheGetRequest) Reset() {
	*x = CacheGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CacheGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CacheGetRequest) ProtoMessage() {}

func (x *CacheGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_cache_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CacheGetRequest.ProtoReflect.Descriptor instead.
func (*CacheGetRequest) Descriptor() ([]byte, []int) {
	return file_cache_proto_rawDescGZIP(), []int{5}
}

func (x *CacheGetRequest) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

type CacheGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Datas []*CacheGetResponse_Data `protobuf:"bytes,1,rep,name=datas,proto3" json:"datas,omitempty"`
}

func (x *CacheGetResponse) Reset() {
	*x = CacheGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CacheGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CacheGetResponse) ProtoMessage() {}

func (x *CacheGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cache_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CacheGetResponse.ProtoReflect.Descriptor instead.
func (*CacheGetResponse) Descriptor() ([]byte, []int) {
	return file_cache_proto_rawDescGZIP(), []int{6}
}

func (x *CacheGetResponse) GetDatas() []*CacheGetResponse_Data {
	if x != nil {
		return x.Datas
	}
	return nil
}

type CacheGetResponse_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//	*CacheGetResponse_Data_Null
	//	*CacheGetResponse_Data_ByteData
	Data isCacheGetResponse_Data_Data `protobuf_oneof:"data"`
}

func (x *CacheGetResponse_Data) Reset() {
	*x = CacheGetResponse_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CacheGetResponse_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CacheGetResponse_Data) ProtoMessage() {}

func (x *CacheGetResponse_Data) ProtoReflect() protoreflect.Message {
	mi := &file_cache_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CacheGetResponse_Data.ProtoReflect.Descriptor instead.
func (*CacheGetResponse_Data) Descriptor() ([]byte, []int) {
	return file_cache_proto_rawDescGZIP(), []int{6, 0}
}

func (m *CacheGetResponse_Data) GetData() isCacheGetResponse_Data_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *CacheGetResponse_Data) GetNull() _struct.NullValue {
	if x, ok := x.GetData().(*CacheGetResponse_Data_Null); ok {
		return x.Null
	}
	return _struct.NullValue(0)
}

func (x *CacheGetResponse_Data) GetByteData() []byte {
	if x, ok := x.GetData().(*CacheGetResponse_Data_ByteData); ok {
		return x.ByteData
	}
	return nil
}

type isCacheGetResponse_Data_Data interface {
	isCacheGetResponse_Data_Data()
}

type CacheGetResponse_Data_Null struct {
	Null _struct.NullValue `protobuf:"varint,1,opt,name=null,proto3,enum=google.protobuf.NullValue,oneof"`
}

type CacheGetResponse_Data_ByteData struct {
	ByteData []byte `protobuf:"bytes,2,opt,name=byteData,proto3,oneof"`
}

func (*CacheGetResponse_Data_Null) isCacheGetResponse_Data_Data() {}

func (*CacheGetResponse_Data_ByteData) isCacheGetResponse_Data_Data() {}

var File_cache_proto protoreflect.FileDescriptor

var file_cache_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6e,
	0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4c, 0x0a, 0x0a, 0x43, 0x61, 0x63,
	0x68, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a,
	0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07,
	0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x45, 0x0a, 0x0f, 0x43, 0x61, 0x63, 0x68, 0x65,
	0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x07, 0x65, 0x6e,
	0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6e, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x43, 0x61, 0x63, 0x68, 0x65,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x12,
	0x0a, 0x10, 0x43, 0x61, 0x63, 0x68, 0x65, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x28, 0x0a, 0x12, 0x43, 0x61, 0x63, 0x68, 0x65, 0x52, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0x15, 0x0a, 0x13,
	0x43, 0x61, 0x63, 0x68, 0x65, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x25, 0x0a, 0x0f, 0x43, 0x61, 0x63, 0x68, 0x65, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x22, 0xad, 0x01, 0x0a, 0x10, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x39, 0x0a, 0x05, 0x64, 0x61, 0x74, 0x61, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23,
	0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x43, 0x61,
	0x63, 0x68, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x05, 0x64, 0x61, 0x74, 0x61, 0x73, 0x1a, 0x5e, 0x0a, 0x04, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x75, 0x6c, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4e, 0x75, 0x6c, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x04,
	0x6e, 0x75, 0x6c, 0x6c, 0x12, 0x1c, 0x0a, 0x08, 0x62, 0x79, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x08, 0x62, 0x79, 0x74, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xef, 0x01, 0x0a, 0x0c, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x46, 0x0a, 0x03, 0x53,
	0x65, 0x74, 0x12, 0x1d, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68,
	0x65, 0x2e, 0x43, 0x61, 0x63, 0x68, 0x65, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1e, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65,
	0x2e, 0x43, 0x61, 0x63, 0x68, 0x65, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x06, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x12, 0x20, 0x2e,
	0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x43, 0x61, 0x63,
	0x68, 0x65, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x21, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x43,
	0x61, 0x63, 0x68, 0x65, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x1d, 0x2e, 0x6e, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x43, 0x61, 0x63, 0x68, 0x65,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x6e, 0x61, 0x74,
	0x69, 0x76, 0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x43, 0x61, 0x63, 0x68, 0x65, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c,
	0x73, 0x6c, 0x61, 0x6d, 0x79, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x43, 0x52, 0x4d, 0x2f, 0x6e, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x2f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x3b, 0x6e, 0x61, 0x74, 0x69, 0x76,
	0x65, 0x5f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cache_proto_rawDescOnce sync.Once
	file_cache_proto_rawDescData = file_cache_proto_rawDesc
)

func file_cache_proto_rawDescGZIP() []byte {
	file_cache_proto_rawDescOnce.Do(func() {
		file_cache_proto_rawDescData = protoimpl.X.CompressGZIP(file_cache_proto_rawDescData)
	})
	return file_cache_proto_rawDescData
}

var file_cache_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_cache_proto_goTypes = []interface{}{
	(*CacheEntry)(nil),            // 0: native_cache.CacheEntry
	(*CacheSetRequest)(nil),       // 1: native_cache.CacheSetRequest
	(*CacheSetResponse)(nil),      // 2: native_cache.CacheSetResponse
	(*CacheRemoveRequest)(nil),    // 3: native_cache.CacheRemoveRequest
	(*CacheRemoveResponse)(nil),   // 4: native_cache.CacheRemoveResponse
	(*CacheGetRequest)(nil),       // 5: native_cache.CacheGetRequest
	(*CacheGetResponse)(nil),      // 6: native_cache.CacheGetResponse
	(*CacheGetResponse_Data)(nil), // 7: native_cache.CacheGetResponse.Data
	(_struct.NullValue)(0),        // 8: google.protobuf.NullValue
}
var file_cache_proto_depIdxs = []int32{
	0, // 0: native_cache.CacheSetRequest.entries:type_name -> native_cache.CacheEntry
	7, // 1: native_cache.CacheGetResponse.datas:type_name -> native_cache.CacheGetResponse.Data
	8, // 2: native_cache.CacheGetResponse.Data.null:type_name -> google.protobuf.NullValue
	1, // 3: native_cache.CacheService.Set:input_type -> native_cache.CacheSetRequest
	3, // 4: native_cache.CacheService.Remove:input_type -> native_cache.CacheRemoveRequest
	5, // 5: native_cache.CacheService.Get:input_type -> native_cache.CacheGetRequest
	2, // 6: native_cache.CacheService.Set:output_type -> native_cache.CacheSetResponse
	4, // 7: native_cache.CacheService.Remove:output_type -> native_cache.CacheRemoveResponse
	6, // 8: native_cache.CacheService.Get:output_type -> native_cache.CacheGetResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_cache_proto_init() }
func file_cache_proto_init() {
	if File_cache_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cache_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CacheEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cache_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CacheSetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cache_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CacheSetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cache_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CacheRemoveRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cache_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CacheRemoveResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cache_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CacheGetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cache_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CacheGetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cache_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CacheGetResponse_Data); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_cache_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*CacheGetResponse_Data_Null)(nil),
		(*CacheGetResponse_Data_ByteData)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cache_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cache_proto_goTypes,
		DependencyIndexes: file_cache_proto_depIdxs,
		MessageInfos:      file_cache_proto_msgTypes,
	}.Build()
	File_cache_proto = out.File
	file_cache_proto_rawDesc = nil
	file_cache_proto_goTypes = nil
	file_cache_proto_depIdxs = nil
}