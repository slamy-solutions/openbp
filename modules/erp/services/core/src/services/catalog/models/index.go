package models

import (
	"log/slog"
	"strings"

	grpc "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CatalogIndexType string

const (
	catalogIndexNamePrefix                      = "openbp_erp_catalog_"
	CatalogIndexTypeHashed     CatalogIndexType = "hashed"
	CatalogIndexTypeAscending  CatalogIndexType = "ascending"
	CatalogIndexTypeDescending CatalogIndexType = "descending"
)

type CatalogIndexField struct {
	Name string
	Type CatalogIndexType
}

type CatalogIndex struct {
	Namespace string
	Catalog   string

	Name   string
	Unique bool

	Fields []CatalogIndexField
}

// Will return nil if the index is not a catalog index (for example defined by internal system services).
func CatalogIndexFromMongoModel(index *mongo.IndexModel, namespace string, catalog string) *CatalogIndex {
	name := ""
	if index.Options.Name != nil {
		name = *index.Options.Name
	}

	if name == "" || !strings.HasPrefix(name, catalogIndexNamePrefix) {
		return nil
	}

	var fields []CatalogIndexField
	for _, key := range index.Keys.(bson.D) {
		switch key.Value {
		case 1:
			fields = append(fields, CatalogIndexField{
				Name: key.Key,
				Type: CatalogIndexTypeAscending,
			})
		case -1:
			fields = append(fields, CatalogIndexField{
				Name: key.Key,
				Type: CatalogIndexTypeDescending,
			})
		case "hashed":
			fields = append(fields, CatalogIndexField{
				Name: key.Key,
				Type: CatalogIndexTypeHashed,
			})
		}
	}

	unique := false
	if index.Options.Unique != nil {
		unique = *index.Options.Unique
	}

	return &CatalogIndex{
		Namespace: namespace,
		Catalog:   catalog,

		Name:   name,
		Unique: unique,

		Fields: fields,
	}
}

func (c *CatalogIndex) ToMongoIndex() mongo.IndexModel {
	var keys bson.D
	for _, field := range c.Fields {
		switch field.Type {
		case CatalogIndexTypeHashed:
			keys = append(keys, bson.E{Key: field.Name, Value: "hashed"})
		case CatalogIndexTypeAscending:
			keys = append(keys, bson.E{Key: field.Name, Value: 1})
		case CatalogIndexTypeDescending:
			keys = append(keys, bson.E{Key: field.Name, Value: -1})
		}
	}

	return mongo.IndexModel{
		Keys:    keys,
		Options: options.Index().SetName(catalogIndexNamePrefix + c.Name).SetUnique(c.Unique),
	}
}

func (c *CatalogIndex) ToLog(key string) slog.Attr {
	if key == "" {
		key = "catalogIndex"
	}
	return slog.Group(
		key,
		slog.String("namespace", c.Namespace),
		slog.String("catalog", c.Catalog),
		slog.String("name", c.Name),
	)
}

func (c *CatalogIndex) ToGRPCatalogIndex() *grpc.CatalogIndex {
	var fields []*grpc.CatalogIndex_IndexField
	for _, field := range c.Fields {
		fieldType := grpc.CatalogIndex_IndexField_ASCENDING
		switch field.Type {
		case CatalogIndexTypeHashed:
			fieldType = grpc.CatalogIndex_IndexField_HASHED
		case CatalogIndexTypeAscending:
			fieldType = grpc.CatalogIndex_IndexField_ASCENDING
		case CatalogIndexTypeDescending:
			fieldType = grpc.CatalogIndex_IndexField_DESCENDING
		}

		fields = append(fields, &grpc.CatalogIndex_IndexField{
			Name: field.Name,
			Type: fieldType,
		})
	}

	return &grpc.CatalogIndex{
		Name:      c.Name,
		Unique:    c.Unique,
		Fields:    fields,
		Namespace: c.Namespace,
		Catalog:   c.Catalog,
	}
}

func CatalogIndexFromGRPC(in *grpc.CatalogIndex) CatalogIndex {
	var fields []CatalogIndexField
	for _, field := range in.Fields {
		fieldType := CatalogIndexTypeAscending
		switch field.Type {
		case grpc.CatalogIndex_IndexField_HASHED:
			fieldType = CatalogIndexTypeHashed
		case grpc.CatalogIndex_IndexField_ASCENDING:
			fieldType = CatalogIndexTypeAscending
		case grpc.CatalogIndex_IndexField_DESCENDING:
			fieldType = CatalogIndexTypeDescending
		}

		fields = append(fields, CatalogIndexField{
			Name: field.Name,
			Type: fieldType,
		})
	}

	return CatalogIndex{
		Name:      in.Name,
		Unique:    in.Unique,
		Fields:    fields,
		Namespace: in.Namespace,
		Catalog:   in.Catalog,
	}
}
