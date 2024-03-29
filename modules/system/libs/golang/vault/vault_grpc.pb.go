// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: vault.proto

package vault

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

// VaultServiceClient is the client API for VaultService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VaultServiceClient interface {
	// Close and encrypt vault. After sealing, most of the operations will not be accessible.
	Seal(ctx context.Context, in *SealRequest, opts ...grpc.CallOption) (*SealResponse, error)
	// Decrypt and open vault. Must be done before most of the operations with vault secrets.
	Unseal(ctx context.Context, in *UnsealRequest, opts ...grpc.CallOption) (*UnsealResponse, error)
	// Set up new seal secret and reincrypt vault. The vault must be unsealed before this operation. You don't need to unseal vault after this operation.
	// This operation requires you to have administrator access to the HSM. Check PKCS11 spec. If you are using emulated HSM (by default) this will be the same as the seal/unseal secret by default ("12345678"). Change it.
	UpdateSealSecret(ctx context.Context, in *UpdateSealSecretRequest, opts ...grpc.CallOption) (*UpdateSealSecretResponse, error)
	// Get current status of the vault.
	GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error)
	// Creates RSA key pair if it doesnt exist. Private key never leaves the HSM (hardware security module).
	EnsureRSAKeyPair(ctx context.Context, in *EnsureRSAKeyPairRequest, opts ...grpc.CallOption) (*EnsureRSAKeyPairResponse, error)
	// Get public key of the RSA keypair.
	GetRSAPublicKey(ctx context.Context, in *GetRSAPublicKeyRequest, opts ...grpc.CallOption) (*GetRSAPublicKeyResponse, error)
	// Sign message stream with RSA. It will use SHA512_RSA_PKCS (RS512) algorithm to sign the message.
	RSASignStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_RSASignStreamClient, error)
	// Validate signature of the message stream using RSA key-pairs public key. It will use SHA512_RSA_PKCS (RS512) algorithm to verify the message.
	RSAVerifyStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_RSAVerifyStreamClient, error)
	// Sign message with RSA. The data must be short (max several kilobytes). If i is longer - use `RSASignStream` instead. It will use SHA512_RSA_PKCS (RS512) algorithm to sign the message.
	RSASign(ctx context.Context, in *RSASignRequest, opts ...grpc.CallOption) (*RSASignResponse, error)
	// Validate signature of the message using RSA key-pairs public key. The data must be short (max several kilobytes). If i is longer - use `RSAVerifyStream` instead. It will use SHA512_RSA_PKCS (RS512) algorithm to verify the message.
	RSAVerify(ctx context.Context, in *RSAVerifyRequest, opts ...grpc.CallOption) (*RSAVerifyResponse, error)
	// Calculate HMAC signature for input data stream. HMAC secret never leaves the HSM (hardware security module). It automatically uses the best available HMAC algorithm for currently used HSM.
	HMACSignStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_HMACSignStreamClient, error)
	// Verify HMAC signature of the data stream.
	HMACVerifyStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_HMACVerifyStreamClient, error)
	// Calculate HMAC signature for input data. The data must be short (max several kilobytes). If it is longer - use `HMACSignStream` instead. HMAC secret never leaves the HSM (hardware security module). It automatically uses the best available HMAC algorithm for currently used HSM.
	HMACSign(ctx context.Context, in *HMACSignRequest, opts ...grpc.CallOption) (*HMACSignResponse, error)
	// Verify HMAC signature of the data. The data must be short (max several kilobytes). If it is longer - use `HMACVerifyStream` instead.
	HMACVerify(ctx context.Context, in *HMACVerifyRequest, opts ...grpc.CallOption) (*HMACVerifyResponse, error)
	// Encrypt data stream. Encryption secret never leaves the HSM (hardware security module). It automatically uses the best available encryption algorithm for currently used HSM. It will only select the algorithm that is capable of decrypting from the middle of the whole data (partial decryption).
	EncryptStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_EncryptStreamClient, error)
	// Decrypt data stream. If you want to decrypt part of the information from the middle of the whole data - ensure, that the first chunk of the data, that you are sending is padded by 2048 bit.
	DecryptStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_DecryptStreamClient, error)
	// Encrypt data. The data must be short (max several kilobytes). If it is longer - use `EncryptStream` instead. Encryption secret never leaves the HSM (hardware security module). It automatically uses the best available encryption algorithm for currently used HSM. It will only select the algorithm that is capable of decrypting from the middle of the whole data (partial decryption).
	Encrypt(ctx context.Context, in *EncryptRequest, opts ...grpc.CallOption) (*EncryptResponse, error)
	// Decrypt data. The data must be short (max several kilobytes). If it is longer - use `DecryptStream` instead.
	Decrypt(ctx context.Context, in *DecryptRequest, opts ...grpc.CallOption) (*DecryptResponse, error)
}

type vaultServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVaultServiceClient(cc grpc.ClientConnInterface) VaultServiceClient {
	return &vaultServiceClient{cc}
}

func (c *vaultServiceClient) Seal(ctx context.Context, in *SealRequest, opts ...grpc.CallOption) (*SealResponse, error) {
	out := new(SealResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/Seal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) Unseal(ctx context.Context, in *UnsealRequest, opts ...grpc.CallOption) (*UnsealResponse, error) {
	out := new(UnsealResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/Unseal", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) UpdateSealSecret(ctx context.Context, in *UpdateSealSecretRequest, opts ...grpc.CallOption) (*UpdateSealSecretResponse, error) {
	out := new(UpdateSealSecretResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/UpdateSealSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*GetStatusResponse, error) {
	out := new(GetStatusResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/GetStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) EnsureRSAKeyPair(ctx context.Context, in *EnsureRSAKeyPairRequest, opts ...grpc.CallOption) (*EnsureRSAKeyPairResponse, error) {
	out := new(EnsureRSAKeyPairResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/EnsureRSAKeyPair", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) GetRSAPublicKey(ctx context.Context, in *GetRSAPublicKeyRequest, opts ...grpc.CallOption) (*GetRSAPublicKeyResponse, error) {
	out := new(GetRSAPublicKeyResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/GetRSAPublicKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) RSASignStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_RSASignStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &VaultService_ServiceDesc.Streams[0], "/system_vault.VaultService/RSASignStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &vaultServiceRSASignStreamClient{stream}
	return x, nil
}

type VaultService_RSASignStreamClient interface {
	Send(*RSASignStreamRequest) error
	CloseAndRecv() (*RSASignStreamResponse, error)
	grpc.ClientStream
}

type vaultServiceRSASignStreamClient struct {
	grpc.ClientStream
}

func (x *vaultServiceRSASignStreamClient) Send(m *RSASignStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *vaultServiceRSASignStreamClient) CloseAndRecv() (*RSASignStreamResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RSASignStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *vaultServiceClient) RSAVerifyStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_RSAVerifyStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &VaultService_ServiceDesc.Streams[1], "/system_vault.VaultService/RSAVerifyStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &vaultServiceRSAVerifyStreamClient{stream}
	return x, nil
}

type VaultService_RSAVerifyStreamClient interface {
	Send(*RSAVerifyStreamRequest) error
	CloseAndRecv() (*RSAVerifyStreamResponse, error)
	grpc.ClientStream
}

type vaultServiceRSAVerifyStreamClient struct {
	grpc.ClientStream
}

func (x *vaultServiceRSAVerifyStreamClient) Send(m *RSAVerifyStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *vaultServiceRSAVerifyStreamClient) CloseAndRecv() (*RSAVerifyStreamResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RSAVerifyStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *vaultServiceClient) RSASign(ctx context.Context, in *RSASignRequest, opts ...grpc.CallOption) (*RSASignResponse, error) {
	out := new(RSASignResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/RSASign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) RSAVerify(ctx context.Context, in *RSAVerifyRequest, opts ...grpc.CallOption) (*RSAVerifyResponse, error) {
	out := new(RSAVerifyResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/RSAVerify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) HMACSignStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_HMACSignStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &VaultService_ServiceDesc.Streams[2], "/system_vault.VaultService/HMACSignStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &vaultServiceHMACSignStreamClient{stream}
	return x, nil
}

type VaultService_HMACSignStreamClient interface {
	Send(*HMACSignStreamRequest) error
	CloseAndRecv() (*HMACSignStreamResponse, error)
	grpc.ClientStream
}

type vaultServiceHMACSignStreamClient struct {
	grpc.ClientStream
}

func (x *vaultServiceHMACSignStreamClient) Send(m *HMACSignStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *vaultServiceHMACSignStreamClient) CloseAndRecv() (*HMACSignStreamResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HMACSignStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *vaultServiceClient) HMACVerifyStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_HMACVerifyStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &VaultService_ServiceDesc.Streams[3], "/system_vault.VaultService/HMACVerifyStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &vaultServiceHMACVerifyStreamClient{stream}
	return x, nil
}

type VaultService_HMACVerifyStreamClient interface {
	Send(*HMACVerifyStreamRequest) error
	CloseAndRecv() (*HMACVerifyStreamResponse, error)
	grpc.ClientStream
}

type vaultServiceHMACVerifyStreamClient struct {
	grpc.ClientStream
}

func (x *vaultServiceHMACVerifyStreamClient) Send(m *HMACVerifyStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *vaultServiceHMACVerifyStreamClient) CloseAndRecv() (*HMACVerifyStreamResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HMACVerifyStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *vaultServiceClient) HMACSign(ctx context.Context, in *HMACSignRequest, opts ...grpc.CallOption) (*HMACSignResponse, error) {
	out := new(HMACSignResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/HMACSign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) HMACVerify(ctx context.Context, in *HMACVerifyRequest, opts ...grpc.CallOption) (*HMACVerifyResponse, error) {
	out := new(HMACVerifyResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/HMACVerify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) EncryptStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_EncryptStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &VaultService_ServiceDesc.Streams[4], "/system_vault.VaultService/EncryptStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &vaultServiceEncryptStreamClient{stream}
	return x, nil
}

type VaultService_EncryptStreamClient interface {
	Send(*EncryptStreamRequest) error
	Recv() (*EncryptStreamResponse, error)
	grpc.ClientStream
}

type vaultServiceEncryptStreamClient struct {
	grpc.ClientStream
}

func (x *vaultServiceEncryptStreamClient) Send(m *EncryptStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *vaultServiceEncryptStreamClient) Recv() (*EncryptStreamResponse, error) {
	m := new(EncryptStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *vaultServiceClient) DecryptStream(ctx context.Context, opts ...grpc.CallOption) (VaultService_DecryptStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &VaultService_ServiceDesc.Streams[5], "/system_vault.VaultService/DecryptStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &vaultServiceDecryptStreamClient{stream}
	return x, nil
}

type VaultService_DecryptStreamClient interface {
	Send(*DecryptStreamRequest) error
	Recv() (*DecryptStreamResponse, error)
	grpc.ClientStream
}

type vaultServiceDecryptStreamClient struct {
	grpc.ClientStream
}

func (x *vaultServiceDecryptStreamClient) Send(m *DecryptStreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *vaultServiceDecryptStreamClient) Recv() (*DecryptStreamResponse, error) {
	m := new(DecryptStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *vaultServiceClient) Encrypt(ctx context.Context, in *EncryptRequest, opts ...grpc.CallOption) (*EncryptResponse, error) {
	out := new(EncryptResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/Encrypt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vaultServiceClient) Decrypt(ctx context.Context, in *DecryptRequest, opts ...grpc.CallOption) (*DecryptResponse, error) {
	out := new(DecryptResponse)
	err := c.cc.Invoke(ctx, "/system_vault.VaultService/Decrypt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VaultServiceServer is the server API for VaultService service.
// All implementations must embed UnimplementedVaultServiceServer
// for forward compatibility
type VaultServiceServer interface {
	// Close and encrypt vault. After sealing, most of the operations will not be accessible.
	Seal(context.Context, *SealRequest) (*SealResponse, error)
	// Decrypt and open vault. Must be done before most of the operations with vault secrets.
	Unseal(context.Context, *UnsealRequest) (*UnsealResponse, error)
	// Set up new seal secret and reincrypt vault. The vault must be unsealed before this operation. You don't need to unseal vault after this operation.
	// This operation requires you to have administrator access to the HSM. Check PKCS11 spec. If you are using emulated HSM (by default) this will be the same as the seal/unseal secret by default ("12345678"). Change it.
	UpdateSealSecret(context.Context, *UpdateSealSecretRequest) (*UpdateSealSecretResponse, error)
	// Get current status of the vault.
	GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error)
	// Creates RSA key pair if it doesnt exist. Private key never leaves the HSM (hardware security module).
	EnsureRSAKeyPair(context.Context, *EnsureRSAKeyPairRequest) (*EnsureRSAKeyPairResponse, error)
	// Get public key of the RSA keypair.
	GetRSAPublicKey(context.Context, *GetRSAPublicKeyRequest) (*GetRSAPublicKeyResponse, error)
	// Sign message stream with RSA. It will use SHA512_RSA_PKCS (RS512) algorithm to sign the message.
	RSASignStream(VaultService_RSASignStreamServer) error
	// Validate signature of the message stream using RSA key-pairs public key. It will use SHA512_RSA_PKCS (RS512) algorithm to verify the message.
	RSAVerifyStream(VaultService_RSAVerifyStreamServer) error
	// Sign message with RSA. The data must be short (max several kilobytes). If i is longer - use `RSASignStream` instead. It will use SHA512_RSA_PKCS (RS512) algorithm to sign the message.
	RSASign(context.Context, *RSASignRequest) (*RSASignResponse, error)
	// Validate signature of the message using RSA key-pairs public key. The data must be short (max several kilobytes). If i is longer - use `RSAVerifyStream` instead. It will use SHA512_RSA_PKCS (RS512) algorithm to verify the message.
	RSAVerify(context.Context, *RSAVerifyRequest) (*RSAVerifyResponse, error)
	// Calculate HMAC signature for input data stream. HMAC secret never leaves the HSM (hardware security module). It automatically uses the best available HMAC algorithm for currently used HSM.
	HMACSignStream(VaultService_HMACSignStreamServer) error
	// Verify HMAC signature of the data stream.
	HMACVerifyStream(VaultService_HMACVerifyStreamServer) error
	// Calculate HMAC signature for input data. The data must be short (max several kilobytes). If it is longer - use `HMACSignStream` instead. HMAC secret never leaves the HSM (hardware security module). It automatically uses the best available HMAC algorithm for currently used HSM.
	HMACSign(context.Context, *HMACSignRequest) (*HMACSignResponse, error)
	// Verify HMAC signature of the data. The data must be short (max several kilobytes). If it is longer - use `HMACVerifyStream` instead.
	HMACVerify(context.Context, *HMACVerifyRequest) (*HMACVerifyResponse, error)
	// Encrypt data stream. Encryption secret never leaves the HSM (hardware security module). It automatically uses the best available encryption algorithm for currently used HSM. It will only select the algorithm that is capable of decrypting from the middle of the whole data (partial decryption).
	EncryptStream(VaultService_EncryptStreamServer) error
	// Decrypt data stream. If you want to decrypt part of the information from the middle of the whole data - ensure, that the first chunk of the data, that you are sending is padded by 2048 bit.
	DecryptStream(VaultService_DecryptStreamServer) error
	// Encrypt data. The data must be short (max several kilobytes). If it is longer - use `EncryptStream` instead. Encryption secret never leaves the HSM (hardware security module). It automatically uses the best available encryption algorithm for currently used HSM. It will only select the algorithm that is capable of decrypting from the middle of the whole data (partial decryption).
	Encrypt(context.Context, *EncryptRequest) (*EncryptResponse, error)
	// Decrypt data. The data must be short (max several kilobytes). If it is longer - use `DecryptStream` instead.
	Decrypt(context.Context, *DecryptRequest) (*DecryptResponse, error)
	mustEmbedUnimplementedVaultServiceServer()
}

// UnimplementedVaultServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVaultServiceServer struct {
}

func (UnimplementedVaultServiceServer) Seal(context.Context, *SealRequest) (*SealResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Seal not implemented")
}
func (UnimplementedVaultServiceServer) Unseal(context.Context, *UnsealRequest) (*UnsealResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unseal not implemented")
}
func (UnimplementedVaultServiceServer) UpdateSealSecret(context.Context, *UpdateSealSecretRequest) (*UpdateSealSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSealSecret not implemented")
}
func (UnimplementedVaultServiceServer) GetStatus(context.Context, *GetStatusRequest) (*GetStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedVaultServiceServer) EnsureRSAKeyPair(context.Context, *EnsureRSAKeyPairRequest) (*EnsureRSAKeyPairResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnsureRSAKeyPair not implemented")
}
func (UnimplementedVaultServiceServer) GetRSAPublicKey(context.Context, *GetRSAPublicKeyRequest) (*GetRSAPublicKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRSAPublicKey not implemented")
}
func (UnimplementedVaultServiceServer) RSASignStream(VaultService_RSASignStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method RSASignStream not implemented")
}
func (UnimplementedVaultServiceServer) RSAVerifyStream(VaultService_RSAVerifyStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method RSAVerifyStream not implemented")
}
func (UnimplementedVaultServiceServer) RSASign(context.Context, *RSASignRequest) (*RSASignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RSASign not implemented")
}
func (UnimplementedVaultServiceServer) RSAVerify(context.Context, *RSAVerifyRequest) (*RSAVerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RSAVerify not implemented")
}
func (UnimplementedVaultServiceServer) HMACSignStream(VaultService_HMACSignStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method HMACSignStream not implemented")
}
func (UnimplementedVaultServiceServer) HMACVerifyStream(VaultService_HMACVerifyStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method HMACVerifyStream not implemented")
}
func (UnimplementedVaultServiceServer) HMACSign(context.Context, *HMACSignRequest) (*HMACSignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HMACSign not implemented")
}
func (UnimplementedVaultServiceServer) HMACVerify(context.Context, *HMACVerifyRequest) (*HMACVerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HMACVerify not implemented")
}
func (UnimplementedVaultServiceServer) EncryptStream(VaultService_EncryptStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method EncryptStream not implemented")
}
func (UnimplementedVaultServiceServer) DecryptStream(VaultService_DecryptStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method DecryptStream not implemented")
}
func (UnimplementedVaultServiceServer) Encrypt(context.Context, *EncryptRequest) (*EncryptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Encrypt not implemented")
}
func (UnimplementedVaultServiceServer) Decrypt(context.Context, *DecryptRequest) (*DecryptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Decrypt not implemented")
}
func (UnimplementedVaultServiceServer) mustEmbedUnimplementedVaultServiceServer() {}

// UnsafeVaultServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VaultServiceServer will
// result in compilation errors.
type UnsafeVaultServiceServer interface {
	mustEmbedUnimplementedVaultServiceServer()
}

func RegisterVaultServiceServer(s grpc.ServiceRegistrar, srv VaultServiceServer) {
	s.RegisterService(&VaultService_ServiceDesc, srv)
}

func _VaultService_Seal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SealRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).Seal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/Seal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).Seal(ctx, req.(*SealRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_Unseal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsealRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).Unseal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/Unseal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).Unseal(ctx, req.(*UnsealRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_UpdateSealSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSealSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).UpdateSealSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/UpdateSealSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).UpdateSealSecret(ctx, req.(*UpdateSealSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).GetStatus(ctx, req.(*GetStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_EnsureRSAKeyPair_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnsureRSAKeyPairRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).EnsureRSAKeyPair(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/EnsureRSAKeyPair",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).EnsureRSAKeyPair(ctx, req.(*EnsureRSAKeyPairRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_GetRSAPublicKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRSAPublicKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).GetRSAPublicKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/GetRSAPublicKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).GetRSAPublicKey(ctx, req.(*GetRSAPublicKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_RSASignStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VaultServiceServer).RSASignStream(&vaultServiceRSASignStreamServer{stream})
}

type VaultService_RSASignStreamServer interface {
	SendAndClose(*RSASignStreamResponse) error
	Recv() (*RSASignStreamRequest, error)
	grpc.ServerStream
}

type vaultServiceRSASignStreamServer struct {
	grpc.ServerStream
}

func (x *vaultServiceRSASignStreamServer) SendAndClose(m *RSASignStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *vaultServiceRSASignStreamServer) Recv() (*RSASignStreamRequest, error) {
	m := new(RSASignStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _VaultService_RSAVerifyStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VaultServiceServer).RSAVerifyStream(&vaultServiceRSAVerifyStreamServer{stream})
}

type VaultService_RSAVerifyStreamServer interface {
	SendAndClose(*RSAVerifyStreamResponse) error
	Recv() (*RSAVerifyStreamRequest, error)
	grpc.ServerStream
}

type vaultServiceRSAVerifyStreamServer struct {
	grpc.ServerStream
}

func (x *vaultServiceRSAVerifyStreamServer) SendAndClose(m *RSAVerifyStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *vaultServiceRSAVerifyStreamServer) Recv() (*RSAVerifyStreamRequest, error) {
	m := new(RSAVerifyStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _VaultService_RSASign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RSASignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).RSASign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/RSASign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).RSASign(ctx, req.(*RSASignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_RSAVerify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RSAVerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).RSAVerify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/RSAVerify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).RSAVerify(ctx, req.(*RSAVerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_HMACSignStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VaultServiceServer).HMACSignStream(&vaultServiceHMACSignStreamServer{stream})
}

type VaultService_HMACSignStreamServer interface {
	SendAndClose(*HMACSignStreamResponse) error
	Recv() (*HMACSignStreamRequest, error)
	grpc.ServerStream
}

type vaultServiceHMACSignStreamServer struct {
	grpc.ServerStream
}

func (x *vaultServiceHMACSignStreamServer) SendAndClose(m *HMACSignStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *vaultServiceHMACSignStreamServer) Recv() (*HMACSignStreamRequest, error) {
	m := new(HMACSignStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _VaultService_HMACVerifyStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VaultServiceServer).HMACVerifyStream(&vaultServiceHMACVerifyStreamServer{stream})
}

type VaultService_HMACVerifyStreamServer interface {
	SendAndClose(*HMACVerifyStreamResponse) error
	Recv() (*HMACVerifyStreamRequest, error)
	grpc.ServerStream
}

type vaultServiceHMACVerifyStreamServer struct {
	grpc.ServerStream
}

func (x *vaultServiceHMACVerifyStreamServer) SendAndClose(m *HMACVerifyStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *vaultServiceHMACVerifyStreamServer) Recv() (*HMACVerifyStreamRequest, error) {
	m := new(HMACVerifyStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _VaultService_HMACSign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HMACSignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).HMACSign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/HMACSign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).HMACSign(ctx, req.(*HMACSignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_HMACVerify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HMACVerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).HMACVerify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/HMACVerify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).HMACVerify(ctx, req.(*HMACVerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_EncryptStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VaultServiceServer).EncryptStream(&vaultServiceEncryptStreamServer{stream})
}

type VaultService_EncryptStreamServer interface {
	Send(*EncryptStreamResponse) error
	Recv() (*EncryptStreamRequest, error)
	grpc.ServerStream
}

type vaultServiceEncryptStreamServer struct {
	grpc.ServerStream
}

func (x *vaultServiceEncryptStreamServer) Send(m *EncryptStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *vaultServiceEncryptStreamServer) Recv() (*EncryptStreamRequest, error) {
	m := new(EncryptStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _VaultService_DecryptStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(VaultServiceServer).DecryptStream(&vaultServiceDecryptStreamServer{stream})
}

type VaultService_DecryptStreamServer interface {
	Send(*DecryptStreamResponse) error
	Recv() (*DecryptStreamRequest, error)
	grpc.ServerStream
}

type vaultServiceDecryptStreamServer struct {
	grpc.ServerStream
}

func (x *vaultServiceDecryptStreamServer) Send(m *DecryptStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *vaultServiceDecryptStreamServer) Recv() (*DecryptStreamRequest, error) {
	m := new(DecryptStreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _VaultService_Encrypt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EncryptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).Encrypt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/Encrypt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).Encrypt(ctx, req.(*EncryptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VaultService_Decrypt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecryptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VaultServiceServer).Decrypt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/system_vault.VaultService/Decrypt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VaultServiceServer).Decrypt(ctx, req.(*DecryptRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VaultService_ServiceDesc is the grpc.ServiceDesc for VaultService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VaultService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "system_vault.VaultService",
	HandlerType: (*VaultServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Seal",
			Handler:    _VaultService_Seal_Handler,
		},
		{
			MethodName: "Unseal",
			Handler:    _VaultService_Unseal_Handler,
		},
		{
			MethodName: "UpdateSealSecret",
			Handler:    _VaultService_UpdateSealSecret_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _VaultService_GetStatus_Handler,
		},
		{
			MethodName: "EnsureRSAKeyPair",
			Handler:    _VaultService_EnsureRSAKeyPair_Handler,
		},
		{
			MethodName: "GetRSAPublicKey",
			Handler:    _VaultService_GetRSAPublicKey_Handler,
		},
		{
			MethodName: "RSASign",
			Handler:    _VaultService_RSASign_Handler,
		},
		{
			MethodName: "RSAVerify",
			Handler:    _VaultService_RSAVerify_Handler,
		},
		{
			MethodName: "HMACSign",
			Handler:    _VaultService_HMACSign_Handler,
		},
		{
			MethodName: "HMACVerify",
			Handler:    _VaultService_HMACVerify_Handler,
		},
		{
			MethodName: "Encrypt",
			Handler:    _VaultService_Encrypt_Handler,
		},
		{
			MethodName: "Decrypt",
			Handler:    _VaultService_Decrypt_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "RSASignStream",
			Handler:       _VaultService_RSASignStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "RSAVerifyStream",
			Handler:       _VaultService_RSAVerifyStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "HMACSignStream",
			Handler:       _VaultService_HMACSignStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "HMACVerifyStream",
			Handler:       _VaultService_HMACVerifyStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "EncryptStream",
			Handler:       _VaultService_EncryptStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "DecryptStream",
			Handler:       _VaultService_DecryptStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "vault.proto",
}
