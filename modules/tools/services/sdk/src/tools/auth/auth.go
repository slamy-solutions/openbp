package auth

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/oauth"
	tools "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools"
)

type Scope struct {
	Namespace string
	Resources []string
	Actions   []string
}

func extractAuthHeader(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values, ok := md["authorization"]
	if !ok {
		return ""
	}
	return values[0]
}

func getAuthToken(ctx context.Context) (string, error) {
	header := extractAuthHeader(ctx)
	if header == "" {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated. No authorization header provided.")
	}

	splits := strings.SplitN(header, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "Bad authorization header. Expected \"Bearer <token>\".")
	}

	if !strings.EqualFold(splits[0], "bearer") {
		return "", status.Errorf(codes.Unauthenticated, "Request unauthenticated. Bearer token expected.")
	}
	return splits[1], nil
}

func AuthorizeRPC(ctx context.Context, modulesStub *tools.ModulesStub, requiredScope Scope) error {
	token, err := getAuthToken(ctx)
	if err != nil {
		return err
	}

	response, err := modulesStub.Native.Services.IamOAuth.CheckAccess(ctx, &oauth.CheckAccessRequest{
		AccessToken: token,
		Scopes: []*oauth.Scope{
			{
				Namespace: requiredScope.Namespace,
				Resources: requiredScope.Resources,
				Actions:   requiredScope.Actions,
			},
		},
	})

	if err != nil {
		log.Error("gRPC error while checking access token: " + err.Error())
		return status.Error(codes.Internal, "")
	}

	if response.Status == oauth.CheckAccessResponse_OK {
		return nil
	} else {
		return status.Error(codes.Unauthenticated, response.Message)
	}
}
