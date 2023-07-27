package authTools

import (
	"errors"
	"net/textproto"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
)

var authHeaderName = textproto.CanonicalMIMEHeaderKey("authorization")

type CheckAuthData struct {
	AccessGranted bool
	StatusCode    int
	ErrorMessage  string

	Namespace    string
	TokenUUID    string
	IdentityUUID string
}

func CheckAuth(ctx *gin.Context, nativeStub *native.NativeStub, scopes []*auth.Scope) (*CheckAuthData, error) {
	authHeader := ctx.Request.Header.Get(authHeaderName)
	if authHeader == "" {
		return &CheckAuthData{AccessGranted: false, StatusCode: 401, ErrorMessage: "Auth header not found"}, nil
	}

	splitedHeader := strings.Split(authHeader, " ")
	if len(splitedHeader) != 2 {
		return &CheckAuthData{AccessGranted: false, StatusCode: 401, ErrorMessage: "Auth header has invalid format"}, nil
	}

	authResponse, err := nativeStub.Services.IAM.Auth.CheckAccessWithToken(ctx.Request.Context(), &auth.CheckAccessWithTokenRequest{
		AccessToken: splitedHeader[1],
		Scopes:      scopes,
	})
	if err != nil {
		return &CheckAuthData{AccessGranted: false, StatusCode: 500, ErrorMessage: "Internal server error"}, errors.New("failed to check access with token: " + err.Error())
	}

	accessGranted := authResponse.Status == auth.CheckAccessWithTokenResponse_OK
	statusCode := 200
	switch authResponse.Status {
	case auth.CheckAccessWithTokenResponse_OK:
		break
	case auth.CheckAccessWithTokenResponse_UNAUTHORIZED:
		statusCode = 403
	default:
		statusCode = 401
	}

	return &CheckAuthData{
		AccessGranted: accessGranted,
		StatusCode:    statusCode,
		ErrorMessage:  authResponse.Message,
		Namespace:     authResponse.GetNamespace(),
		TokenUUID:     authResponse.GetTokenUUID(),
		IdentityUUID:  authResponse.GetIdentityUUID(),
	}, nil
}

func FillLoggerWithAuthMetadata(logger *logrus.Entry, authData *CheckAuthData) *logrus.Entry {
	logger = logger.WithFields(logrus.Fields{
		"auth.namespace":    authData.Namespace,
		"auth.tokenUUID":    authData.TokenUUID,
		"auth.identityUUID": authData.IdentityUUID,
	})

	if authData.StatusCode != 200 {
		logger = logger.WithFields(logrus.Fields{
			"auth.statusCode": authData.StatusCode,
			"auth.message":    authData.ErrorMessage,
		})
	}

	return logger
}
