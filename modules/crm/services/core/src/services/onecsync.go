package services

import (
	"context"
	"errors"
	"log/slog"

	onecSync "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/onecsync"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/sync"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OneCSyncService struct {
	onecSync.UnimplementedOneCSyncServiceServer

	onecSyncEngine sync.SyncEngine

	logger *slog.Logger
}

func NewOneCSyncServer(onecSyncEngine sync.SyncEngine, logger *slog.Logger) *OneCSyncService {
	return &OneCSyncService{
		onecSyncEngine: onecSyncEngine,
		logger:         logger,
	}
}

func (s *OneCSyncService) Sync(ctx context.Context, in *onecSync.SyncRequest) (*onecSync.SyncResponse, error) {
	err := s.onecSyncEngine.SyncNow(ctx, in.Namespace)
	ok := err == nil
	message := ""
	if err != nil {
		message = err.Error()
	}

	return &onecSync.SyncResponse{
		Ok:           ok,
		ErrorMessage: message,
	}, status.Error(codes.OK, "")
}
func (s *OneCSyncService) GetLog(ctx context.Context, in *onecSync.GetLogRequest) (*onecSync.GetLogResponse, error) {
	events, totalCount, err := s.onecSyncEngine.ListSyncEvents(in.Namespace, int(in.Skip), int(in.Limit))
	if err != nil {
		err := errors.Join(errors.New("failed to get log"), err)
		s.logger.With(slog.String("route", "GetLog")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get log: %s", err.Error())
	}

	var responseEvents []*onecSync.OneCSyncEvent = make([]*onecSync.OneCSyncEvent, len(events))
	for i, event := range events {
		responseEvents[i] = event.ToGRPC()
	}

	return &onecSync.GetLogResponse{
		Events:     responseEvents,
		TotalCount: int32(totalCount),
	}, status.Error(codes.OK, "")
}
