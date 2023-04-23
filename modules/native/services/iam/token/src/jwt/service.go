package jwt

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"sync"
	"time"

	goJWT "github.com/golang-jwt/jwt"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const jwtVaultKeyName = "native_iam_token_jwt"

type jwtService struct {
	rsaKey       []byte
	keyLoadMutex sync.RWMutex

	systemStub *system.SystemStub
}

type JWTService interface {
	JWTDataFromString(ctx context.Context, input string) (*JWTData, error)
	JWTDataToSignedString(ctx context.Context, data *JWTData) (string, error)
}

var ErrInvalidToken = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")
var ErrVaultSealed = errors.New("vault is sealed. Its impossible to perform any operations")

func NewJWTService(systemStub *system.SystemStub) JWTService {
	return &jwtService{
		rsaKey:       []byte{},
		keyLoadMutex: sync.RWMutex{},
		systemStub:   systemStub,
	}
}

func (s *jwtService) loadPublicKey(ctx context.Context) error {
	s.keyLoadMutex.Lock()
	defer s.keyLoadMutex.Unlock()

	if len(s.rsaKey) != 0 {
		return nil
	}

	getRSAKeyResponse, err := s.systemStub.Vault.GetRSAPublicKey(ctx, &vault.GetRSAPublicKeyRequest{
		KeyName: jwtVaultKeyName,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return ErrVaultSealed
			}
			if st.Code() == codes.NotFound {
				_, err = s.systemStub.Vault.EnsureRSAKeyPair(ctx, &vault.EnsureRSAKeyPairRequest{
					KeyName: jwtVaultKeyName,
				})
				if err != nil {
					return errors.New("JWT RSA public key not found system_vault. Unexpected error while creating it: " + err.Error())
				}

				getRSAKeyResponse, err := s.systemStub.Vault.GetRSAPublicKey(ctx, &vault.GetRSAPublicKeyRequest{
					KeyName: jwtVaultKeyName,
				})
				if err != nil {
					return errors.New("JWT RSA public key not found system_vault. Key-pair was created, but still failed to get it: " + err.Error())
				}
				s.rsaKey = getRSAKeyResponse.PublicKey
			}
		}

		return errors.New("unexpected error while trying to get JWT RSA public key from the system_vault: " + err.Error())
	} else {
		s.rsaKey = getRSAKeyResponse.PublicKey
	}

	// Remove key from the memory after small period of time. This will force even public key to be deleted when vault is sealed.
	go func() {
		time.Sleep(time.Second * 30)

		s.keyLoadMutex.Lock()
		defer s.keyLoadMutex.Unlock()

		s.rsaKey = []byte{}
	}()

	return nil
}

func (s *jwtService) getPublicKey(ctx context.Context) ([]byte, error) {
	{
		s.keyLoadMutex.RLock()
		defer s.keyLoadMutex.RUnlock()

		if len(s.rsaKey) != 0 {
			return s.rsaKey, nil
		}
	}

	err := s.loadPublicKey(ctx)
	if err != nil {
		return []byte{}, err
	}
	return s.rsaKey, nil
}

func (s *jwtService) JWTDataFromString(ctx context.Context, input string) (*JWTData, error) {
	rsaKey, err := s.getPublicKey(ctx)
	if err != nil {
		return nil, err
	}

	data := &JWTData{}
	_, err = goJWT.ParseWithClaims(input, data, func(token *goJWT.Token) (interface{}, error) {
		return rsaKey, nil
	})
	if err != nil {
		if err, ok := err.(*goJWT.ValidationError); ok {
			if err.Errors&goJWT.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrInvalidToken
	}

	return data, nil
}

func (s *jwtService) JWTDataToSignedString(ctx context.Context, data *JWTData) (string, error) {
	token := goJWT.NewWithClaims(goJWT.SigningMethodRS512, data)
	stringToSign, err := token.SigningString()
	if err != nil {
		return "", errors.New("failed to generate signing string: " + err.Error())
	}
	bytesToSign := []byte(stringToSign)

	signStream, err := s.systemStub.Vault.RSASign(ctx)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return "", ErrVaultSealed
			}
		}
		return "", errors.New("failed to open system_vault signing stream for JWT token data: " + err.Error())
	}
	for i := 0; i < len(bytesToSign); i += 4096 {
		end := i + 4096
		if end > len(bytesToSign) {
			end = len(bytesToSign)
		}
		err := signStream.Send(&vault.RSASignRequest{
			KeyName: jwtVaultKeyName,
			Data:    bytesToSign[i:end],
		})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.FailedPrecondition {
					return "", ErrVaultSealed
				}
			}

			return "", errors.New("failed to send chunk of JWT token to sign to system_vault: " + err.Error())
		}
	}
	signResponse, err := signStream.CloseAndRecv()
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return "", ErrVaultSealed
			}
		}
		return "", errors.New("unexpected error while signing JWT token in the syste_vault: " + err.Error())
	}

	signature := base64.StdEncoding.EncodeToString(signResponse.Signature)
	return strings.Join([]string{stringToSign, signature}, "."), nil
}
