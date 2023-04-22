package pkcs

import (
	"context"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type otelPKCSWraper struct {
	innerPKCS PKCS
	tracer    trace.Tracer
}

func WrapPKCSInOTel(pkcs PKCS) PKCS {
	wrappedPkcs := &otelPKCSWraper{
		innerPKCS: pkcs,
		tracer:    otel.GetTracerProvider().Tracer("github.com/slamy-solutions/openbp/modules/system/services/vault/src/pkcs/otel"),
	}

	return wrappedPkcs
}

func (wraper *otelPKCSWraper) Initialize() error {
	return wraper.innerPKCS.Initialize()
}

func (wraper *otelPKCSWraper) GetProviderName() string {
	return wraper.innerPKCS.GetProviderName()
}

func (wraper *otelPKCSWraper) IsLoggedIn() bool {
	return wraper.innerPKCS.IsLoggedIn()
}

func (wraper *otelPKCSWraper) EnsureSessionAndLogIn(password string) error {
	return wraper.innerPKCS.EnsureSessionAndLogIn(password)
}
func (wraper *otelPKCSWraper) LogOutAndCloseSession() error {
	return wraper.innerPKCS.LogOutAndCloseSession()
}

func (wraper *otelPKCSWraper) UpdatePins(ctx context.Context, adminPin string, newAdminPin string, newPin string) error {
	ctx, span := wraper.tracer.Start(ctx, "pkcs.UpdatePins")
	defer span.End()

	return wraper.innerPKCS.UpdatePins(ctx, adminPin, newAdminPin, newPin)
}

func (wraper *otelPKCSWraper) EnsureRSAKeyPair(ctx context.Context, name string) error {
	ctx, span := wraper.tracer.Start(ctx, "pkcs.EnsureRSAKeyPair")
	defer span.End()

	return wraper.innerPKCS.EnsureRSAKeyPair(ctx, name)
}

func (wraper *otelPKCSWraper) GetRSAPublicKey(ctx context.Context, name string) ([]byte, error) {
	ctx, span := wraper.tracer.Start(ctx, "pkcs.GetRSAPublicKey")
	defer span.End()

	return wraper.innerPKCS.GetRSAPublicKey(ctx, name)
}

func (wraper *otelPKCSWraper) SignRSA(ctx context.Context, name string, message *io.PipeReader) ([]byte, error) {
	ctx, span := wraper.tracer.Start(ctx, "pkcs.SignRSA")
	defer span.End()

	return wraper.innerPKCS.SignRSA(ctx, name, message)
}

func (wraper *otelPKCSWraper) VerifyRSA(ctx context.Context, name string, message *io.PipeReader, signature []byte) (bool, error) {
	ctx, span := wraper.tracer.Start(ctx, "pkcs.VerifyRSA")
	defer span.End()

	return wraper.innerPKCS.VerifyRSA(ctx, name, message, signature)
}

func (wraper *otelPKCSWraper) Close() error {
	return wraper.innerPKCS.Close()
}
