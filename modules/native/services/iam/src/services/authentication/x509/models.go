package x509

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"time"

	x509GRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/x509"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type X509InMongo struct {
	UUID        primitive.ObjectID `bson:"_id,omitempty"`
	Identity    primitive.ObjectID `bson:"identity"`
	Disabled    bool               `bson:"disabled"`
	Description string             `bson:"description"`
	PublicKey   []byte             `bson:"publicKey"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func (c *X509InMongo) ToGRPCCertificate(namespace string) *x509GRPC.Certificate {
	return &x509GRPC.Certificate{
		Namespace:   namespace,
		Uuid:        c.UUID.Hex(),
		Identity:    c.Identity.Hex(),
		Disabled:    c.Disabled,
		Description: c.Description,
		PublicKey:   c.PublicKey,
		Created:     timestamppb.New(c.Created),
		Updated:     timestamppb.New(c.Updated),
		Version:     c.Version,
	}
}

func (c *X509InMongo) ToSignedX509(ctx context.Context, namespace string, signer *x509Signer) ([]byte, error) {
	publicKey, err := x509.ParsePKCS1PublicKey(c.PublicKey)
	if err != nil {
		return nil, errors.New("invalid public key format: DER format expected")
	}

	serial := new(big.Int)
	serial.SetString(c.UUID.Hex(), 16)
	cert := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   c.Identity.Hex(),
			Organization: []string{namespace},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	return x509.CreateCertificate(rand.Reader, cert, signer.ca, publicKey, signer.GetContextedSigner(ctx))
}
