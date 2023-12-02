package bucket

import (
	"errors"
	"log/slog"
	"time"

	bucketGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/bucket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrBucketNameInvalid = errors.New("bucket name is invalid")
var ErrBucketAlreadyExists = errors.New("bucket already exists")
var ErrBucketNotFound = errors.New("bucket not found")

type Bucket struct {
	Namespace string             `bson:"-"`
	UUID      primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`

	Hidden bool `bson:"hidden"`

	Created time.Time `bson:"_created"`
	Updated time.Time `bson:"_updated"`
	Version int64     `bson:"_version"`
}

func (b *Bucket) ToSlogAttr(groupName string) slog.Attr {
	if groupName == "" {
		groupName = "bucket"
	}

	return slog.Group(groupName,
		slog.String("namespace", b.Namespace),
		slog.String("uuid", b.UUID.Hex()),
		slog.String("name", b.Name),
		slog.Bool("hidden", b.Hidden),
		slog.Time("created", b.Created),
		slog.Time("updated", b.Updated),
		slog.Int64("version", b.Version),
	)
}

func (b *Bucket) ToGRPC() *bucketGRPC.Bucket {
	return &bucketGRPC.Bucket{
		Namespace: b.Namespace,
		Uuid:      b.UUID.Hex(),
		Name:      b.Name,

		Hidden: b.Hidden,

		XCreated: timestamppb.New(b.Created),
		XUpdated: timestamppb.New(b.Updated),
		XVersion: b.Version,
	}
}
