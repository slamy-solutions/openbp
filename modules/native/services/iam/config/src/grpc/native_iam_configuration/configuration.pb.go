// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.12.4
// source: configuration.proto

package native_iam_configuration_grpc

import (
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

// Configuration of the IAM
type Configuration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Time to live of access token in milliseconds
	AccessTokenTTL uint32 `protobuf:"varint,1,opt,name=accessTokenTTL,proto3" json:"accessTokenTTL,omitempty"`
	// Time to live ot refresh token in milliseconds
	RefreshTokenTTL uint32 `protobuf:"varint,2,opt,name=refreshTokenTTL,proto3" json:"refreshTokenTTL,omitempty"`
	// Password authentication configuration
	PasswordAuth *Configuration_PasswordAuth `protobuf:"bytes,10,opt,name=passwordAuth,proto3" json:"passwordAuth,omitempty"`
	// Google oauth2 configuration
	GoogleOAuth2 *Configuration_OAuth2 `protobuf:"bytes,11,opt,name=googleOAuth2,proto3" json:"googleOAuth2,omitempty"`
	// Facebook oauth2 configuration
	FacebookOAuth2 *Configuration_OAuth2 `protobuf:"bytes,12,opt,name=facebookOAuth2,proto3" json:"facebookOAuth2,omitempty"`
	// Github oauth2 configuration
	GithubOAuth2 *Configuration_OAuth2 `protobuf:"bytes,13,opt,name=githubOAuth2,proto3" json:"githubOAuth2,omitempty"`
	// Github oauth2 configuration
	GitlabOAuth2 *Configuration_OAuth2 `protobuf:"bytes,14,opt,name=gitlabOAuth2,proto3" json:"gitlabOAuth2,omitempty"`
}

func (x *Configuration) Reset() {
	*x = Configuration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration) ProtoMessage() {}

func (x *Configuration) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration.ProtoReflect.Descriptor instead.
func (*Configuration) Descriptor() ([]byte, []int) {
	return file_configuration_proto_rawDescGZIP(), []int{0}
}

func (x *Configuration) GetAccessTokenTTL() uint32 {
	if x != nil {
		return x.AccessTokenTTL
	}
	return 0
}

func (x *Configuration) GetRefreshTokenTTL() uint32 {
	if x != nil {
		return x.RefreshTokenTTL
	}
	return 0
}

func (x *Configuration) GetPasswordAuth() *Configuration_PasswordAuth {
	if x != nil {
		return x.PasswordAuth
	}
	return nil
}

func (x *Configuration) GetGoogleOAuth2() *Configuration_OAuth2 {
	if x != nil {
		return x.GoogleOAuth2
	}
	return nil
}

func (x *Configuration) GetFacebookOAuth2() *Configuration_OAuth2 {
	if x != nil {
		return x.FacebookOAuth2
	}
	return nil
}

func (x *Configuration) GetGithubOAuth2() *Configuration_OAuth2 {
	if x != nil {
		return x.GithubOAuth2
	}
	return nil
}

func (x *Configuration) GetGitlabOAuth2() *Configuration_OAuth2 {
	if x != nil {
		return x.GitlabOAuth2
	}
	return nil
}

type GetConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Use cache or not. Cache have a very low chance to be invalid. Cache invalidates after short period of thime (60 seconds). Cache can only be invalid on multiple simultanious read and writes. Its safe to use cache in most of the cases.
	UseCache bool `protobuf:"varint,1,opt,name=useCache,proto3" json:"useCache,omitempty"`
}

func (x *GetConfigRequest) Reset() {
	*x = GetConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConfigRequest) ProtoMessage() {}

func (x *GetConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConfigRequest.ProtoReflect.Descriptor instead.
func (*GetConfigRequest) Descriptor() ([]byte, []int) {
	return file_configuration_proto_rawDescGZIP(), []int{1}
}

func (x *GetConfigRequest) GetUseCache() bool {
	if x != nil {
		return x.UseCache
	}
	return false
}

type GetConfigresponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Current configuration
	Configuration *Configuration `protobuf:"bytes,1,opt,name=configuration,proto3" json:"configuration,omitempty"`
}

func (x *GetConfigresponse) Reset() {
	*x = GetConfigresponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetConfigresponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConfigresponse) ProtoMessage() {}

func (x *GetConfigresponse) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConfigresponse.ProtoReflect.Descriptor instead.
func (*GetConfigresponse) Descriptor() ([]byte, []int) {
	return file_configuration_proto_rawDescGZIP(), []int{2}
}

func (x *GetConfigresponse) GetConfiguration() *Configuration {
	if x != nil {
		return x.Configuration
	}
	return nil
}

type SetConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Configuration to set
	Configuration *Configuration `protobuf:"bytes,1,opt,name=configuration,proto3" json:"configuration,omitempty"`
}

func (x *SetConfigRequest) Reset() {
	*x = SetConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetConfigRequest) ProtoMessage() {}

