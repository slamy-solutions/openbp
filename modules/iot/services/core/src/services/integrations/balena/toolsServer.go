package balena

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/integrations/balena/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ToolsServer struct {
	balena.UnimplementedBalenaToolsServiceServer

	apiClient api.Client
	logger    *logrus.Entry
}

func NewToolsServer(logger *logrus.Entry, apiClient api.Client) *ToolsServer {
	return &ToolsServer{
		apiClient: apiClient,
		logger:    logger,
	}
}

func (s *ToolsServer) VerifyConnectionData(ctx context.Context, in *balena.VerifyConnectionDataRequest) (*balena.VerifyConnectionDataResponse, error) {
	pingError := s.apiClient.Ping(ctx, api.BalenaServerInfo{
		BaseURL:  in.Url,
		APIToken: in.Token,
	})
	if pingError == nil {
		return &balena.VerifyConnectionDataResponse{
			Status:  balena.VerifyConnectionDataResponse_OK,
			Message: "",
		}, status.Error(codes.OK, "")
	}

	if errors.Is(pingError, api.ErrServerURLInvalid) || errors.Is(pingError, api.ErrFailedToCreateRequest) {
		return &balena.VerifyConnectionDataResponse{
			Status:  balena.VerifyConnectionDataResponse_BAD_URL,
			Message: pingError.Error(),
		}, status.Error(codes.OK, "")
	}

	if errors.Is(pingError, api.ErrServerConnectionError) {
		return &balena.VerifyConnectionDataResponse{
			Status:  balena.VerifyConnectionDataResponse_SERVER_UNAVAILABLE,
			Message: pingError.Error(),
		}, status.Error(codes.OK, "")
	}

	if errors.Is(pingError, api.ErrInvalidStatusCode) {
		return &balena.VerifyConnectionDataResponse{
			Status:  balena.VerifyConnectionDataResponse_SERVER_BAD_RESPONSE,
			Message: pingError.Error(),
		}, status.Error(codes.OK, "")
	}

	return nil, status.Error(codes.Internal, "unhandled ping error: "+pingError.Error())
}
