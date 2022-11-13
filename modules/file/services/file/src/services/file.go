package services

import (
	"context"
	"crypto/sha512"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/slamy-solutions/openbp/modules/system/libs/go/cache"

	fileGRPC "github.com/slamy-solutions/openbp/modules/native/services/file/src/grpc/native_file"
	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/services/file/src/grpc/native_namespace"
)

type FileServer struct {
	fileGRPC.UnimplementedFileServiceServer

	mongoClient     *mongo.Client
	cacheClient     cache.Cache
	namespaceClient namespaceGRPC.NamespaceServiceClient
}

type fileInMongo struct {
	Id       primitive.ObjectID `bson:"_id"`
	DataId   primitive.ObjectID `bson:"dataId"`
	Readonly bool               `bson:"readonly"`
	Size     uint64             `bson:"size"`
	MimeType string             `bson:"mimeType"`

	SHA512HashCalculated bool   `bson:"sha512HashCalculated"`
	SHA512Hash           []byte `bson:"sha512hash"`

	DisableCache bool `bson:"disableCache"`
	ForceCaching bool `bson:"forceCaching"`

	XCreated time.Time `bson:"_created"`
	XUpdated time.Time `bson:"_updated"`
	XVersion int64     `bson:"_version"`
}

func (f *fileInMongo) ToProtoFile(namespace string) *fileGRPC.File {
	return &fileGRPC.File{
		Namespace:            namespace,
		Uuid:                 f.Id.Hex(),
		Readonly:             f.Readonly,
		MimeType:             f.MimeType,
		Size:                 f.Size,
		SHA512HashCalculated: f.SHA512HashCalculated,
		SHA512Hash:           f.SHA512Hash,
		DisableCache:         f.DisableCache,
		ForceCaching:         f.ForceCaching,
		XCreated:             timestamppb.New(f.XCreated),
		XUpdated:             timestamppb.New(f.XUpdated),
		XVersion:             f.XVersion,
	}
}

func New(mongoClient *mongo.Client, cacheClient cache.Cache, namespaceClient namespaceGRPC.NamespaceServiceClient) *FileServer {
	return &FileServer{
		mongoClient:     mongoClient,
		cacheClient:     cacheClient,
		namespaceClient: namespaceClient,
	}
}

func (s *FileServer) getDB(namespace string) *mongo.Database {
	return s.mongoClient.Database(fmt.Sprintf("openbp_namespace_%s", namespace))
}

func (s *FileServer) getFileCollection(namespace string) *mongo.Collection {
	return s.getDB(namespace).Collection("native_file")
}

func (s *FileServer) getBucket(namespace string) (*gridfs.Bucket, error) {
	return gridfs.NewBucket(
		s.getDB(namespace),
		options.GridFSBucket().SetName("native_file_bucket"),
	)
}