func (x *SetConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetConfigRequest.ProtoReflect.Descriptor instead.
func (*SetConfigRequest) Descriptor() ([]byte, []int) {
	return file_configuration_proto_rawDescGZIP(), []int{3}
}

func (x *SetConfigRequest) GetConfiguration() *Configuration {
	if x != nil {
		return x.Configuration
	}
	return nil
}

type SetConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SetConfigResponse) Reset() {
	*x = SetConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetConfigResponse) ProtoMessage() {}

func (x *SetConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetConfigResponse.ProtoReflect.Descriptor instead.
func (*SetConfigResponse) Descriptor() ([]byte, []int) {
	return file_configuration_proto_rawDescGZIP(), []int{4}
}

// Configuration of specific OAuth2 provider
type Configuration_OAuth2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Enable or disable this provider of OAuth2
	Enabled bool `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// OAuth2 client ID
	ClientId string `protobuf:"bytes,2,opt,name=clientId,proto3" json:"clientId,omitempty"`
	// OAuth2 client secret
	ClientSecret string `protobuf:"bytes,3,opt,name=clientSecret,proto3" json:"clientSecret,omitempty"`
	// Allow registration using this OAuth2 provider
	AllowRegistration bool `protobuf:"varint,4,opt,name=allowRegistration,proto3" json:"allowRegistration,omitempty"`
}

func (x *Configuration_OAuth2) Reset() {
	*x = Configuration_OAuth2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration_OAuth2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration_OAuth2) ProtoMessage() {}

func (x *Configuration_OAuth2) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration_OAuth2.ProtoReflect.Descriptor instead.
func (*Configuration_OAuth2) Descriptor() ([]byte, []int) {
	return file_configuration_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Configuration_OAuth2) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *Configuration_OAuth2) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *Configuration_OAuth2) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

func (x *Configuration_OAuth2) GetAllowRegistration() bool {
	if x != nil {
		return x.AllowRegistration
	}
	return false
}

type Configuration_PasswordAuth struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Allow password authorization or not
	Enabled bool `protobuf:"varint,1,opt,name=enabled,proto3" json:"enabled,omitempty"`
	// Allow registration using password method
	AllowRegistration bool `protobuf:"varint,2,opt,name=allowRegistration,proto3" json:"allowRegistration,omitempty"`
}

func (x *Configuration_PasswordAuth) Reset() {
	*x = Configuration_PasswordAuth{}
	if protoimpl.UnsafeEnabled {
		mi := &file_configuration_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Configuration_PasswordAuth) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Configuration_PasswordAuth) ProtoMessage() {}

func (x *Configuration_PasswordAuth) ProtoReflect() protoreflect.Message {
	mi := &file_configuration_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Configuration_PasswordAuth.ProtoReflect.Descriptor instead.
func (*Configuration_PasswordAuth) Descriptor() ([]byte, []int) {
	return file_configuration_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Configuration_PasswordAuth) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *Configuration_PasswordAuth) GetAllowRegistration() bool {
	if x != nil {
		return x.AllowRegistration
	}
	return false
}

var File_configuration_proto protoreflect.FileDescriptor

var file_configuration_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61,
	0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0xfa, 0x05, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x26, 0x0a, 0x0e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x54, 0x54, 0x4c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x54, 0x4c, 0x12, 0x28, 0x0a, 0x0f, 0x72, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x54, 0x4c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0f, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x54, 0x54, 0x4c, 0x12, 0x58, 0x0a, 0x0c, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x41,
	0x75, 0x74, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x6e, 0x61, 0x74, 0x69,
	0x76, 0x65, 0x5f, 0x69, 0x61, 0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x41, 0x75, 0x74, 0x68, 0x52,
	0x0c, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x41, 0x75, 0x74, 0x68, 0x12, 0x52, 0x0a,
	0x0c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61, 0x6d,
	0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x41, 0x75,
	0x74, 0x68, 0x32, 0x52, 0x0c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x4f, 0x41, 0x75, 0x74, 0x68,
	0x32, 0x12, 0x56, 0x0a, 0x0e, 0x66, 0x61, 0x63, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x4f, 0x41, 0x75,
	0x74, 0x68, 0x32, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6e, 0x61, 0x74, 0x69,
	0x76, 0x65, 0x5f, 0x69, 0x61, 0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x52, 0x0e, 0x66, 0x61, 0x63, 0x65, 0x62,
	0x6f, 0x6f, 0x6b, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x12, 0x52, 0x0a, 0x0c, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2e, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61, 0x6d, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x52,
	0x0c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x12, 0x52, 0x0a,
	0x0c, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61, 0x6d,
	0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x4f, 0x41, 0x75,
	0x74, 0x68, 0x32, 0x52, 0x0c, 0x67, 0x69, 0x74, 0x6c, 0x61, 0x62, 0x4f, 0x41, 0x75, 0x74, 0x68,
	0x32, 0x1a, 0x90, 0x01, 0x0a, 0x06, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32, 0x12, 0x18, 0x0a, 0x07,
	0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65,
	0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x2c, 0x0a, 0x11, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x11, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x56, 0x0a, 0x0c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x41, 0x75, 0x74, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x2c,
	0x0a, 0x11, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x61, 0x6c, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2e, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x43, 0x61, 0x63, 0x68, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x75, 0x73, 0x65, 0x43, 0x61, 0x63, 0x68, 0x65, 0x22, 0x62, 0x0a, 0x11,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x4d, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76,
	0x65, 0x5f, 0x69, 0x61, 0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0x61, 0x0a, 0x10, 0x53, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x4d, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6e, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61, 0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x22, 0x13, 0x0a, 0x11, 0x53, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xd2, 0x01, 0x0a, 0x10, 0x49, 0x41, 0x4d,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5e, 0x0a,
	0x03, 0x47, 0x65, 0x74, 0x12, 0x2a, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61,
	0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2b, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61, 0x6d, 0x5f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5e, 0x0a,
	0x03, 0x53, 0x65, 0x74, 0x12, 0x2a, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61,
	0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x53, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2b, 0x2e, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69, 0x61, 0x6d, 0x5f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x74, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x46, 0x5a,
	0x44, 0x73, 0x6c, 0x61, 0x6d, 0x79, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x43, 0x52, 0x4d, 0x2f, 0x6e,
	0x61, 0x74, 0x69, 0x76, 0x65, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x3b, 0x6e, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x69,
	0x61, 0x6d, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_configuration_proto_rawDescOnce sync.Once
	file_configuration_proto_rawDescData = file_configuration_proto_rawDesc
)

func file_configuration_proto_rawDescGZIP() []byte {
	file_configuration_proto_rawDescOnce.Do(func() {
		file_configuration_proto_rawDescData = protoimpl.X.CompressGZIP(file_configuration_proto_rawDescData)
	})
	return file_configuration_proto_rawDescData
}

var file_configuration_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_configuration_proto_goTypes = []interface{}{
	(*Configuration)(nil),              // 0: native_iam_configuration.Configuration
	(*GetConfigRequest)(nil),           // 1: native_iam_configuration.GetConfigRequest
	(*GetConfigresponse)(nil),          // 2: native_iam_configuration.GetConfigresponse
	(*SetConfigRequest)(nil),           // 3: native_iam_configuration.SetConfigRequest
	(*SetConfigResponse)(nil),          // 4: native_iam_configuration.SetConfigResponse
	(*Configuration_OAuth2)(nil),       // 5: native_iam_configuration.Configuration.OAuth2
	(*Configuration_PasswordAuth)(nil), // 6: native_iam_configuration.Configuration.PasswordAuth
}
var file_configuration_proto_depIdxs = []int32{
	6, // 0: native_iam_configuration.Configuration.passwordAuth:type_name -> native_iam_configuration.Configuration.PasswordAuth
	5, // 1: native_iam_configuration.Configuration.googleOAuth2:type_name -> native_iam_configuration.Configuration.OAuth2
	5, // 2: native_iam_configuration.Configuration.facebookOAuth2:type_name -> native_iam_configuration.Configuration.OAuth2
	5, // 3: native_iam_configuration.Configuration.githubOAuth2:type_name -> native_iam_configuration.Configuration.OAuth2
	5, // 4: native_iam_configuration.Configuration.gitlabOAuth2:type_name -> native_iam_configuration.Configuration.OAuth2
	0, // 5: native_iam_configuration.GetConfigresponse.configuration:type_name -> native_iam_configuration.Configuration
	0, // 6: native_iam_configuration.SetConfigRequest.configuration:type_name -> native_iam_configuration.Configuration
	1, // 7: native_iam_configuration.IAMConfigService.Get:input_type -> native_iam_configuration.GetConfigRequest
	3, // 8: native_iam_configuration.IAMConfigService.Set:input_type -> native_iam_configuration.SetConfigRequest
	2, // 9: native_iam_configuration.IAMConfigService.Get:output_type -> native_iam_configuration.GetConfigresponse
	4, // 10: native_iam_configuration.IAMConfigService.Set:output_type -> native_iam_configuration.SetConfigResponse
	9, // [9:11] is the sub-list for method output_type
	7, // [7:9] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_configuration_proto_init() }
func file_configuration_proto_init() {
	if File_configuration_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_configuration_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration); i {
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
		file_configuration_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetConfigRequest); i {
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
		file_configuration_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetConfigresponse); i {
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
		file_configuration_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetConfigRequest); i {
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
		file_configuration_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetConfigResponse); i {
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
		file_configuration_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration_OAuth2); i {
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
		file_configuration_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Configuration_PasswordAuth); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_configuration_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_configuration_proto_goTypes,
		DependencyIndexes: file_configuration_proto_depIdxs,
		MessageInfos:      file_configuration_proto_msgTypes,
	}.Build()
	File_configuration_proto = out.File
	file_configuration_proto_rawDesc = nil
	file_configuration_proto_goTypes = nil
	file_configuration_proto_depIdxs = nil
}