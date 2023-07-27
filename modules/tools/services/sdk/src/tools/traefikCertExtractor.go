package tools

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"google.golang.org/grpc/metadata"
)

type traefikCertificateExctractor struct {
}

func NewTraefikCertificateExtractor() CertificateExtractor {
	return &traefikCertificateExctractor{}
}

func (e *traefikCertificateExctractor) Exctract(ctx context.Context) (*x509.Certificate, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("failed to extract metadata from the context")
	}

	headerValues := md.Get("X-Forwarded-Tls-Client-Cert")
	if len(headerValues) == 0 {
		return nil, errors.New("failed to extract X-Forwarded-Tls-Client-Cert header from context. Header not found")
	}
	pemCertString := headerValues[0]

	pemBlock, _ := pem.Decode([]byte(pemCertString))
	if pemBlock == nil {
		return nil, errors.New("failed to parse PEM block from certificate header")
	}
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse certificate from PEM block: " + err.Error())
	}

	return cert, nil
}

/*
func (e *traefikCertificateExctractor) Start() error {
	return nil
}
func (e *traefikCertificateExctractor) Close() error {
	return nil
}
*/
