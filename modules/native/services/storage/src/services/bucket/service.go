package bucket

import (
	"context"
	"errors"
	"io"
	"log/slog"

	bucketGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/bucket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	bucketGRPC.UnimplementedBucketServiceServer

	repository *BucketRepository
	logger     *slog.Logger
}

func NewService(repository *BucketRepository, logger *slog.Logger) bucketGRPC.BucketServiceServer {
	return &service{
		repository: repository,
		logger:     logger.With("service", "bucket"),
	}
}

func (s *service) Create(ctx context.Context, in *bucketGRPC.CreateBucketRequest) (*bucketGRPC.CreateBucketResponse, error) {
	bucket, err := s.repository.Create(ctx, in.Namespace, in.Name, in.Hidden)
	if err != nil {
		if err == ErrBucketAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "bucket with same name already exists")
		}

		if err == ErrBucketNameInvalid {
			return nil, status.Error(codes.InvalidArgument, "bucket name is invalid")
		}

		s.logger.ErrorContext(ctx, "failed to create bucket", "error", err)
		return nil, status.Error(codes.Internal, "failed to create bucket: "+err.Error())
	}

	return &bucketGRPC.CreateBucketResponse{
		Bucket: bucket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) Ensure(ctx context.Context, in *bucketGRPC.EnsureBucketRequest) (*bucketGRPC.EnsureBucketResponse, error) {
	bucket, err := s.repository.Ensure(ctx, in.Namespace, in.Name, in.Hidden)
	if err != nil {
		if err == ErrBucketNameInvalid {
			return nil, status.Error(codes.InvalidArgument, "bucket name is invalid")
		}

		s.logger.ErrorContext(ctx, "failed to ensure bucket", "error", err)
		return nil, status.Error(codes.Internal, "failed to ensure bucket: "+err.Error())
	}

	return &bucketGRPC.EnsureBucketResponse{
		Bucket: bucket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) Get(ctx context.Context, in *bucketGRPC.GetBucketRequest) (*bucketGRPC.GetBucketResponse, error) {
	bucket, err := s.repository.Get(ctx, in.Namespace, in.Name)
	if err != nil {
		if err == ErrBucketNotFound {
			return nil, status.Error(codes.NotFound, "bucket not found")
		}

		s.logger.ErrorContext(ctx, "failed to get bucket", "error", err)
		return nil, status.Error(codes.Internal, "failed to get bucket: "+err.Error())
	}

	return &bucketGRPC.GetBucketResponse{
		Bucket: bucket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) GetByUUID(ctx context.Context, in *bucketGRPC.GetBucketByUUIDRequest) (*bucketGRPC.GetBucketByUUIDResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "bucket not found: invalid uuid")
	}

	bucket, err := s.repository.GetByUUID(ctx, in.Namespace, uuid)
	if err != nil {
		if err == ErrBucketNotFound {
			return nil, status.Error(codes.NotFound, "bucket not found")
		}

		s.logger.ErrorContext(ctx, "failed to get bucket", "error", err)
		return nil, status.Error(codes.Internal, "failed to get bucket: "+err.Error())
	}

	return &bucketGRPC.GetBucketByUUIDResponse{
		Bucket: bucket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) List(in *bucketGRPC.ListBucketsRequest, out bucketGRPC.BucketService_ListServer) error {
	ctx := out.Context()
	cursor, err := s.repository.List(ctx, in.Namespace, int32(in.Skip), int32(in.Limit))
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list buckets", "error", err)
		return status.Error(codes.Internal, "failed to list buckets: "+err.Error())
	}
	defer cursor.Close()

	for {
		bucket, err := cursor.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Join(errors.New("failed to get next bucket"), err)
			s.logger.ErrorContext(ctx, "failed to list buckets", "error", err)
			return status.Error(codes.Internal, "failed to list buckets: "+err.Error())
		}

		err = out.Send(&bucketGRPC.ListBucketsResponse{
			Bucket: bucket.ToGRPC(),
		})
		if err != nil {
			err = errors.Join(errors.New("failed to send bucket"), err)
			s.logger.ErrorContext(ctx, "failed to list buckets", "error", err)
			return status.Error(codes.Internal, "failed to list buckets: "+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *service) Update(ctx context.Context, in *bucketGRPC.UpdateBucketRequest) (*bucketGRPC.UpdateBucketResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "bucket not found: invalid uuid")
	}

	bucket, err := s.repository.Update(ctx, in.Namespace, uuid, in.Name, in.Hidden)
	if err != nil {
		if err == ErrBucketNotFound {
			return nil, status.Error(codes.NotFound, "bucket not found")
		}
		if err == ErrBucketAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "bucket with same name already exists")
		}
		if err == ErrBucketNameInvalid {
			return nil, status.Error(codes.InvalidArgument, "bucket name is invalid")
		}

		s.logger.ErrorContext(ctx, "failed to update bucket", "error", err)
		return nil, status.Error(codes.Internal, "failed to update bucket: "+err.Error())
	}

	return &bucketGRPC.UpdateBucketResponse{
		Bucket: bucket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) Delete(ctx context.Context, in *bucketGRPC.DeleteBucketRequest) (*bucketGRPC.DeleteBucketResponse, error) {
	bucket, err := s.repository.Delete(ctx, in.Namespace, in.Name)
	if err != nil {
		if err == ErrBucketNotFound {
			return nil, status.Error(codes.NotFound, "bucket not found")
		}

		s.logger.ErrorContext(ctx, "failed to delete bucket", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete bucket: "+err.Error())
	}

	return &bucketGRPC.DeleteBucketResponse{
		Bucket: bucket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) DeleteByUUID(ctx context.Context, in *bucketGRPC.DeleteBucketByUUIDRequest) (*bucketGRPC.DeleteBucketByUUIDResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "bucket not found: invalid uuid")
	}

	bucket, err := s.repository.DeleteByUUID(ctx, in.Namespace, uuid)
	if err != nil {
		if err == ErrBucketNotFound {
			return nil, status.Error(codes.NotFound, "bucket not found")
		}

		s.logger.ErrorContext(ctx, "failed to delete bucket", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete bucket: "+err.Error())
	}

	return &bucketGRPC.DeleteBucketByUUIDResponse{
		Bucket: bucket.ToGRPC(),
	}, status.Error(codes.OK, "")
}
