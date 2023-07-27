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

func rsaSignMechanismFromGRPC(in vault.RSASignMechanism) uint {
	switch in {
	case vault.RSASignMechanism_DEFAULT:
		return pkcs.RSASignAlgoDefault
	case vault.RSASignMechanism_SHA512_RSA:
		return pkcs.RSASignAlgoSHA512
	case vault.RSASignMechanism_RSA_PKCS:
		return pkcs.RSASignAlgoRSAPKCS
	case vault.RSASignMechanism_SHA256_RSA:
		return pkcs.RSASignAlgoSHA256
	}

	return pkcs.RSASignAlgoDefault
}

func (s *VaultService) RSASignStream(srv vault.VaultService_RSASignStreamServer) error {
	ctx := srv.Context()

	dataReader, dataWriter := io.Pipe()
	defer dataWriter.Close()
	defer dataReader.Close()

	keyPairNameChanel := make(chan struct {
		name      string
		mechanism vault.RSASignMechanism
		err       error
	})
	defer close(keyPairNameChanel)

	go func() {
		nameSended := false
		for {
			message, err := srv.Recv()
			if err != nil {
				if !nameSended {
					keyPairNameChanel <- struct {
						name      string
						mechanism vault.RSASignMechanism
						err       error
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
					name      string
					mechanism vault.RSASignMechanism
					err       error
				}{name: message.KeyName, err: nil, mechanism: message.Mechanism}
				nameSended = true
			}
			_, err = dataWriter.Write(message.Data)
			if err != nil {
				return
			}
		}
	}()

	metadata := <-keyPairNameChanel
	if metadata.err != nil {
		log.Error("[GRPC Vault Service]-(RSASignStream) Internal error while reading RSA keypair name: " + metadata.err.Error())
		return status.Error(codes.Internal, "error while reading RSA keypair name: "+metadata.err.Error())
	}

	signature, err := s.pkcsHandle.SignRSAStream(ctx, metadata.name, dataReader, rsaSignMechanismFromGRPC(metadata.mechanism))

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		if err == pkcs.ErrRSAKeyDoesntExist {
			return status.Error(codes.NotFound, "RSA key-pair doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(RSASignStream) Internal error while signing message with PKCS11 RSA: " + err.Error())
		return status.Error(codes.Internal, "error while signing message with PKCS11 RSA: "+err.Error())
	}

	srv.SendAndClose(&vault.RSASignStreamResponse{
		Signature: signature,
	})
	return nil
}
func (s *VaultService) RSAVerifyStream(srv vault.VaultService_RSAVerifyStreamServer) error {
	ctx := srv.Context()

	dataReader, dataWriter := io.Pipe()
	defer dataWriter.Close()
	defer dataReader.Close()

	metadataChanel := make(chan struct {
		name      string
		signature []byte
		mechanism vault.RSASignMechanism
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
						mechanism vault.RSASignMechanism
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
					mechanism vault.RSASignMechanism
					err       error
				}{name: message.KeyName, signature: message.Signature, err: nil, mechanism: message.Mechanism}
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
		log.Error("[GRPC Vault Service]-(RSAVerifyStream) Internal error while reading RSA keypair name: " + metadata.err.Error())
		return status.Error(codes.Internal, "error while reading RSA keypair name: "+metadata.err.Error())
	}

	valid, err := s.pkcsHandle.VerifyRSAStream(ctx, metadata.name, dataReader, metadata.signature, rsaSignMechanismFromGRPC(metadata.mechanism))

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		if err == pkcs.ErrRSAKeyDoesntExist {
			return status.Error(codes.NotFound, "RSA key-pair doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(RSAVerifyStream) Internal error while verifiying message with PKCS11 RSA: " + err.Error())
		return status.Error(codes.Internal, "error while verifiying message with PKCS11 RSA: "+err.Error())
	}

	srv.SendAndClose(&vault.RSAVerifyStreamResponse{
		Valid: valid,
	})
	return nil
}

func (s *VaultService) RSASign(ctx context.Context, in *vault.RSASignRequest) (*vault.RSASignResponse, error) {
	signature, err := s.pkcsHandle.SignRSA(ctx, in.KeyName, in.Data, rsaSignMechanismFromGRPC(in.Mechanism))

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return nil, status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		//TODO: This is for the future
		if err == pkcs.ErrRSAKeyDoesntExist {
			return nil, status.Error(codes.NotFound, "RSA key doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(RSASign) Internal error while signing message with PKCS11 RSA: " + err.Error())
		return nil, status.Error(codes.Internal, "error while signing message with PKCS11 RSA: "+err.Error())
	}

	return &vault.RSASignResponse{
		Signature: signature,
	}, status.Error(codes.OK, "")
}
func (s *VaultService) RSAVerify(ctx context.Context, in *vault.RSAVerifyRequest) (*vault.RSAVerifyResponse, error) {
	valid, err := s.pkcsHandle.VerifyRSA(ctx, in.KeyName, in.Data, in.Signature, rsaSignMechanismFromGRPC(in.Mechanism))

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return nil, status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		if err == pkcs.ErrRSAKeyDoesntExist {
			return nil, status.Error(codes.NotFound, "RSA key doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(RSAVerify) Internal error while verifiying message with PKCS11 RSA: " + err.Error())
		return nil, status.Error(codes.Internal, "error while verifiying message with PKCS11 RSA: "+err.Error())
	}

	return &vault.RSAVerifyResponse{
		Valid: valid,
	}, status.Error(codes.OK, "")
}

func (s *VaultService) HMACSignStream(srv vault.VaultService_HMACSignStreamServer) error {
	ctx := srv.Context()

	dataReader, dataWriter := io.Pipe()
	defer dataWriter.Close()
	defer dataReader.Close()

	go func() {
		for {
			message, err := srv.Recv()
			if err != nil {
				if err == io.EOF {
					dataWriter.Close()
					return
				}
				dataWriter.CloseWithError(errors.New("error while receiving message from grpc: " + err.Error()))
				return
			}
			_, err = dataWriter.Write(message.Data)
			if err != nil {
				return
			}
		}
	}()

	signature, err := s.pkcsHandle.SignHMACStream(ctx, dataReader)

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		//TODO: This is for the future
		if err == pkcs.ErrHMACKeyDoesntExist {
			return status.Error(codes.NotFound, "HMAC key doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(HMACSignStream) Internal error while signing message with PKCS11 HMAC: " + err.Error())
		return status.Error(codes.Internal, "error while signing message with PKCS11 HMAC: "+err.Error())
	}

	srv.SendAndClose(&vault.HMACSignStreamResponse{
		Signature: signature,
	})
	return nil
}

func (s *VaultService) HMACVerifyStream(srv vault.VaultService_HMACVerifyStreamServer) error {
	ctx := srv.Context()

	dataReader, dataWriter := io.Pipe()
	defer dataWriter.Close()
	defer dataReader.Close()

	metadataChanel := make(chan struct {
		signature []byte
		err       error
	})
	defer close(metadataChanel)

	go func() {
		metaSended := false
		for {
			message, err := srv.Recv()
			if err != nil {
				if !metaSended {
					metadataChanel <- struct {
						signature []byte
						err       error
					}{signature: []byte{}, err: errors.New("failed to receive first chunk of data from grpc stream: " + err.Error())}
				}
				if err == io.EOF {
					dataWriter.Close()
					return
				}
				dataWriter.CloseWithError(errors.New("error while receiving message from grpc: " + err.Error()))
				return
			}
			if !metaSended {
				metadataChanel <- struct {
					signature []byte
					err       error
				}{signature: message.Signature, err: nil}
				metaSended = true
			}
			_, err = dataWriter.Write(message.Data)
			if err != nil {
				return
			}
		}
	}()

	metadata := <-metadataChanel
	if metadata.err != nil {
		log.Error("[GRPC Vault Service]-(HMACVerifyStream) Internal error while reading signature: " + metadata.err.Error())
		return status.Error(codes.Internal, "error while reading HMAC signature: "+metadata.err.Error())
	}

	valid, err := s.pkcsHandle.VerifyHMACStream(ctx, dataReader, metadata.signature)

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		if err == pkcs.ErrHMACKeyDoesntExist {
			return status.Error(codes.NotFound, "HMAC key doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(HMACVerifyStream) Internal error while verifiying message with PKCS11 HMAC: " + err.Error())
		return status.Error(codes.Internal, "error while verifiying message with PKCS11 HMAC: "+err.Error())
	}

	srv.SendAndClose(&vault.HMACVerifyStreamResponse{
		Valid: valid,
	})
	return nil
}

func (s *VaultService) HMACSign(ctx context.Context, in *vault.HMACSignRequest) (*vault.HMACSignResponse, error) {
	signature, err := s.pkcsHandle.SignHMAC(ctx, in.Data)

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return nil, status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		//TODO: This is for the future
		if err == pkcs.ErrHMACKeyDoesntExist {
			return nil, status.Error(codes.NotFound, "HMAC key doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(HMACSign) Internal error while signing message with PKCS11 HMAC: " + err.Error())
		return nil, status.Error(codes.Internal, "error while signing message with PKCS11 HMAC: "+err.Error())
	}

	return &vault.HMACSignResponse{
		Signature: signature,
	}, status.Error(codes.OK, "")
}
func (s *VaultService) HMACVerify(ctx context.Context, in *vault.HMACVerifyRequest) (*vault.HMACVerifyResponse, error) {
	valid, err := s.pkcsHandle.VerifyHMAC(ctx, in.Data, in.Signature)

	if err != nil {
		if err == pkcs.ErrPKCSNotLoggedIn {
			return nil, status.Error(codes.FailedPrecondition, "the vault is sealed")
		}

		if err == pkcs.ErrHMACKeyDoesntExist {
			return nil, status.Error(codes.NotFound, "HMAC key doesnt exist")
		}

		log.Error("[GRPC Vault Service]-(HMACVerify) Internal error while verifiying message with PKCS11 HMAC: " + err.Error())
		return nil, status.Error(codes.Internal, "error while verifiying message with PKCS11 HMAC: "+err.Error())
	}

	return &vault.HMACVerifyResponse{
		Valid: valid,
	}, status.Error(codes.OK, "")
}
