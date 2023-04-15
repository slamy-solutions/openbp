package service

import (
	"context"
	"sync"

	"github.com/miekg/pkcs11"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VaultService struct {
	vault.UnimplementedVaultServiceServer

	pkcsCtx         *pkcs11.Ctx
	masterKeyMutex  sync.Mutex
	masterKey       []byte
	masterKeyLoaded bool
}

func NewVaultGRPCService(pkcsCtx *pkcs11.Ctx, masterKey []byte, masterKeyLoaded bool) *VaultService {
	return &VaultService{
		pkcsCtx:         pkcsCtx,
		masterKey:       masterKey,
		masterKeyLoaded: masterKeyLoaded,
		masterKeyMutex:  sync.Mutex{},
	}
}

func (s *VaultService) setMasterKey(key []byte) {
	s.masterKeyMutex.Lock()
	defer s.masterKeyMutex.Unlock()

	s.masterKey = key
	s.masterKeyLoaded = false
}

func (s *VaultService) resetMasterKey() {
	s.masterKeyMutex.Lock()
	defer s.masterKeyMutex.Unlock()

	s.masterKey = []byte{}
	s.masterKeyLoaded = false
}

func (s *VaultService) Seal(ctx context.Context, in *vault.SealRequest) (*vault.SealResponse, error) {
	s.resetMasterKey()

	return nil, status.Errorf(codes.Unimplemented, "method Seal not implemented")
}
func (s *VaultService) Unseal(ctx context.Context, in *vault.UnsealRequest) (*vault.UnsealResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unseal not implemented")
}
func (s *VaultService) UpdateSealSecret(ctx context.Context, in *vault.UpdateSealSecretRequest) (*vault.UpdateSealSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSealSecret not implemented")
}
func (s *VaultService) GetStatus(ctx context.Context, in *vault.GetStatusRequest) (*vault.GetStatusResponse, error) {
	return &vault.GetStatusResponse{
		Sealed: !s.masterKeyLoaded,
	}, status.Error(codes.OK, "")
}
func (s *VaultService) EnsureRSAKeyPair(ctx context.Context, in *vault.EnsureRSAKeyPairRequest) (*vault.EnsureRSAKeyPairResponse, error) {
	freeSlots, err := s.pkcsCtx.GetSlotList(false)
	s.pkcsCtx.GetTokenInfo()

	return nil, status.Errorf(codes.Unimplemented, "method EnsureRSAKeyPair not implemented")
}
func (s *VaultService) GetRSAPublicKey(ctx context.Context, in *vault.GetRSAPublicKeyRequest) (*vault.GetRSAPublicKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRSAPublicKey not implemented")
}
func (s *VaultService) RSASign(ctx context.Context, in *vault.RSASignRequest) (*vault.RSASignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RSASign not implemented")
}
func (s *VaultService) RSAValidatePublic(ctx context.Context, in *vault.RSAValidatePublicRequest) (*vault.RSAValidatePublicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RSAValidatePublic not implemented")
}