func (s *FileServer) Create(in fileGRPC.FileService_CreateServer) error {
	ctx := in.Context()

	// Receive first package with file information
	pkg, err := in.Recv()
	if err != nil {
		if err == io.EOF {
			return status.Error(grpccodes.Aborted, "Package with file information wasnt received")
		}
		return status.Error(grpccodes.Internal, err.Error())
	}
	info := pkg.GetInfo()
	if info == nil {
		return status.Error(grpccodes.DataLoss, "First package must be file info")
	}

	// Check if namespace exist.
	existResponse, err := s.namespaceClient.Exists(ctx, &namespaceGRPC.IsNamespaceExistRequest{Name: info.Namespace, UseCache: true})
	if err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}
	if !existResponse.Exist {
		return status.Error(grpccodes.FailedPrecondition, "Namespace does not exist")
	}

	bucket, err := s.getBucket(info.Namespace)
	if err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}

	hasher := sha512.New()

	// Write file to the GridFS
	stream, err := bucket.OpenUploadStream("")
	if err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}

	totalSize := 0
	for {
		pkg, err := in.Recv()
		if err != nil {
			if err == io.EOF {
				// All data received.
				stream.Close()
				break
			}
			//
			stream.Close()
			bucket.Delete(stream.FileID)
			return status.Error(grpccodes.DataLoss, err.Error())
		}
		chunk := pkg.GetChunk()
		if chunk == nil {
			stream.Close()
			bucket.Delete(stream.FileID)
			return status.Error(grpccodes.DataLoss, "All the packages except first must have chunk data")
		}

		n, err := stream.Write(chunk.Data)
		if err != nil {
			stream.Close()
			bucket.Delete(stream.FileID)
			return status.Error(grpccodes.Internal, err.Error())
		}
		totalSize += n

		// Hasher never returns error
		hasher.Write(chunk.Data)
	}

	sha512sum := hasher.Sum([]byte{})

	// Create DB entry with file information and link to the data
	creationTime := time.Now()
	infoToSave := fileInMongo{
		DataId:   stream.FileID.(primitive.ObjectID),
		Readonly: info.Readonly,
		Size:     uint64(totalSize),
		MimeType: info.MimeType,

		SHA512HashCalculated: true,
		SHA512Hash:           sha512sum,

		DisableCache: info.DisableCache,
		ForceCaching: info.ForceCaching,

		XCreated: creationTime,
		XUpdated: creationTime,
		XVersion: 0,
	}

	insertResult, err := s.getFileCollection(info.Namespace).InsertOne(ctx, infoToSave)
	if err != nil {
		bucket.Delete(stream.FileID)
		return status.Error(grpccodes.Internal, err.Error())
	}
	infoToSave.Id = insertResult.InsertedID.(primitive.ObjectID)

	err = in.SendAndClose(&fileGRPC.FileCreateResponse{File: infoToSave.ToProtoFile(info.Namespace)})
	return err
}
func (s *FileServer) Exists(ctx context.Context, in *fileGRPC.FileExistRequest) (*fileGRPC.FileExistResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Cant convert UUID in object ID. Bad format.")
	}

	var info fileInMongo
	err = s.getFileCollection(in.Namespace).FindOne(ctx, bson.M{"_id": id}).Decode(&info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &fileGRPC.FileExistResponse{Exist: false}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &fileGRPC.FileExistResponse{Exist: true}, status.Error(grpccodes.OK, "")
}
func (s *FileServer) Stat(ctx context.Context, in *fileGRPC.StatFileRequest) (*fileGRPC.StatFileResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Cant convert UUID in object ID. Bad format.")
	}

	var info fileInMongo
	err = s.getFileCollection(in.Namespace).FindOne(ctx, bson.M{"_id": id}).Decode(&info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, err.Error())
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &fileGRPC.StatFileResponse{File: info.ToProtoFile(in.Namespace)}, status.Error(grpccodes.OK, "")
}
func (s *FileServer) Read(in *fileGRPC.ReadFileRequest, out fileGRPC.FileService_ReadServer) error {
	ctx := out.Context()

	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return status.Error(grpccodes.InvalidArgument, "Cant convert UUID in object ID. Bad format.")
	}

	var info fileInMongo
	err = s.getFileCollection(in.Namespace).FindOne(ctx, bson.M{"_id": id}).Decode(&info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return status.Error(grpccodes.NotFound, err.Error())
		}
		return status.Error(grpccodes.Internal, err.Error())
	}

	if in.Start > info.Size {
		return status.Error(grpccodes.InvalidArgument, "Start index is greater than file size")
	}
	if in.Start+in.ToRead > info.Size {
		return status.Error(grpccodes.InvalidArgument, "Number of bytes to read is out of file range. (start + toRead > fileSize)")
	}

	bytesToTransfer := info.Size - in.Start
	if in.ToRead > 0 {
		bytesToTransfer = in.ToRead
	}

	bucket, err := s.getBucket(in.Namespace)
	if err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}

	stream, err := bucket.OpenDownloadStream(info.DataId)
	if err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}
	defer stream.Close()

	if in.Start > 0 {
		_, err = stream.Skip(int64(in.Start))
		if err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
	}
	buf := make([]byte, 32*1024)
	transfered := uint64(0)
	for {
		readed, err := stream.Read(buf)
		if uint64(readed) > bytesToTransfer-transfered {
			readed = int(bytesToTransfer - transfered)
		}
		transfered += uint64(readed)
		if err != nil {
			if err == io.EOF {
				return status.Error(codes.OK, "")
			}
			return status.Error(codes.Internal, err.Error())
		}
		out.Send(&fileGRPC.ReadFileResponse{
			TotalSize:  bytesToTransfer,
			Transfered: transfered,
			ChunkStart: in.Start + transfered - uint64(readed),
			Chunk:      buf[:readed],
		})
		if transfered == bytesToTransfer {
			break
		}
	}

	return status.Error(codes.OK, "")
}
func (s *FileServer) Delete(ctx context.Context, in *fileGRPC.DeleteFileRequest) (*fileGRPC.DeleteFileResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Cant convert UUID in object ID. Bad format.")
	}

	var info fileInMongo
	err = s.getFileCollection(in.Namespace).FindOne(ctx, bson.M{"_id": id}).Decode(&info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &fileGRPC.DeleteFileResponse{}, nil
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	bucket, err := s.getBucket(in.Namespace)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	err = bucket.Delete(info.DataId)
	if err != nil {
		if err != gridfs.ErrFileNotFound {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}
	}

	_, err = s.getFileCollection(in.Namespace).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &fileGRPC.DeleteFileResponse{}, nil
}
func (s *FileServer) CalculateSHA512(ctx context.Context, in *fileGRPC.CalculateFileSHA512Request) (*fileGRPC.CalculateFileSHA512Response, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Cant convert UUID in object ID. Bad format.")
	}

	var info fileInMongo
	err = s.getFileCollection(in.Namespace).FindOne(ctx, bson.M{"_id": id}).Decode(&info)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(grpccodes.NotFound, "File not found")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	// Dont calculate SHA512 second time if it is already calculated
	if info.SHA512HashCalculated {
		return &fileGRPC.CalculateFileSHA512Response{SHA512: info.SHA512Hash}, nil
	}

	bucket, err := s.getBucket(in.Namespace)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	stream, err := bucket.OpenDownloadStream(info.DataId)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	defer stream.Close()

	hasher := sha512.New()
	_, err = io.Copy(hasher, stream)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	sha512Hash := hasher.Sum([]byte{})

	_, err = s.getFileCollection(in.Namespace).UpdateOne(ctx, bson.M{"_id": id}, bson.M{"sha512": sha512Hash})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &fileGRPC.CalculateFileSHA512Response{SHA512: sha512Hash}, nil
}
