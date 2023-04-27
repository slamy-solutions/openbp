package pkcs

import (
	"context"
	"errors"
	"io"
	"os"
	"strconv"
)

var ErrPKCSNotLoggedIn = errors.New("pkcs is not logged in")
var ErrPKCSBadLoginPassword = errors.New("bad pkcs password")

var ErrRSAKeyDoesntExist = errors.New("RSA key doesnt exist")
var ErrHMACKeyDoesntExist = errors.New("HMAC key doesnt exist")
var ErrEncryptionKeyDoesntExist = errors.New("Encryption key doesnt exist")

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
	// Signs message using RSA private key
	SignRSA(ctx context.Context, name string, message *io.PipeReader) ([]byte, error)
	// Verifies message using RSA public key
	VerifyRSA(ctx context.Context, name string, message *io.PipeReader, signature []byte) (bool, error)

	// Create HMAC based on message
	SignHMAC(ctx context.Context, message *io.PipeReader) ([]byte, error)
	// Checks if HMAC corresponds to provided message
	VerifyHMAC(ctx context.Context, message *io.PipeReader, signature []byte) (bool, error)

	// Encrypt message
	Encrypt(ctx context.Context, plain *io.PipeReader, encrypted *io.PipeWriter) error
	// Decrypt message
	Decrypt(ctx context.Context, encrypted *io.PipeReader, plain *io.PipeWriter) error

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
