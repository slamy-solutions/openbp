package models

import (
	"encoding/json"
	"log/slog"
	"time"

	grpc "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CatalogEntry struct {
	Namespace string `bson:"-"`
	Catalog   string `bson:"-"`

	UUID primitive.ObjectID `bson:"uuid,omitempty"`

	Data interface{} `bson:"data"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version int64     `bson:"version"`
}

func (c *CatalogEntry) ToLog(key string) slog.Attr {
	if key == "" {
		key = "catalogEntry"
	}
	return slog.Group(
		key,
		slog.String("uuid", c.UUID.Hex()),
		slog.String("namespace", c.Namespace),
		slog.String("catalog", c.Catalog),
	)
}

func (c *CatalogEntry) ToGRPCCatalogEntry() *grpc.CatalogEntry {
	jsonData, _ := json.Marshal(c.Data)

	return &grpc.CatalogEntry{
		Uuid:      c.UUID.Hex(),
		Data:      jsonData,
		Namespace: c.Namespace,
		Catalog:   c.Catalog,
		Created:   timestamppb.New(c.Created),
		Updated:   timestamppb.New(c.Updated),
		Version:   c.Version,
	}
}
