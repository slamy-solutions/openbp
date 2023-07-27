package x509

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const x509_root_ca_key_name = "openbp_native_iam_authentication_x509_rootCA"
const x509_root_ca_redis_key = "openbp_native_iam_authentication_x509_rootCA"

type x509Signer struct {
	systemStub *system.SystemStub

	publicKey *rsa.PublicKey
	ca        *x509.Certificate
	mut       *sync.RWMutex

	//Public() PublicKey
	//Sign(rand io.Reader, digest []byte, opts SignerOpts) (signature []byte, err error)
}
type x509ContextedSigner struct {
	x509Signer

	ctx context.Context
}

var ErrX509SignerVaultSealed = errors.New("the vault is sealed")

func newX509Signer(systemStub *system.SystemStub) *x509Signer {
	return &x509Signer{
		systemStub: systemStub,
		publicKey:  nil,
		ca:         nil,
		mut:        &sync.RWMutex{},
	}
}

func (s *x509Signer) EnsureReady(ctx context.Context) error {
	s.mut.Lock()
	defer s.mut.Unlock()

	if s.publicKey != nil {
		return nil
	}

	// Load public key
	getRSAKeyResponse, err := s.systemStub.Vault.GetRSAPublicKey(ctx, &vault.GetRSAPublicKeyRequest{
		KeyName: x509_root_ca_key_name,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return ErrX509SignerVaultSealed
			}
			if st.Code() == codes.NotFound {
				_, err = s.systemStub.Vault.EnsureRSAKeyPair(ctx, &vault.EnsureRSAKeyPairRequest{
					KeyName: x509_root_ca_key_name,
				})
				if err != nil {
					return errors.New("public key not found in system_vault. Unexpected error while creating it: " + err.Error())
				}

				getRSAKeyResponse, err := s.systemStub.Vault.GetRSAPublicKey(ctx, &vault.GetRSAPublicKeyRequest{
					KeyName: x509_root_ca_key_name,
				})
				if err != nil {
					return errors.New("public key not found in system_vault. Key-pair was created, but still failed to get it: " + err.Error())
				}
				key, err := x509.ParsePKCS1PublicKey(getRSAKeyResponse.PublicKey)
				if err != nil {
					return errors.New("failed to parse RSA public key: " + err.Error())
				}
				s.publicKey = key
			}
		} else {
			return errors.New("unexpected error while trying to get public key from the system_vault: " + err.Error())
		}
	} else {
		key, err := x509.ParsePKCS1PublicKey(getRSAKeyResponse.PublicKey)
		if err != nil {
			return errors.New("failed to parse RSA public key: " + err.Error())
		}
		s.publicKey = key
	}

	// Load CA certificate
	caBytes, err := s.systemStub.Redis.Get(ctx, x509_root_ca_redis_key).Bytes()
	if err != nil {
		if err != redis.Nil {
			s.publicKey = nil
			return errors.New("error while searching in redis for root CA certificate: " + err.Error())
		}

		// Try create certifica and add it to the Redis
		ca := &x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Organization:  []string{"OpenBP"},
				Country:       []string{},
				Province:      []string{},
				Locality:      []string{},
				StreetAddress: []string{},
				PostalCode:    []string{},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(30, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
		}

		caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, s.publicKey, s.GetContextedSigner(ctx))
		if err != nil {
			s.publicKey = nil
			return errors.New("failed to sign certificate: " + err.Error())
		}

		// Someone may already set it to the redis. Lets check if there are no race conditions
		rsp := s.systemStub.Redis.SetNX(ctx, x509_root_ca_redis_key, caBytes, 0)
		seted, err := rsp.Result()
		if err != nil {
			s.publicKey = nil
			return errors.New("failed to update certificate in redis: " + err.Error())
		}

		if seted {
			newCA, _ := x509.ParseCertificate(caBytes)
			s.ca = newCA
		} else {
			// Someone change CA certificate in meantime. We have to load
			caBytes, err := s.systemStub.Redis.Get(ctx, x509_root_ca_redis_key).Bytes()
			if err != nil {
				s.publicKey = nil
				return errors.New("failed to get certificate from redis: " + err.Error())
			}

			loadedCA, err := x509.ParseCertificate(caBytes)
			if err != nil {
				s.publicKey = nil
				return errors.New("error while parsing CA cerfiticate from redis: " + err.Error())
			}
			s.ca = loadedCA
		}
	} else {
		loadedCA, err := x509.ParseCertificate(caBytes)
		if err != nil {
			s.publicKey = nil
			return errors.New("error while parsing CA cerfiticate from redis: " + err.Error())
		}
		s.ca = loadedCA
	}

	return nil
}

func (s *x509Signer) GetContextedSigner(ctx context.Context) *x509ContextedSigner {
	return &x509ContextedSigner{
		x509Signer: x509Signer{
			publicKey:  s.publicKey,
			systemStub: s.systemStub,
			ca:         s.ca,
			mut:        s.mut,
		},
		ctx: ctx,
	}
}

func (s *x509ContextedSigner) Public() crypto.PublicKey {
	return s.publicKey
}

// DER1 prefix for SHA256
var rsaSignerPrefix = []byte{0x30, 0x31, 0x30, 0x0d, 0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01, 0x05, 0x00, 0x04, 0x20}

func (s *x509ContextedSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {

	// We will leave padding for vault -----
	//hashLen := opts.HashFunc().Size()
	//if hashLen != len(digest) {
	//	return nil, errors.New("dugest lenght is not the same as expected from hash function")
	//}

	//tLen := len(rsaSignerPrefix) + hashLen
	//k := s.publicKey.Size()
	//if k < tLen+11 {
	//	return nil, errors.New("message too long")
	//}

	// EM = 0x00 || 0x01 || PS || 0x00 || T
	// em := make([]byte, k)
	// em[1] = 1
	// for i := 2; i < k-tLen-1; i++ {
	//	em[i] = 0xff
	//}
	//copy(em[k-tLen:k-hashLen], rsaSignerPrefix)
	//copy(em[k-hashLen:k], digest)
	// -----------

	em := make([]byte, len(rsaSignerPrefix)+len(digest))
	copy(em[:len(rsaSignerPrefix)], rsaSignerPrefix)
	copy(em[len(rsaSignerPrefix):], digest)

	signResult, err := s.systemStub.Vault.RSASign(s.ctx, &vault.RSASignRequest{
		KeyName:   x509_root_ca_key_name,
		Data:      em,
		Mechanism: vault.RSASignMechanism_RSA_PKCS,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return nil, ErrX509SignerVaultSealed
			}
		}

		return nil, err
	}

	return signResult.Signature, nil
}
