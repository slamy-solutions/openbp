package balena

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SyncServer struct {
	balena.UnimplementedBalenaSyncServiceServer

	logger *logrus.Entry
}

func NewSyncServer(logger *logrus.Entry) *SyncServer {
	return &SyncServer{
		logger: logger,
	}
}

func (s *SyncServer) SyncNow(ctx context.Context, in *balena.SyncNowRequest) (*balena.SyncNowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncNow not implemented")
}
func (s *SyncServer) ListLog(in *balena.ListSyncLogRequest, out balena.BalenaSyncService_ListLogServer) error {
	return status.Errorf(codes.Unimplemented, "method ListLog not implemented")
}
func (s *SyncServer) CountLog(ctx context.Context, in *balena.CountSyncLogRequest) (*balena.CountSyncLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountLog not implemented")
}
func (s *SyncServer) GetLastSyncLog(ctx context.Context, in *balena.GetLastSyncLogRequest) (*balena.GetLastSyncLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastSyncLog not implemented")
}
