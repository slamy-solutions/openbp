package service

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"github.com/slamy-solutions/openbp/modules/system/services/vault/src/pkcs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VaultService struct {
	vault.UnimplementedVaultServiceServer

	pkcsHandle pkcs.PKCS
	sealer     *pkcs.Sealer
}

func NewVaultGRPCService(pkcsHandle pkcs.PKCS, sealer *pkcs.Sealer) *VaultService {
	return &VaultService{
		pkcsHandle: pkcsHandle,
		sealer:     sealer,
	}
}

func (s *VaultService) Seal(ctx context.Context, in *vault.SealRequest) (*vault.SealResponse, error) {
	s.sealer.Seal()
	return nil, status.Error(codes.OK, "")
}
func (s *VaultService) Unseal(ctx context.Context, in *vault.UnsealRequest) (*vault.UnsealResponse, error) {
	unsealed, err := s.sealer.Unseal(in.Secret)
	if err != nil {
		log.Error("[GRPC Vault Service]-(Unseal) Internal error while unsealing: " + err.Error())
		return nil, status.Error(codes.Internal, "error while unsealing: "+err.Error())
	}

	return &vault.UnsealResponse{Success: unsealed}, status.Error(codes.OK, "")
}
func (s *VaultService) UpdateSealSecret(ctx context.Context, in *vault.UpdateSealSecretRequest) (*vault.UpdateSealSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSealSecret not implemented")
}
func (s *VaultService) GetStatus(ctx context.Context, in *vault.GetStatusRequest) (*vault.GetStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSealSecret not implemented")
}
func (s *VaultService) EnsureRSAKeyPair(ctx context.Context, in *vault.EnsureRSAKeyPairRequest) (*vault.EnsureRSAKeyPairResponse, error) {
	err := s.pkcsHandle.EnsureRSAKeyPair(ctx, in.KeyName)
	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return nil, status.Error(codes.FailedPrecondition, "the vault is sealed")
		}
		log.Error("[GRPC Vault Service]-(EnsureRSAKeyPair) Internal error while ensuring RSA keys via PKCS: " + err.Error())
		return nil, status.Error(codes.Internal, "error while ensuring RSA keys via PKCS: "+err.Error())
	}

	return &vault.EnsureRSAKeyPairResponse{}, status.Error(codes.OK, "")
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
