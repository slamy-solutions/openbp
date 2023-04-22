package service

import (
	"context"
	"errors"
	"io"

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
	return &vault.GetStatusResponse{
		Sealed: !s.pkcsHandle.IsLoggedIn(),
	}, status.Error(codes.OK, "")
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
	key, err := s.pkcsHandle.GetRSAPublicKey(ctx, in.KeyName)
	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return nil, status.Error(codes.FailedPrecondition, "the vault is sealed")
		}
		if err == pkcs.ErrRSAKeyDoesntExist {
			return nil, status.Error(codes.NotFound, "RSA key-pair doesnt exist")
		}
		log.Error("[GRPC Vault Service]-(GetRSAPublicKey) Internal error while getting public key of the RSA key-pair using PKCS11: " + err.Error())
		return nil, status.Error(codes.Internal, "error while getting public key of the RSA key-pair using PKCS11: "+err.Error())
	}

	return &vault.GetRSAPublicKeyResponse{
		PublicKey: key,
	}, nil
}

func (s *VaultService) RSASign(srv vault.VaultService_RSASignServer) error {
	ctx := srv.Context()

	dataReader, dataWriter := io.Pipe()
	defer dataWriter.Close()
	defer dataReader.Close()

	keyPairNameChanel := make(chan struct {
		name string
		err  error
	})
	defer close(keyPairNameChanel)

	go func() {
		nameSended := false
		for {
			message, err := srv.Recv()
			if err != nil {
				if !nameSended {
					keyPairNameChanel <- struct {
						name string
						err  error
					}{name: "", err: errors.New("failed to receive first chunk of data from grpc stream: " + err.Error())}
				}
				if err == io.EOF {
					dataWriter.Close()
					return
				}
				dataWriter.CloseWithError(errors.New("error while receiving message from grpc: " + err.Error()))
				return
			}
			if !nameSended {
				keyPairNameChanel <- struct {
					name string
					err  error
				}{name: message.KeyName, err: nil}
				nameSended = true
			}
			_, err = dataWriter.Write(message.Data)
			if err != nil {
				return
			}
		}
	}()

	nameData := <-keyPairNameChanel
	if nameData.err != nil {
		log.Error("[GRPC Vault Service]-(RSASign) Internal error while reading RSA keypair name: " + nameData.err.Error())
		return status.Error(codes.Internal, "error while reading RSA keypair name: "+nameData.err.Error())
	}

	signature, err := s.pkcsHandle.SignRSA(ctx, nameData.name, dataReader)

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		if err == pkcs.ErrRSAKeyDoesntExist {
			return status.Error(codes.NotFound, "RSA key-pair doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(RSASign) Internal error while signing message with PKCS11 RSA: " + err.Error())
		return status.Error(codes.Internal, "error while signing message with PKCS11 RSA: "+err.Error())
	}

	srv.SendAndClose(&vault.RSASignResponse{
		Signature: signature,
	})
	return nil
}
func (s *VaultService) RSAVerify(srv vault.VaultService_RSAVerifyServer) error {
	ctx := srv.Context()

	dataReader, dataWriter := io.Pipe()
	defer dataWriter.Close()
	defer dataReader.Close()

	metadataChanel := make(chan struct {
		name      string
		signature []byte
		err       error
	})
	defer close(metadataChanel)

	go func() {
		nameSended := false
		for {
			message, err := srv.Recv()
			if err != nil {
				if !nameSended {
					metadataChanel <- struct {
						name      string
						signature []byte
						err       error
					}{name: "", signature: []byte{}, err: errors.New("failed to receive first chunk of data from grpc stream: " + err.Error())}
				}
				if err == io.EOF {
					dataWriter.Close()
					return
				}
				dataWriter.CloseWithError(errors.New("error while receiving message from grpc: " + err.Error()))
				return
			}
			if !nameSended {
				metadataChanel <- struct {
					name      string
					signature []byte
					err       error
				}{name: message.KeyName, signature: message.Signature, err: nil}
				nameSended = true
			}
			_, err = dataWriter.Write(message.Data)
			if err != nil {
				return
			}
		}
	}()

	metadata := <-metadataChanel
	if metadata.err != nil {
		log.Error("[GRPC Vault Service]-(RSAVerify) Internal error while reading RSA keypair name: " + metadata.err.Error())
		return status.Error(codes.Internal, "error while reading RSA keypair name: "+metadata.err.Error())
	}

	valid, err := s.pkcsHandle.VerifyRSA(ctx, metadata.name, dataReader, metadata.signature)

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		if err == pkcs.ErrRSAKeyDoesntExist {
			return status.Error(codes.NotFound, "RSA key-pair doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(RSAVerify) Internal error while verifiying message with PKCS11 RSA: " + err.Error())
		return status.Error(codes.Internal, "error while verifiying message with PKCS11 RSA: "+err.Error())
	}

	srv.SendAndClose(&vault.RSAVerifyResponse{
		Valid: valid,
	})
	return nil
}
