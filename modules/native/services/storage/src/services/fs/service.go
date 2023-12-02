package fs

import (
	"context"
	"io"
	"log/slog"

	fsGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/fs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	fsGRPC.UnimplementedFSServiceServer

	repository *FileRepository
	logger     *slog.Logger
}

func NewService(repository *FileRepository, logger *slog.Logger) fsGRPC.FSServiceServer {
	return &service{
		repository: repository,
		logger:     logger.With("service", "fs"),
	}
}

func (s *service) CreateFile(ctx context.Context, in *fsGRPC.CreateFileRequest) (*fsGRPC.CreateFileResponse, error) {
	bucketUUID, err := primitive.ObjectIDFromHex(in.Bucket)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid bucket id")
	}

	file, err := s.repository.Create(ctx, in.Namespace, bucketUUID, in.Path, in.MimeType)
	if err != nil {
		if err == ErrFileAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "file already exists")
		}

		s.logger.ErrorContext(ctx, "failed to create file", "error", err)
		return nil, status.Error(codes.Internal, "failed to create file: "+err.Error())
	}

	return &fsGRPC.CreateFileResponse{
		File: file.ToGRPC(),
	}, status.Error(codes.OK, "")
}

type uploadStreamReader struct {
	firstChunk       []byte
	firstChunkSended bool
	srv              fsGRPC.FSService_UploadFileServer
}

func (r *uploadStreamReader) Read(p []byte) (n int, err error) {
	if !r.firstChunkSended {
		r.firstChunkSended = true
		return copy(p, r.firstChunk), nil
	}

	frame, err := r.srv.Recv()
	if err != nil {
		return 0, err
	}

	return copy(p, frame.DataChunk), nil
}

func (r *uploadStreamReader) WriteTo(w io.Writer) (n int64, err error) {
	var bytesSended int64 = 0

	if !r.firstChunkSended {
		written, err := w.Write(r.firstChunk)
		if err != nil {
			return 0, err
		}
		bytesSended += int64(written)
	}

	for {
		frame, err := r.srv.Recv()
		if err != nil {
			if err == io.EOF {
				return bytesSended, nil
			}

			return 0, err
		}

		written, err := w.Write(frame.DataChunk)
		if err != nil {
			return 0, err
		}
		bytesSended += int64(written)
	}
}

