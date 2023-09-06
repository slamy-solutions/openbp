package models

import (
	"errors"
	"log/slog"
	"time"

	grpc "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CatalogFieldType string

const (
	CatalogFieldTypeInt           CatalogFieldType = "int"
	CatalogFieldTypeString        CatalogFieldType = "string"
	CatalogFieldTypeBool          CatalogFieldType = "bool"
	CatalogFieldTypeFloat         CatalogFieldType = "float"
	CatalogFieldTypeObject        CatalogFieldType = "object"
	CatalogFieldTypeTable         CatalogFieldType = "table"
	CatalogFieldTypeLinkToCatalog CatalogFieldType = "catalog"
)

type CatalogFieldBase struct {
	Type       CatalogFieldType `bson:"type"`
	Name       string           `bson:"name"`
	PublicName string           `bson:"publicName"`
}

type CatalogIntField struct {
	CatalogFieldBase `bson:",inline"`
}

type toGRPCFieldSchemaConvertable interface {
	ToGRPCFieldSchema() *grpc.FieldSchema
}

func (f *CatalogIntField) ToGRPCFieldSchema() *grpc.FieldSchema {
	return &grpc.FieldSchema{
		Schema:     &grpc.FieldSchema_IntData_{IntData: &grpc.FieldSchema_IntData{}},
		Name:       f.Name,
		PublicName: f.PublicName,
	}
}

type CatalogStringField struct {
	CatalogFieldBase `bson:",inline"`
}

func (f *CatalogStringField) ToGRPCFieldSchema() *grpc.FieldSchema {
	return &grpc.FieldSchema{
		Schema:     &grpc.FieldSchema_StringData_{StringData: &grpc.FieldSchema_StringData{}},
		Name:       f.Name,
		PublicName: f.PublicName,
	}
}

type CatalogBoolField struct {
	CatalogFieldBase `bson:",inline"`
}

func (f *CatalogBoolField) ToGRPCFieldSchema() *grpc.FieldSchema {
	return &grpc.FieldSchema{
		Schema:     &grpc.FieldSchema_BooleanData_{BooleanData: &grpc.FieldSchema_BooleanData{}},
		Name:       f.Name,
		PublicName: f.PublicName,
	}
}

type CatalogFloatField struct {
	CatalogFieldBase `bson:",inline"`
}

func (f *CatalogFloatField) ToGRPCFieldSchema() *grpc.FieldSchema {
	return &grpc.FieldSchema{
		Schema:     &grpc.FieldSchema_FloatData_{FloatData: &grpc.FieldSchema_FloatData{}},
		Name:       f.Name,
		PublicName: f.PublicName,
	}
}

type CatalogObjectField struct {
	CatalogFieldBase `bson:",inline"`

	Fields []interface{} `bson:"fields"`
}

func (f *CatalogObjectField) ToGRPCFieldSchema() *grpc.FieldSchema {
	var fields []*grpc.FieldSchema
	for _, field := range f.Fields {
		fields = append(fields, field.(toGRPCFieldSchemaConvertable).ToGRPCFieldSchema())
	}

	return &grpc.FieldSchema{
		Schema:     &grpc.FieldSchema_ObjectData_{ObjectData: &grpc.FieldSchema_ObjectData{Fields: fields}},
		Name:       f.Name,
		PublicName: f.PublicName,
	}
}

type CatalogTableField struct {
	CatalogFieldBase `bson:",inline"`

	Columns []interface{} `bson:"fields"`
}

func (f *CatalogTableField) ToGRPCFieldSchema() *grpc.FieldSchema {
	var columns []*grpc.FieldSchema
	for _, column := range f.Columns {
		columns = append(columns, column.(toGRPCFieldSchemaConvertable).ToGRPCFieldSchema())
	}

	return &grpc.FieldSchema{
		Schema:     &grpc.FieldSchema_TableData_{TableData: &grpc.FieldSchema_TableData{Columns: columns}},
		Name:       f.Name,
		PublicName: f.PublicName,
	}
}

type CatalogLinkToCatalogField struct {
	CatalogFieldBase `bson:",inline"`

	TargetCatalogName string `bson:"targetCatalogName"`
}

func (f *CatalogLinkToCatalogField) ToGRPCFieldSchema() *grpc.FieldSchema {
	return &grpc.FieldSchema{
		Schema:     &grpc.FieldSchema_CatalogLinkData_{CatalogLinkData: &grpc.FieldSchema_CatalogLinkData{CatalogName: f.TargetCatalogName}},
		Name:       f.Name,
		PublicName: f.PublicName,
	}
}

func parseRawCatalogFieldBaseSchema(raw bson.M, fieldType CatalogFieldType) (*CatalogFieldBase, error) {
	nameBSON, ok := raw["name"]
	if !ok {
		return nil, errors.New("missing name field")
	}
	name, ok := nameBSON.(string)
	if !ok {
		return nil, errors.New("invalid name field (not string format)")
	}

	publicNameBSON, ok := raw["publicName"]
	if !ok {
		return nil, errors.New("missing publicName field")
	}
	publicName, ok := publicNameBSON.(string)
	if !ok {
		return nil, errors.New("invalid publicName field (not string format)")
	}

	return &CatalogFieldBase{
		Type:       fieldType,
		Name:       name,
		PublicName: publicName,
	}, nil
}

func parseRawCatalogIntFieldSchema(raw bson.M) (*CatalogIntField, error) {
	catalogField, err := parseRawCatalogFieldBaseSchema(raw, CatalogFieldTypeInt)
	if err != nil {
		return nil, err
	}

	return &CatalogIntField{
		CatalogFieldBase: *catalogField,
	}, nil
}

func parseRawCatalogStringFieldSchema(raw bson.M) (*CatalogStringField, error) {
	catalogField, err := parseRawCatalogFieldBaseSchema(raw, CatalogFieldTypeInt)
	if err != nil {
		return nil, err
	}

	return &CatalogStringField{
		CatalogFieldBase: *catalogField,
	}, nil
}

func parseRawCatalogBoolFieldSchema(raw bson.M) (*CatalogBoolField, error) {
	catalogField, err := parseRawCatalogFieldBaseSchema(raw, CatalogFieldTypeInt)
	if err != nil {
		return nil, err
	}

	return &CatalogBoolField{
		CatalogFieldBase: *catalogField,
	}, nil
}

func parseRawCatalogFloatFieldSchema(raw bson.M) (*CatalogFloatField, error) {
	catalogField, err := parseRawCatalogFieldBaseSchema(raw, CatalogFieldTypeInt)
	if err != nil {
		return nil, err
	}

	return &CatalogFloatField{
		CatalogFieldBase: *catalogField,
	}, nil
}

func parseRawCatalogObjectFieldSchema(raw bson.M) (*CatalogObjectField, error) {
	catalogField, err := parseRawCatalogFieldBaseSchema(raw, CatalogFieldTypeInt)
	if err != nil {
		return nil, err
	}

	var fields []interface{}
	fieldsBSON, ok := raw["fields"]
	if !ok {
		return nil, errors.New("missing fields field")
	}
	fieldsBsonArray, ok := fieldsBSON.(bson.A)
	if !ok {
		return nil, errors.New("invalid fields field (not array format)")
	}
	for _, field := range fieldsBsonArray {
		fieldBsonMap, ok := field.(bson.M)
		if !ok {
			return nil, errors.New("invalid field (not map format)")
		}

		fieldSchema, err := parseRawCatalogFieldSchema(fieldBsonMap)
		if err != nil {
			return nil, errors.Join(errors.New("error while parsing fields of object "+catalogField.Name), err)
		}

		fields = append(fields, fieldSchema)
	}

	return &CatalogObjectField{
		CatalogFieldBase: *catalogField,
		Fields:           fields,
	}, nil
}

func parseRawCatalogTableFieldSchema(raw bson.M) (*CatalogTableField, error) {
	catalogField, err := parseRawCatalogFieldBaseSchema(raw, CatalogFieldTypeInt)
	if err != nil {
		return nil, err
	}

	var columns []interface{}
	columnsBSON, ok := raw["columns"]
	if !ok {
		return nil, errors.New("missing columns field")
	}
	columnsBsonArray, ok := columnsBSON.(bson.A)
	if !ok {
		return nil, errors.New("invalid columns field (not array format)")
	}
	for _, column := range columnsBsonArray {
		columnBsonMap, ok := column.(bson.M)
		if !ok {
			return nil, errors.New("invalid column (not map format)")
		}

		columnSchema, err := parseRawCatalogFieldSchema(columnBsonMap)
		if err != nil {
			return nil, errors.Join(errors.New("error while parsing columns of table "+catalogField.Name), err)
		}

		columns = append(columns, columnSchema)
	}

	return &CatalogTableField{
		CatalogFieldBase: *catalogField,
		Columns:          columns,
	}, nil
}

func parseCatalogLinkToCatalogFieldSchema(raw bson.M) (*CatalogLinkToCatalogField, error) {
	catalogField, err := parseRawCatalogFieldBaseSchema(raw, CatalogFieldTypeInt)
	if err != nil {
		return nil, err
	}

	targetCatalogNameBSON, ok := raw["targetCatalogName"]
	if !ok {
		return nil, errors.New("missing targetCatalogName field")
	}
	targetCatalogName, ok := targetCatalogNameBSON.(string)
	if !ok {
		return nil, errors.New("invalid targetCatalogName field (not string format)")
	}

	return &CatalogLinkToCatalogField{
		CatalogFieldBase:  *catalogField,
		TargetCatalogName: targetCatalogName,
	}, nil
}

func parseRawCatalogFieldSchema(raw bson.M) (interface{}, error) {
	typeBSON, ok := raw["type"]
	if !ok {
		return nil, errors.New("missing type field")
	}

	typeString, ok := typeBSON.(string)
	if !ok {
		return nil, errors.New("invalid type field (not string format)")
	}

	switch CatalogFieldType(typeString) {
	case CatalogFieldTypeInt:
		return parseRawCatalogIntFieldSchema(raw)
	case CatalogFieldTypeString:
		return parseRawCatalogStringFieldSchema(raw)
	case CatalogFieldTypeBool:
		return parseRawCatalogBoolFieldSchema(raw)
	case CatalogFieldTypeFloat:
		return parseRawCatalogFloatFieldSchema(raw)
	case CatalogFieldTypeObject:
		return parseRawCatalogObjectFieldSchema(raw)
	case CatalogFieldTypeTable:
		return parseRawCatalogTableFieldSchema(raw)
	case CatalogFieldTypeLinkToCatalog:
		return parseCatalogLinkToCatalogFieldSchema(raw)
	default:
		return nil, errors.New("invalid type field (unknown type)")
	}
}

type Catalog struct {
	Namespace string `bson:"-"`

	Name       string `bson:"name"`
	PublicName string `bson:"publicName"`

	Fields []interface{} `bson:"fields"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version int64     `bson:"version"`
}

func BSONToCatalog(raw bson.M) (*Catalog, error) {
	nameBSON, ok := raw["name"]
	if !ok {
		return nil, errors.New("missing name field")
	}
	name, ok := nameBSON.(string)
	if !ok {
		return nil, errors.New("invalid name field (not string format)")
	}

	publicNameBSON, ok := raw["publicName"]
	if !ok {
		return nil, errors.New("missing publicName field")
	}
	publicName, ok := publicNameBSON.(string)
	if !ok {
		return nil, errors.New("invalid publicName field (not string format)")
	}

	var fields []interface{}
	fieldsBSON, ok := raw["fields"]
	if !ok {
		return nil, errors.New("missing fields field")
	}
	fieldsBsonArray, ok := fieldsBSON.(bson.A)
	if !ok {
		return nil, errors.New("invalid fields field (not array format)")
	}
	for _, field := range fieldsBsonArray {
		fieldBsonMap, ok := field.(bson.M)
		if !ok {
			return nil, errors.New("invalid field (not map format)")
		}

		fieldSchema, err := parseRawCatalogFieldSchema(fieldBsonMap)
		if err != nil {
			return nil, err
		}

		fields = append(fields, fieldSchema)
	}

	createdBSON, ok := raw["created"]
	if !ok {
		return nil, errors.New("missing created field")
	}
	created, ok := createdBSON.(time.Time)
	if !ok {
		return nil, errors.New("invalid created field (not time.Time format)")
	}

	updatedBSON, ok := raw["updated"]
	if !ok {
		return nil, errors.New("missing updated field")
	}
	updated, ok := updatedBSON.(time.Time)
	if !ok {
		return nil, errors.New("invalid updated field (not time.Time format)")
	}

	versionBSON, ok := raw["version"]
	if !ok {
		return nil, errors.New("missing version field")
	}
	version, ok := versionBSON.(int64)
	if !ok {
		return nil, errors.New("invalid version field (not int64 format)")
	}

	return &Catalog{
		Name:       name,
		PublicName: publicName,
		Fields:     fields,
		Created:    created,
		Updated:    updated,
		Version:    version,
	}, nil
}

func (c *Catalog) ToGRPCatalog() *grpc.Catalog {
	var fields []*grpc.FieldSchema
	for _, field := range c.Fields {
		fields = append(fields, field.(toGRPCFieldSchemaConvertable).ToGRPCFieldSchema())
	}

	return &grpc.Catalog{
		Namespace:  c.Namespace,
		Name:       c.Name,
		PublicName: c.PublicName,

		Fields: fields,

		Created: timestamppb.New(c.Created),
		Updated: timestamppb.New(c.Updated),
		Version: c.Version,
	}
}

func (c *Catalog) ToLog(key string) slog.Attr {
	if key == "" {
		key = "catalog"
	}
	return slog.Group(
		key,
		slog.String("namespace", c.Namespace),
		slog.String("name", c.Name),
	)
}

func CatalogFieldFromGRPCFieldSchema(fieldSchema *grpc.FieldSchema) (interface{}, error) {
	switch fieldSchema.Schema.(type) {
	case *grpc.FieldSchema_IntData_:
		return &CatalogIntField{
			CatalogFieldBase: CatalogFieldBase{
				Type:       CatalogFieldTypeInt,
				Name:       fieldSchema.Name,
				PublicName: fieldSchema.PublicName,
			},
		}, nil
	case *grpc.FieldSchema_StringData_:
		return &CatalogStringField{
			CatalogFieldBase: CatalogFieldBase{
				Type:       CatalogFieldTypeString,
				Name:       fieldSchema.Name,
				PublicName: fieldSchema.PublicName,
			},
		}, nil
	case *grpc.FieldSchema_BooleanData_:
		return &CatalogBoolField{
			CatalogFieldBase: CatalogFieldBase{
				Type:       CatalogFieldTypeBool,
				Name:       fieldSchema.Name,
				PublicName: fieldSchema.PublicName,
			},
		}, nil
	case *grpc.FieldSchema_FloatData_:
		return &CatalogFloatField{
			CatalogFieldBase: CatalogFieldBase{
				Type:       CatalogFieldTypeFloat,
				Name:       fieldSchema.Name,
				PublicName: fieldSchema.PublicName,
			},
		}, nil
	case *grpc.FieldSchema_ObjectData_:
		var fields []interface{}
		for _, field := range fieldSchema.GetObjectData().Fields {
			fieldSchema, err := CatalogFieldFromGRPCFieldSchema(field)
			if err != nil {
				return nil, errors.Join(errors.New("error while parsing fields of object "), err)
			}

			fields = append(fields, fieldSchema)
		}

		return &CatalogObjectField{
			CatalogFieldBase: CatalogFieldBase{
				Type:       CatalogFieldTypeObject,
				Name:       fieldSchema.Name,
				PublicName: fieldSchema.PublicName,
			},
			Fields: fields,
		}, nil
	case *grpc.FieldSchema_TableData_:
		var columns []interface{}
		for _, column := range fieldSchema.GetTableData().Columns {
			columnSchema, err := CatalogFieldFromGRPCFieldSchema(column)
			if err != nil {
				return nil, errors.Join(errors.New("error while parsing columns of table "), err)
			}

			columns = append(columns, columnSchema)
		}

		return &CatalogTableField{
			CatalogFieldBase: CatalogFieldBase{
				Type:       CatalogFieldTypeTable,
				Name:       fieldSchema.Name,
				PublicName: fieldSchema.PublicName,
			},
			Columns: columns,
		}, nil
	case *grpc.FieldSchema_CatalogLinkData_:
		return &CatalogLinkToCatalogField{
			CatalogFieldBase: CatalogFieldBase{
				Type:       CatalogFieldTypeLinkToCatalog,
				Name:       fieldSchema.Name,
				PublicName: fieldSchema.PublicName,
			},
			TargetCatalogName: fieldSchema.GetCatalogLinkData().CatalogName,
		}, nil
	default:
		return nil, errors.New("invalid type field (unknown type)")
	}
}
