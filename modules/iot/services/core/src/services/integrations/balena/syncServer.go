package balena

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SyncServer struct {
	balena.UnimplementedBalenaSyncServiceServer

	systemStub  *system.SystemStub
	syncManager SyncManager
	logger      *logrus.Entry
}

func NewSyncServer(logger *logrus.Entry, systemStub *system.SystemStub, syncManager SyncManager) *SyncServer {
	return &SyncServer{
		logger:      logger,
		systemStub:  systemStub,
		syncManager: syncManager,
	}
}

func (s *SyncServer) SyncNow(ctx context.Context, in *balena.SyncNowRequest) (*balena.SyncNowResponse, error) {
	serverUUID, err := primitive.ObjectIDFromHex(in.BalenaServerUUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid server uuid")
	}

	// Find the server to sync
	serverCollection := getBalenaServerCollection(s.systemStub)
	var server BalenaServerInMongo
	err = serverCollection.FindOne(ctx, bson.M{"_id": serverUUID}).Decode(&server)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "server not found")
		}

		err = errors.Join(errors.New("error while finding server in database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	syncLog, err := s.syncManager.SyncServer(ctx, server)
	if err != nil && errors.Is(err, ErrSyncInternal) {
		err = errors.Join(errors.New("error while syncing server"), err)
		s.logger.Error(err.Error())
		return &balena.SyncNowResponse{Log: syncLog.ToGRPCSyncLog()}, status.Error(codes.OK, "")
	}

	return &balena.SyncNowResponse{Log: syncLog.ToGRPCSyncLog()}, status.Error(codes.OK, "")
}
func (s *SyncServer) ListLog(in *balena.ListSyncLogRequest, out balena.BalenaSyncService_ListLogServer) error {
	serverUUID, err := primitive.ObjectIDFromHex(in.ServerUUID)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid server uuid")
	}

	ctx := out.Context()
	logCollection := getSyncLogCollection(s.systemStub)

	findOptions := options.Find().SetSort(bson.M{"timestamp": -1})
	if in.Skip != 0 {
		findOptions = findOptions.SetSkip(int64(in.Skip))
	}
	if in.Limit != 0 {
		findOptions = findOptions.SetLimit(int64(in.Limit))
	}

	cur, err := logCollection.Find(ctx, bson.M{"serverUUID": serverUUID}, findOptions)
	if err != nil {
		err = errors.Join(errors.New("error while setting up cursor for listing logs in database"), err)
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	defer cur.Close(context.Background())

	for cur.Next(ctx) {
		var log SyncLogInMongo
		err := cur.Decode(&log)
		if err != nil {
			err = errors.Join(errors.New("error while decoding log from database"), err)
			s.logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		err = out.Send(&balena.ListSyncLogResponse{
			Log: log.ToGRPCSyncLog(),
		})
		if err != nil {
			err = errors.Join(errors.New("error while sending log to client"), err)
			s.logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}

	if err := cur.Err(); err != nil {
		err = errors.Join(errors.New("error while iterating over logs in database"), err)
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	return status.Error(codes.OK, "")
}
func (s *SyncServer) CountLog(ctx context.Context, in *balena.CountSyncLogRequest) (*balena.CountSyncLogResponse, error) {
	serverUUID, err := primitive.ObjectIDFromHex(in.ServerUUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid server uuid")
	}
	collection := getSyncLogCollection(s.systemStub)
	count, err := collection.CountDocuments(ctx, bson.M{"serverUUID": serverUUID})
	if err != nil {
		err = errors.Join(errors.New("error while counting logs in database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.CountSyncLogResponse{TotalCount: uint64(count)}, status.Error(codes.OK, "")
}
func (s *SyncServer) GetLastSyncLog(ctx context.Context, in *balena.GetLastSyncLogRequest) (*balena.GetLastSyncLogResponse, error) {
	serverUUID, err := primitive.ObjectIDFromHex(in.ServerUUID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid server uuid")
	}
	collection := getSyncLogCollection(s.systemStub)

	var log SyncLogInMongo
	err = collection.FindOne(ctx, bson.M{"serverUUID": serverUUID}, options.FindOne().SetSort(bson.M{"timestamp": -1})).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "no logs found")
		}

		err = errors.Join(errors.New("error while finding last log in database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.GetLastSyncLogResponse{Log: log.ToGRPCSyncLog()}, status.Error(codes.OK, "")
}