func (s *service) UploadFile(srv fsGRPC.FSService_UploadFileServer) error {
	ctx := srv.Context()

	rq, err := srv.Recv()
	if err != nil {
		if err == io.EOF {
			return status.Error(codes.InvalidArgument, "empty file. no first chunk received")
		}

		return status.Error(codes.Internal, "failed to receive first chunk: "+err.Error())
	}

	fileUUID, err := primitive.ObjectIDFromHex(rq.Uuid)
	if err != nil {
		return status.Error(codes.InvalidArgument, "invalid file id")
	}

	fileInfo, err := s.repository.Upload(ctx, rq.Namespace, fileUUID, &uploadStreamReader{
		firstChunk:       rq.DataChunk,
		firstChunkSended: false,
		srv:              srv,
	})
	if err != nil {
		if err == ErrFileNotFound {
			return status.Error(codes.NotFound, "file not found")
		}

		s.logger.ErrorContext(ctx, "failed to upload file", "error", err)
		return status.Error(codes.Internal, "failed to upload file: "+err.Error())
	}

	return srv.SendAndClose(&fsGRPC.UploadFileResponse{
		File: fileInfo.ToGRPC(),
	})
}
func (s *service) StatFile(ctx context.Context, in *fsGRPC.StatFileRequest) (*fsGRPC.StatFileResponse, error) {
	fileUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "file not found. invalid file id")
	}

	file, err := s.repository.Stat(ctx, in.Namespace, fileUUID)
	if err != nil {
		if err == ErrFileNotFound {
			return nil, status.Error(codes.NotFound, "file not found")
		}

		s.logger.ErrorContext(ctx, "failed to stat file", "error", err)
		return nil, status.Error(codes.Internal, "failed to stat file: "+err.Error())
	}

	return &fsGRPC.StatFileResponse{
		File: file.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) UpdateFile(ctx context.Context, in *fsGRPC.UpdateFileRequest) (*fsGRPC.UpdateFileResponse, error) {
	fileUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "file not found. invalid file id")
	}

	file, err := s.repository.Update(ctx, in.Namespace, fileUUID, in.Path, in.MimeType)
	if err != nil {
		if err == ErrFileNotFound {
			return nil, status.Error(codes.NotFound, "file not found")
		}

		if err == ErrFileAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "file already exists")
		}

		s.logger.ErrorContext(ctx, "failed to update file", "error", err)
		return nil, status.Error(codes.Internal, "failed to update file: "+err.Error())
	}

	return &fsGRPC.UpdateFileResponse{
		File: file.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) DeleteFile(ctx context.Context, in *fsGRPC.DeleteFileRequest) (*fsGRPC.DeleteFileResponse, error) {
	fileUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "file not found. invalid file id")
	}

	file, err := s.repository.Delete(ctx, in.Namespace, fileUUID)
	if err != nil {
		if err == ErrFileNotFound {
			return nil, status.Error(codes.NotFound, "file not found")
		}

		s.logger.ErrorContext(ctx, "failed to delete file", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete file: "+err.Error())
	}

	return &fsGRPC.DeleteFileResponse{
		File: file.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *service) ListFiles(in *fsGRPC.ListFilesRequest, out fsGRPC.FSService_ListFilesServer) error {
	ctx := out.Context()

	bucketUUID, err := primitive.ObjectIDFromHex(in.Bucket)
	if err != nil {
		return status.Error(codes.OK, "Invalid bucket UUID")
	}

	filesCursor, err := s.repository.List(ctx, in.Namespace, bucketUUID, int64(in.Skip), int64(in.Limit))
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list files", "error", err)
		return status.Error(codes.Internal, "failed to list files: "+err.Error())
	}
	defer filesCursor.Close()

	for {
		file, err := filesCursor.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			s.logger.ErrorContext(ctx, "failed to get next file", "error", err)
			return status.Error(codes.Internal, "failed to list files: "+err.Error())
		}

		err = out.Send(&fsGRPC.ListFilesResponse{
			File: file.ToGRPC(),
		})
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to send file", "error", err)
			return status.Error(codes.Internal, "failed to list files: "+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *service) CountFiles(ctx context.Context, in *fsGRPC.CountFilesRequest) (*fsGRPC.CountFilesResponse, error) {
	bucketUUID, err := primitive.ObjectIDFromHex(in.Bucket)
	if err != nil {
		return &fsGRPC.CountFilesResponse{Count: 0}, status.Error(codes.OK, "Invalid bucket UUID")
	}

	count, err := s.repository.Count(ctx, in.Namespace, bucketUUID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to count files", "error", err)
		return nil, status.Error(codes.Internal, "failed to count files: "+err.Error())
	}

	return &fsGRPC.CountFilesResponse{
		Count: uint32(count),
	}, status.Error(codes.OK, "")
}
func (s *service) Download(in *fsGRPC.DownloadFileRequest, out fsGRPC.FSService_DownloadServer) error {
	ctx := out.Context()

	fileUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return status.Error(codes.NotFound, "file not found. invalid file id")
	}

	fileInfo, reader, err := s.repository.Download(ctx, in.Namespace, fileUUID, int64(in.Seek))
	if err != nil {
		if err == ErrFileNotFound {
			return status.Error(codes.NotFound, "file not found")
		}

		s.logger.ErrorContext(ctx, "failed to download file", "error", err)
		return status.Error(codes.Internal, "failed to download file: "+err.Error())
	}
	defer reader.Close()

	toRead := fileInfo.Size - int64(in.Seek)
	if in.Limit > 0 && toRead > int64(in.Limit) {
		toRead = int64(in.Limit)
	}

	bufSize := 1024 * 32
	buf := make([]byte, bufSize)
	for {
		chunkSizeToRead := bufSize
		if toRead < int64(bufSize) {
			chunkSizeToRead = int(toRead)
		}

		readedBytes, err := reader.Read(buf[:chunkSizeToRead])
		if err != nil {
			if err == io.EOF {
				break
			}

			s.logger.ErrorContext(ctx, "failed to read file", "error", err)
			return status.Error(codes.Internal, "failed to download file: "+err.Error())
		}

		err = out.Send(&fsGRPC.DownloadFileResponse{
			DataChunk: buf[:readedBytes],
		})
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to send file", "error", err)
			return status.Error(codes.Internal, "failed to download file: "+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *service) DownloadByPath(in *fsGRPC.DownloadFileByPathRequest, out fsGRPC.FSService_DownloadByPathServer) error {
	bucketUUID, err := primitive.ObjectIDFromHex(in.Bucket)
	if err != nil {
		return status.Error(codes.NotFound, "file not found. invalid bucket uuid")
	}

	ctx := out.Context()

	fileInfo, reader, err := s.repository.DownloadByPath(ctx, in.Namespace, bucketUUID, in.Path, int64(in.Seek))
	if err != nil {
		if err == ErrFileNotFound {
			return status.Error(codes.NotFound, "file not found")
		}

		s.logger.ErrorContext(ctx, "failed to download file", "error", err)
		return status.Error(codes.Internal, "failed to download file: "+err.Error())
	}
	defer reader.Close()

	toRead := fileInfo.Size - int64(in.Seek)
	if in.Limit > 0 && toRead > int64(in.Limit) {
		toRead = int64(in.Limit)
	}

	bufSize := 1024 * 32
	buf := make([]byte, bufSize)
	for {
		chunkSizeToRead := bufSize
		if toRead < int64(bufSize) {
			chunkSizeToRead = int(toRead)
		}

		readedBytes, err := reader.Read(buf[:chunkSizeToRead])
		if err != nil {
			if err == io.EOF {
				break
			}

			s.logger.ErrorContext(ctx, "failed to read file", "error", err)
			return status.Error(codes.Internal, "failed to download file: "+err.Error())
		}

		err = out.Send(&fsGRPC.DownloadFileByPathResponse{
			DataChunk: buf[:readedBytes],
		})
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to send file", "error", err)
			return status.Error(codes.Internal, "failed to download file: "+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *service) DownloadDirect(in *fsGRPC.DownloadDirectFileRequest, out fsGRPC.FSService_DownloadDirectServer) error {
	ctx := out.Context()

	fileUUID, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return status.Error(codes.NotFound, "file not found. invalid file id")
	}

	fileInfo, reader, err := s.repository.DownloadDirect(ctx, in.Namespace, fileUUID, in.DirectDownloadSecret, int64(in.Seek))
	if err != nil {
		if err == ErrFileNotFound {
			return status.Error(codes.NotFound, "file not found")
		}

		if err == ErrDirectDownloadSecretInvalid {
			return status.Error(codes.PermissionDenied, "direct download secret invalid")
		}

		s.logger.ErrorContext(ctx, "failed to download file", "error", err)
		return status.Error(codes.Internal, "failed to download file: "+err.Error())
	}
	defer reader.Close()

	toRead := fileInfo.Size - int64(in.Seek)
	if in.Limit > 0 && toRead > int64(in.Limit) {
		toRead = int64(in.Limit)
	}

	bufSize := 1024 * 32
	buf := make([]byte, bufSize)
	for {
		chunkSizeToRead := bufSize
		if toRead < int64(bufSize) {
			chunkSizeToRead = int(toRead)
		}

		readedBytes, err := reader.Read(buf[:chunkSizeToRead])
		if err != nil {
			if err == io.EOF {
				break
			}

			s.logger.ErrorContext(ctx, "failed to read file", "error", err)
			return status.Error(codes.Internal, "failed to download file: "+err.Error())
		}

		err = out.Send(&fsGRPC.DownloadDirectFileResponse{
			DataChunk: buf[:readedBytes],
		})
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to send file", "error", err)
			return status.Error(codes.Internal, "failed to download file: "+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *service) DownloadDirectByPath(in *fsGRPC.DownloadDirectFileByPathRequest, out fsGRPC.FSService_DownloadDirectByPathServer) error {
	bucketUUID, err := primitive.ObjectIDFromHex(in.Bucket)
	if err != nil {
		return status.Error(codes.NotFound, "file not found. invalid bucket uuid")
	}

	ctx := out.Context()

	fileInfo, reader, err := s.repository.DownloadDirectByPath(ctx, in.Namespace, bucketUUID, in.Path, in.DirectDownloadSecret, int64(in.Seek))
	if err != nil {
		if err == ErrFileNotFound {
			return status.Error(codes.NotFound, "file not found")
		}

		if err == ErrDirectDownloadSecretInvalid {
			return status.Error(codes.PermissionDenied, "direct download secret invalid")
		}

		s.logger.ErrorContext(ctx, "failed to download file", "error", err)
		return status.Error(codes.Internal, "failed to download file: "+err.Error())
	}
	defer reader.Close()

	toRead := fileInfo.Size - int64(in.Seek)
	if in.Limit > 0 && toRead > int64(in.Limit) {
		toRead = int64(in.Limit)
	}

	bufSize := 1024 * 32
	buf := make([]byte, bufSize)
	for {
		chunkSizeToRead := bufSize
		if toRead < int64(bufSize) {
			chunkSizeToRead = int(toRead)
		}

		readedBytes, err := reader.Read(buf[:chunkSizeToRead])
		if err != nil {
			if err == io.EOF {
				break
			}

			s.logger.ErrorContext(ctx, "failed to read file", "error", err)
			return status.Error(codes.Internal, "failed to download file: "+err.Error())
		}

		err = out.Send(&fsGRPC.DownloadDirectFileByPathResponse{
			DataChunk: buf[:readedBytes],
		})
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to send file", "error", err)
			return status.Error(codes.Internal, "failed to download file: "+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
