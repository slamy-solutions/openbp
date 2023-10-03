package runtime

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/golang/protobuf/proto"
	grpcRuntime "github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/runtime"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ManagerRuntimeServer struct {
	grpcRuntime.UnimplementedRuntimeServiceServer

	systemStub *system.SystemStub
	logger     *slog.Logger
}

func NewManagerRuntimeServer(logger *slog.Logger, systemStub *system.SystemStub) *ManagerRuntimeServer {
	return &ManagerRuntimeServer{
		systemStub: systemStub,
		logger:     logger,
	}
}

func (s *ManagerRuntimeServer) GetRuntimesForNamespace(ctx context.Context, in *grpcRuntime.GetRuntimesForNamespaceReqeust) (*grpcRuntime.GetRuntimesForNamespaceResponse, error) {
	logger := s.logger.With(slog.String("endpoint", "GetRuntimesForNamespace"), slog.String("namespace", in.Namespace))

	collection := GetRuntimeCollection(s.systemStub)

	var runtimes []RuntimeInMongo
	cur, err := collection.Find(ctx, bson.M{"namespace": in.Namespace})
	if err != nil {
		err = errors.Join(errors.New("failed to find runtimes for namespace"), err)
		logger.Error(err.Error())
		return nil, err
	}

	err = cur.All(ctx, &runtimes)
	if err != nil {
		err = errors.Join(errors.New("failed to decode runtimes for namespace"), err)
		logger.Error(err.Error())
		return nil, err
	}

	var grpcRuntimes []*grpcRuntime.Runtime
	for _, runtime := range runtimes {
		grpcRuntimes = append(grpcRuntimes, runtime.ToGRPCRuntime())
	}

	return &grpcRuntime.GetRuntimesForNamespaceResponse{
		Runtimes: grpcRuntimes,
	}, status.Error(codes.OK, "")
}
func (s *ManagerRuntimeServer) GetRuntime(ctx context.Context, in *grpcRuntime.GetRuntimeRequest) (*grpcRuntime.GetRuntimeResponse, error) {
	logger := s.logger.With(slog.String("endpoint", "GetRuntime"))

	collection := GetRuntimeCollection(s.systemStub)

	var runtime RuntimeInMongo
	err := collection.FindOne(ctx, bson.M{"namespace": in.Namespace, "name": in.Name}).Decode(&runtime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "runtime not found")
		}

		err = errors.Join(errors.New("failed to find runtime"), err)
		logger.Error(err.Error())
		return nil, err
	}

	return &grpcRuntime.GetRuntimeResponse{
		Runtime: runtime.ToGRPCRuntime(),
	}, status.Error(codes.OK, "")
}
func (s *ManagerRuntimeServer) CreateRuntime(ctx context.Context, in *grpcRuntime.CreateRuntimeRequest) (*grpcRuntime.CreateRuntimeResponse, error) {
	logger := s.logger.With(slog.String("endpoint", "CreateRuntime"))

	collection := GetRuntimeCollection(s.systemStub)

	_, err := collection.InsertOne(ctx, RuntimeInMongo{
		Namespace: in.Runtime.Namespace,
		Name:      in.Runtime.Name,
		Run:       in.Runtime.Run,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Errorf(codes.AlreadyExists, "runtime already exists")
		}

		err = errors.Join(errors.New("failed to insert runtime"), err)
		logger.Error(err.Error())
		return nil, err
	}

	// Publish evet about created runtime
	runtimeAsBinary, err := proto.Marshal(in.Runtime)
	if err != nil {
		err = errors.Join(errors.New("failed to publish runtime created event. failed to marshal runtime"), err)
		logger.Error(err.Error())
	} else {
		err = s.systemStub.Nats.Publish(runtimeCreatedEventName, runtimeAsBinary)
		if err != nil {
			err = errors.Join(errors.New("failed to publish runtime created event"), err)
			logger.Error(err.Error())
		}
	}

	return nil, status.Errorf(codes.Unimplemented, "method CreateRuntime not implemented")
}
func (s *ManagerRuntimeServer) UpdateRuntime(ctx context.Context, in *grpcRuntime.UpdateRuntimeRequest) (*grpcRuntime.UpdateRuntimeResponse, error) {
	logger := s.logger.With(slog.String("endpoint", "UpdateRuntime"))

	collection := GetRuntimeCollection(s.systemStub)

	var runtime RuntimeInMongo
	err := collection.FindOneAndUpdate(ctx, bson.M{"namespace": in.Namespace, "name": in.Name}, bson.M{"$set": bson.M{"run": in.NewRun}}).Decode(&runtime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "runtime not found")
		}

		err = errors.Join(errors.New("failed to update runtime"), err)
		logger.Error(err.Error())
		return nil, err
	}

	// Publish evet about changed runtime
	runtimeAsBinary, err := proto.Marshal(runtime.ToGRPCRuntime())
	if err != nil {
		err = errors.Join(errors.New("failed to publish runtime updated event. failed to marshal runtime"), err)
		logger.Error(err.Error())
	} else {
		err = s.systemStub.Nats.Publish(runtimeUpdatedEventName, runtimeAsBinary)
		if err != nil {
			err = errors.Join(errors.New("failed to publish runtime updated event"), err)
			logger.Error(err.Error())
		}
	}

	return nil, status.Error(codes.OK, "")
}
func (s *ManagerRuntimeServer) DeleteRuntime(ctx context.Context, in *grpcRuntime.DeleteRuntimeRequest) (*grpcRuntime.DeleteRuntimeResponse, error) {
	logger := s.logger.With(slog.String("endpoint", "DeleteRuntime"))

	collection := GetRuntimeCollection(s.systemStub)

	var runtime RuntimeInMongo
	err := collection.FindOneAndDelete(ctx, bson.M{"namespace": in.Namespace, "name": in.Name}).Decode(&runtime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "runtime not found")
		}

		err = errors.Join(errors.New("failed to delete runtime"), err)
		logger.Error(err.Error())
		return nil, err
	}

	// Publish evet about deleted runtime
	runtimeAsBinary, err := proto.Marshal(runtime.ToGRPCRuntime())
	if err != nil {
		err = errors.Join(errors.New("failed to publish runtime deleted event. failed to marshal runtime"), err)
		logger.Error(err.Error())
	} else {
		err = s.systemStub.Nats.Publish(runtimeDeletedEventName, runtimeAsBinary)
		if err != nil {
			err = errors.Join(errors.New("failed to publish runtime deleted event"), err)
			logger.Error(err.Error())
		}
	}

	return &grpcRuntime.DeleteRuntimeResponse{}, status.Error(codes.OK, "")
}
func (s *ManagerRuntimeServer) UploadRuntimeBinary(srv grpcRuntime.RuntimeService_UploadRuntimeBinaryServer) error {
	logger := s.logger.With(slog.String("endpoint", "UploadRuntimeBinary"))
	ctx := srv.Context()

	collection := GetRuntimeCollection(s.systemStub)
	var bucket *gridfs.Bucket
	var uploadStream *gridfs.UploadStream
	var runtime RuntimeInMongo
	successfullyUploaded := false

	defer func() {
		if uploadStream != nil {
			uploadStream.Close()
		}
	}()

	defer func() {
		if !successfullyUploaded {
			if bucket != nil {
				_ = bucket.Delete(uploadStream.FileID)
			}
		}
	}()

	loadRuntimeData := func(namespace string, name string) error {
		err := collection.FindOne(ctx, bson.M{"namespace": namespace, "name": name}).Decode(&runtime)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return status.Errorf(codes.NotFound, "runtime not found")
			}

			err = errors.Join(errors.New("failed to find runtime"), err)
			logger.Error(err.Error())
			return err
		}

		bucket, err = GetRuntimeDataBucket(runtime.Namespace, s.systemStub)
		if err != nil {
			err = errors.Join(errors.New("failed to get runtime data bucket"), err)
			logger.Error(err.Error())
			return err
		}

		uploadStream, err = bucket.OpenUploadStream(runtime.Name)
		if err != nil {
			err = errors.Join(errors.New("failed to open upload stream"), err)
			logger.Error(err.Error())
			return err
		}

		return nil
	}

	for {
		rcv, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Join(errors.New("failed to receive package"), err)
			logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		if bucket == nil {
			err := loadRuntimeData(rcv.Namespace, rcv.Name)
			if err != nil {
				// Error is already logged
				return err
			}
		}

		_, err = uploadStream.Write(rcv.Binary)
		if err != nil {
			err = errors.Join(errors.New("failed to write package to database"), err)
			logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}

	if bucket == nil {
		return status.Errorf(codes.InvalidArgument, "no runtime data received")
	}

	err := collection.FindOneAndUpdate(ctx, bson.M{"namespace": runtime.Namespace, "name": runtime.Name}, bson.M{"$set": bson.M{"binaryFile": uploadStream.FileID}}).Err()
	if err != nil {
		err = errors.Join(errors.New("failed to update runtime"), err)
		logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	successfullyUploaded = true

	// Publish evet about changed binary data
	runtimeAsBinary, err := proto.Marshal(runtime.ToGRPCRuntime())
	if err != nil {
		err = errors.Join(errors.New("failed to publish runtime binary updated event. failed to marshal runtime"), err)
		logger.Error(err.Error())
	} else {
		err = s.systemStub.Nats.Publish(runtimeBinaryUpdatedEventName, runtimeAsBinary)
		if err != nil {
			err = errors.Join(errors.New("failed to publish runtime binary updated event"), err)
			logger.Error(err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *ManagerRuntimeServer) DownloadRuntimeBinary(in *grpcRuntime.DownloadRuntimeBinaryRequest, out grpcRuntime.RuntimeService_DownloadRuntimeBinaryServer) error {
	logger := s.logger.With(slog.String("endpoint", "DownloadRuntimeBinary"))

	collection := GetRuntimeCollection(s.systemStub)
	var runtime RuntimeInMongo
	err := collection.FindOne(out.Context(), bson.M{"namespace": in.Namespace, "name": in.Name}).Decode(&runtime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return status.Errorf(codes.NotFound, "runtime not found")
		}

		err = errors.Join(errors.New("failed to find runtime"), err)
		logger.Error(err.Error())
		return err
	}

	if runtime.BinaryFile == primitive.NilObjectID {
		return status.Errorf(codes.InvalidArgument, "runtime has no binary data")
	}

	bucket, err := GetRuntimeDataBucket(runtime.Namespace, s.systemStub)
	if err != nil {
		err = errors.Join(errors.New("failed to get runtime data bucket"), err)
		logger.Error(err.Error())
		return err
	}

	downloadStream, err := bucket.OpenDownloadStream(runtime.BinaryFile)
	if err != nil {
		err = errors.Join(errors.New("failed to open download stream"), err)
		logger.Error(err.Error())
		return err
	}
	defer downloadStream.Close()

	buf := make([]byte, 1024*1024) // 1 megabyte
	for {
		n, err := downloadStream.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Join(errors.New("failed to read package from database"), err)
			logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		err = out.Send(&grpcRuntime.DownloadRuntimeBinaryResponse{
			Binary: buf[:n],
		})
		if err != nil {
			err = errors.Join(errors.New("failed to send package"), err)
			logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}
	return status.Error(codes.OK, "")
}
