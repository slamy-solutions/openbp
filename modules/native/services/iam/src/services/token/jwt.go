package token

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
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
	rsaKey       *rsa.PublicKey
	keyLoadMutex *sync.RWMutex

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
		rsaKey:       nil,
		keyLoadMutex: &sync.RWMutex{},
		systemStub:   systemStub,
	}
}

func (s *jwtService) loadPublicKey(ctx context.Context) error {
	s.keyLoadMutex.Lock()
	defer s.keyLoadMutex.Unlock()

	if s.rsaKey != nil {
		return nil
	}

	getRSAKeyResponse, err := s.systemStub.Vault.GetRSAPublicKey(ctx, &vault.GetRSAPublicKeyRequest{
		KeyName: jwtVaultKeyName,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return ErrVaultSealed
			} else if st.Code() == codes.NotFound {
				_, err = s.systemStub.Vault.EnsureRSAKeyPair(ctx, &vault.EnsureRSAKeyPairRequest{
					KeyName: jwtVaultKeyName,
				})
				if err != nil {
					return errors.New("JWT RSA public key not found in system_vault. Unexpected error while creating it: " + err.Error())
				}

				getRSAKeyResponse, err := s.systemStub.Vault.GetRSAPublicKey(ctx, &vault.GetRSAPublicKeyRequest{
					KeyName: jwtVaultKeyName,
				})
				if err != nil {
					return errors.New("JWT RSA public key not found in system_vault. Key-pair was created, but still failed to get it: " + err.Error())
				}
				key, err := x509.ParsePKCS1PublicKey(getRSAKeyResponse.PublicKey)
				if err != nil {
					return errors.New("failed to parse RSA public key: " + err.Error())
				}
				s.rsaKey = key
			} else {
				return errors.New("unexpected error while trying to get JWT RSA public key from the system_vault: " + st.Message())
			}
		} else {
			return errors.New("unexpected error while trying to get JWT RSA public key from the system_vault: " + err.Error())
		}
	} else {
		key, err := x509.ParsePKCS1PublicKey(getRSAKeyResponse.PublicKey)
		if err != nil {
			return errors.New("failed to parse RSA public key: " + err.Error())
		}
		s.rsaKey = key
	}

	// Remove key from the memory after small period of time. This will force even public key to be deleted when vault is sealed.
	go func() {
		time.Sleep(time.Second * 30)

		s.keyLoadMutex.Lock()
		s.rsaKey = nil
		s.keyLoadMutex.Unlock()
	}()

	return nil
}

func (s *jwtService) getPublicKey(ctx context.Context) (*rsa.PublicKey, error) {
	s.keyLoadMutex.RLock()
	loaded := s.rsaKey != nil
	s.keyLoadMutex.RUnlock()

	if !loaded {
		err := s.loadPublicKey(ctx)
		if err != nil {
			return nil, err
		}
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
		if _, ok := token.Method.(*goJWT.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}

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
	// Make sure Key Pair created. In this case, public key must exist and be loaded.
	s.getPublicKey(ctx)

	token := goJWT.NewWithClaims(goJWT.SigningMethodRS512, data)
	stringToSign, err := token.SigningString()
	if err != nil {
		return "", errors.New("failed to generate signing string: " + err.Error())
	}
	bytesToSign := []byte(stringToSign)

	signResponse, err := s.systemStub.Vault.RSASign(ctx, &vault.RSASignRequest{
		KeyName: jwtVaultKeyName,
		Data:    bytesToSign,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return "", ErrVaultSealed
			}
		}
		return "", errors.New("unexpected error while signing JWT token in the system_vault: " + err.Error())
	}

	signature := base64.RawURLEncoding.EncodeToString(signResponse.Signature)
	return strings.Join([]string{stringToSign, signature}, "."), nil
}
