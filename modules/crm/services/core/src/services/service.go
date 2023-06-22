package services

import (
	"context"

	"github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceServer struct {
	service.UnimplementedServiceServiceServer
}

func NewServiceServer() *ServiceServer {
	return &ServiceServer{}
}

func (*ServiceServer) Ping(ctx context.Context, in *service.PingRequest) (*service.PingResponse, error) {
	return &service.PingResponse{}, status.Error(codes.OK, "")
}
