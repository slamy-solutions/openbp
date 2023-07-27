package auth

import (
	"crypto/rsa"
	cryptoX509 "crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/x509"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CertificateRouter struct {
	nativeStub *native.NativeStub
}

func NewCertificateRouter(nativeStub *native.NativeStub) *CertificateRouter {
	return &CertificateRouter{
		nativeStub: nativeStub,
	}
}

type formatedCertificate struct {
	Namespace   string    `json:"certificate"`
	UUID        string    `json:"uuid"`
	Identity    string    `json:"identity"`
	Disabled    bool      `json:"disabled"`
	Description string    `json:"description"`
	PublicKey   string    `json:"publicKey"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Version     uint64    `json:"version"`
}

func NewFormatedCertificateFromGRPC(cert *x509.Certificate, namespace string) *formatedCertificate {
	return &formatedCertificate{
		Namespace:   namespace,
		UUID:        cert.Uuid,
		Identity:    cert.Identity,
		Disabled:    cert.Disabled,
		Description: cert.Description,
		PublicKey:   base64.StdEncoding.EncodeToString(cert.PublicKey),
		Created:     cert.Created.AsTime(),
		Updated:     cert.Updated.AsTime(),
		Version:     cert.Version,
	}
}

type listCertificateForIdentityRequest struct {
	Namespace    string `form:"namespace" binding:"lte=32"`
	IdentityUUID string `form:"identityUUID" binding:"lte=64,required"`
	Skip         uint32 `form:"skip" binding:"gte=0"`
	Limit        uint32 `form:"limit" binding:"gte=0,lte=100"`
}
type listCertificateForIdentityResponse struct {
	Certificates []formatedCertificate `json:"certificates"`
	TotalCount   uint64                `json:"totalCount"`
}

func (r *CertificateRouter) ListCertificatesForIdentity(ctx *gin.Context) {
	var requestData listCertificateForIdentityRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.identity." + requestData.IdentityUUID},
			Actions:              []string{"native.iam.auth.certificate.list"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	certsStream, err := r.nativeStub.Services.IAM.Authentication.X509.ListForIdentity(ctx.Request.Context(), &x509.ListForIdentityRequest{
		Namespace: requestData.Namespace,
		Identity:  requestData.IdentityUUID,
		Skip:      uint64(requestData.Skip),
		Limit:     uint64(requestData.Limit),
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	certificates := make([]formatedCertificate, 0, 10)
	for {
		responseChunk, err := certsStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		certificates = append(certificates, *NewFormatedCertificateFromGRPC(responseChunk.Certificate, requestData.Namespace))
	}

	countResponse, err := r.nativeStub.Services.IAM.Authentication.X509.CountForIdentity(ctx.Request.Context(), &x509.CountForIdentityRequest{
		Namespace: requestData.Namespace,
		Identity:  requestData.IdentityUUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &listCertificateForIdentityResponse{Certificates: certificates, TotalCount: countResponse.Count})
}

type registerKeyAndGenerateCertificateRequest struct {
	Namespace    string `json:"namespace" binding:"lte=32"`
	IdentityUUID string `json:"identityUUID" binding:"lte=64,required"`
	PublicKey    string `json:"publicKey" binding:"gte=0,lte=10000"`
	Description  string `json:"description" binding:"gte=0,lte=128"`
}
type registerKeyAndGenerateCertificateResponse struct {
	Certificate *formatedCertificate `json:"certificate"`
	Raw         string               `json:"raw"`
}

func (r *CertificateRouter) RegisterKeyAndGenerateCertificate(ctx *gin.Context) {
	var requestData registerKeyAndGenerateCertificateRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.identity." + requestData.IdentityUUID},
			Actions:              []string{"native.iam.auth.certificate.register"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	block, _ := pem.Decode([]byte(requestData.PublicKey))
	if block == nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to parse public key as PEM"})
		return
	}
	key, err := cryptoX509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to parse key as RSA public key in X.509/SPKI format"})
		return
	}
	pubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "The key was parsed successfully, but this is not a RSA public key in X.509/SPKI format"})
		return
	}

	registerResponse, err := r.nativeStub.Services.IAM.Authentication.X509.RegisterAndGenerate(ctx.Request.Context(), &x509.RegisterAndGenerateRequest{
		Namespace:   requestData.Namespace,
		Identity:    requestData.IdentityUUID,
		PublicKey:   cryptoX509.MarshalPKCS1PublicKey(pubKey),
		Description: requestData.Description,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "Key file has invalid format. Maybe unsupported algorithm?"})
				return
			}
			if st.Code() == codes.FailedPrecondition {
				ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"message": "The vault is sealed"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cert, err := cryptoX509.ParseCertificate(registerResponse.Raw)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("error while parsing certificate from the IAM service: "+err.Error()))
		return
	}

	pemCertificate := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	if pemCertificate == nil {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("failed to encode certificate to PEM format"))
		return
	}

	ctx.JSON(http.StatusOK, &registerKeyAndGenerateCertificateResponse{Certificate: NewFormatedCertificateFromGRPC(registerResponse.Info, requestData.Namespace), Raw: string(pemCertificate)})
}

type disableCertificateRequest struct {
	Namespace       string `json:"namespace" binding:"lte=32"`
	CertificateUUID string `json:"certificateUUID" binding:"lte=64,required"`
}

func (r *CertificateRouter) Disable(ctx *gin.Context) {
	var requestData disableCertificateRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.auth.certificate." + requestData.CertificateUUID},
			Actions:              []string{"native.iam.auth.certificate.disable"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	_, err = r.nativeStub.Services.IAM.Authentication.X509.Disable(ctx.Request.Context(), &x509.DisableRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.CertificateUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Certificate not found"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type deleteCertificateRequest struct {
	Namespace       string `form:"namespace" binding:"lte=32"`
	CertificateUUID string `form:"certificateUUID" binding:"lte=64,required"`
}

func (r *CertificateRouter) Delete(ctx *gin.Context) {
	var requestData deleteCertificateRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.auth.certificate." + requestData.CertificateUUID},
			Actions:              []string{"native.iam.auth.certificate.delete"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	_, err = r.nativeStub.Services.IAM.Authentication.X509.Delete(ctx.Request.Context(), &x509.DeleteRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.CertificateUUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
