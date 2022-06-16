package services

import (
	"fmt"
	grpc "slamy/opencrm/native/catalog/grpc/native_catalog_grpc"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	PropertyTypeInt    int = 0
	PropertyTypeString int = 1
	PropertyTypeBool   int = 2
)

func GRPCPropertiesToBSON(schemas []*grpc.PropertySchema) bson.A {
	resultProperties := bson.A{}
	for _, property := range schemas {
		resultProperty := bson.D{
			bson.E{Key: "name", Value: property.Name},
			bson.E{Key: "publicName", Value: property.PublicName},
		}

		if el := property.GetIntProperty(); el != nil {
			resultProperty = append(resultProperty, bson.E{Key: "_propertyType", Value: PropertyTypeInt})
			resultProperty = append(resultProperty, bson.E{Key: "defaultValue", Value: el.DefaultValue})
		} else if el := property.GetStringProperty(); el != nil {
			resultProperty = append(resultProperty, bson.E{Key: "_propertyType", Value: PropertyTypeString})
			resultProperty = append(resultProperty, bson.E{Key: "defaultValue", Value: el.DefaultValue})
		} else if el := property.GetBooleanProperty(); el != nil {
			resultProperty = append(resultProperty, bson.E{Key: "_propertyType", Value: PropertyTypeBool})
			resultProperty = append(resultProperty, bson.E{Key: "defaultValue", Value: el.DefaultValue})
		}

		resultProperties = append(resultProperties, resultProperty)
	}
	return resultProperties
}

func GRPCPropertiesFromBSON(data bson.A) []*grpc.PropertySchema {

}

func CatalogFromBSON(data bson.D, namespace string) *grpc.Catalog {
	m := data.Map()
	return &grpc.Catalog{
		Namespace:  namespace,
		Name:       m["name"].(string),
		PublicName: m["publicName"].(string),
		Properties: GRPCPropertiesFromBSON(m["properties"].(bson.A)),
		XCreated:   timestamppb.New(m["_created"].(time.Time)),
		XUpdated:   timestamppb.New(m["_updated"].(time.Time)),
		XVersion:   m["_version"].(int64),
	}
}

func CatalogIndexFromBSON(data bson.D, namespace string) *grpc.CatalogIndex {

}

func CatalogIndexToBSON(indexData *grpc.CatalogIndex) bson.D {

}

func CatalogIndexFieldsToBSON(indexData *grpc.CatalogIndex) bson.D {
	var result []bson.D
	for _, field := range indexData.Fields {
		grpc.CatalogIndex_IndexField_ASCENDING
		result = append(result, &bson.E{Key: field.Name, Value: field.Type})
	}
}

func MakeDBName(namespace string) string {
	return fmt.Sprintf("opencrm_namespace_%s", namespace)
}

/*
type IntPropertyData struct {
	DefaultValue int64 `bson:"defaultValue"`
}

type StringPropertyData struct {
	DefaultValue int64 `bson:"defaultValue"`
}

type BoolPropertyData struct {
	DefaultValue int64 `bson:"defaultValue"`
}

type PropertySchema struct {
	Name       string `bson:"name"`
	PublicName string `bson:"publicName"`

	IntData    IntPropertyData    `bson:"int,omitempty,inline"`
	StringData StringPropertyData `bson:"string,omitempty,inline"`
	BoolData   BoolPropertyData   `bson:"bool,omitempty,inline"`
}

type catalog struct {
	Name       string `bson:"name"`
	PublicName string `bson:"publicName"`

	Properties PropertySchema `bson:"properties"`

	Created time.Time `bson:"_created"`
	Updated time.Time `bson:"_updated"`
	Version int64     `bson:"_version"`
}

type Catalog interface {
	ToBSON() error
}

func NewCatalogFromGRPC(name string, publicName string)
*/
