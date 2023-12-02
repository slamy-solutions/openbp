package fs

import (
	"errors"
	"time"

	fsGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/fs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrFailedToGetFileBucket = errors.New("failed to get file bucket")
var ErrFileAlreadyExists = errors.New("file at path already exists")
var ErrFileNotFound = errors.New("file not found")
var ErrDirectDownloadSecretInvalid = errors.New("direct download secret invalid")
var ErrFilePathInvalid = errors.New("file path invalid")

type File struct {
	Namespace         string             `bson:"-"`
	Bucket            primitive.ObjectID `bson:"bucket"`
	UUID              primitive.ObjectID `bson:"_id"`
	Path              string             `bson:"path"`
	BaseDirectoryPath string             `bson:"baseDirectoryPath"`

	DirectDownloadSecret string             `bson:"directDownloadSecret"`
	MimeType             string             `bson:"mimeType"`
	Size                 int64              `bson:"size"`
	GridFSFile           primitive.ObjectID `bson:"gridfsFile"`

	Created time.Time `bson:"_created"`
	Updated time.Time `bson:"_updated"`
	Version int       `bson:"_version"`
}

func (f *File) ToGRPC() *fsGRPC.File {
	return &fsGRPC.File{
		Namespace: f.Namespace,
		Bucket:    f.Bucket.Hex(),
		Uuid:      f.UUID.Hex(),
		Path:      f.Path,

		DirectDownloadSecret: f.DirectDownloadSecret,
		MimeType:             f.MimeType,
		Size:                 f.Size,

		XCreated: timestamppb.New(f.Created),
		XUpdated: timestamppb.New(f.Updated),
		XVersion: int64(f.Version),
	}
}

type Directory struct {
	Namespace         string             `bson:"-"`
	Bucket            primitive.ObjectID `bson:"bucket"`
	Path              string             `bson:"path"`
	BaseDirectoryPath string             `bson:"baseDirectoryPath"`
}
