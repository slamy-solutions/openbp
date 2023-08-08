package pkcs

import (
	"context"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/miekg/pkcs11"
)

var ErrPKCSNotLoggedIn = errors.New("pkcs is not logged in")
var ErrPKCSBadLoginPassword = errors.New("bad pkcs password")

var ErrRSAKeyDoesntExist = errors.New("RSA key doesnt exist")
var ErrHMACKeyDoesntExist = errors.New("HMAC key doesnt exist")
var ErrEncryptionKeyDoesntExist = errors.New("encryption key doesnt exist")

const (
	RSASignAlgoDefault = pkcs11.CKM_SHA512_RSA_PKCS
	RSASignAlgoSHA512  = pkcs11.CKM_SHA512_RSA_PKCS
	RSASignAlgoSHA256  = pkcs11.CKM_SHA256_RSA_PKCS
	RSASignAlgoRSAPKCS = pkcs11.CKM_RSA_PKCS
)

type PKCS interface {
	Initialize() error

	GetProviderName() string

	IsLoggedIn() bool
	EnsureSessionAndLogIn(pin string) error
	LogOutAndCloseSession() error
	UpdatePins(ctx context.Context, adminPin string, newAdminPin string, newPin string) error

	// Creates RSA key-pair with specified name if it doesnt exist
	EnsureRSAKeyPair(ctx context.Context, name string) error
	// Returns RSA public key if it exists
	GetRSAPublicKey(ctx context.Context, name string) ([]byte, error)
	// Signs message stream using RSA private key
	SignRSAStream(ctx context.Context, name string, message *io.PipeReader, mechanism uint) ([]byte, error)
	// Verifies message stream using RSA public key
	VerifyRSAStream(ctx context.Context, name string, message *io.PipeReader, signature []byte, mechanism uint) (bool, error)
	// Signs message using RSA private key
	SignRSA(ctx context.Context, name string, message []byte, mechanism uint) ([]byte, error)
	// Verifies message using RSA public key
	VerifyRSA(ctx context.Context, name string, message []byte, signature []byte, mechanism uint) (bool, error)

	// Create HMAC based on message stream
	SignHMACStream(ctx context.Context, message *io.PipeReader) ([]byte, error)
	// Checks if HMAC corresponds to provided message stream
	VerifyHMACStream(ctx context.Context, message *io.PipeReader, signature []byte) (bool, error)
	// Create HMAC based on message
	SignHMAC(ctx context.Context, message []byte) ([]byte, error)
	// Checks if HMAC corresponds to provided message
	VerifyHMAC(ctx context.Context, message []byte, signature []byte) (bool, error)

	// Encrypt message stream
	EncryptStream(ctx context.Context, plain *io.PipeReader, encrypted *io.PipeWriter) error
	// Decrypt message stream
	DecryptStream(ctx context.Context, encrypted *io.PipeReader, plain *io.PipeWriter) error
	// Encrypt message
	Encrypt(ctx context.Context, plain []byte) ([]byte, error)
	// Decrypt message
	Decrypt(ctx context.Context, encrypted []byte) ([]byte, error)

	Close() error
}

// Create new PKCS instance using configuration from environment variables
func NewPKCSFromEnv() (PKCS, error) {
	hsmProvider := os.Getenv("HSM_PROVIDER")

	var selectedPkcs PKCS

	switch hsmProvider {
	case "softhsm2":
		library := os.Getenv("SOFTHSM2_PKCS11_LIBRARY_PATH")
		tokenLabel := os.Getenv("SOFTHSM2_PKCS11_TOKEN_LABEL")
		p, err := NewSoftHSM2PKCSHandle(library, tokenLabel)
		if err != nil {
			return nil, errors.New("failed to setup SOFTHSM2 library: " + err.Error())
		}
		selectedPkcs = p
	case "dynamic":
		fallthrough
	default:
		library := os.Getenv("DYNAMIC_PKCS11_LIBRARY_PATH")
		slot, err := strconv.ParseUint(os.Getenv("DYNAMIC_PKCS11_SLOT"), 10, 32)
		if err != nil {
			return nil, errors.New("failed to parse DYNAMIC_PKCS11_SLOT environment variable: " + err.Error())
		}
		selectedPkcs = NewDynamicPKCSHandle(library, uint(slot))
	}

	selectedPkcs = WrapPKCSInOTel(selectedPkcs)

	return selectedPkcs, nil
}
