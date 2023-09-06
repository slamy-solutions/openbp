package repositories

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrCatalogEntryNotFound = errors.New("catalog entry not found")
var ErrCatalogDataSchemeInvalid = errors.New("catalog data scheme invalid")

type CatalogEntryStream interface {
	Next(ctx context.Context) (*models.CatalogEntry, error)
	Close(ctx context.Context) error
}

type CatalogEntryRepository interface {
	Create(ctx context.Context, catalogEntry models.CatalogEntry) (*models.CatalogEntry, error)
	Get(ctx context.Context, namespace string, catalogName string, uuid primitive.ObjectID) (*models.CatalogEntry, error)
	List(ctx context.Context, namespace string, catalogName string, skip int64, limit int64) (CatalogEntryStream, error)
	Count(ctx context.Context, namespace string, catalogName string) (int64, error)
	Update(ctx context.Context, catalogEntry models.CatalogEntry) (*models.CatalogEntry, error)
	Delete(ctx context.Context, namespace string, catalogName string, uuid primitive.ObjectID) (*models.CatalogEntry, error)
}

type catalogEntryRepository struct {
	logger *slog.Logger
	system *system.SystemStub
	tracer trace.Tracer
}

func NewCatalogEntryRepository(logger *slog.Logger, system *system.SystemStub) CatalogEntryRepository {
	return &catalogEntryRepository{
		logger: logger,
		system: system,
		tracer: otel.GetTracerProvider().Tracer("catalog/repository/entry"),
	}
}

func (r *catalogEntryRepository) Create(ctx context.Context, catalogEntry models.CatalogEntry) (*models.CatalogEntry, error) {
	ctx, span := r.tracer.Start(ctx, "createCatalogEntry")
	defer span.End()

	// Make sure the UUID is empty. MongoDB will generate a new one.
	catalogEntry.UUID = primitive.NilObjectID

	collection := getCatalogEntryCollection(r.system, catalogEntry.Namespace, catalogEntry.Catalog)
	insertResult, err := collection.InsertOne(ctx, catalogEntry)
	if err != nil {
		err = errors.Join(errors.New("failed to insert catalog entry to the database"), err)
		r.logger.Error("failed to insert catalog entry to the database", slog.String("err", err.Error()))
		return nil, err
	}

	catalogEntry.UUID = insertResult.InsertedID.(primitive.ObjectID)
	r.logger.Info("catalog entry created", catalogEntry.ToLog("catalogEntry"))
	return &catalogEntry, nil
}

func (r *catalogEntryRepository) Get(ctx context.Context, namespace string, catalogName string, uuid primitive.ObjectID) (*models.CatalogEntry, error) {
	ctx, span := r.tracer.Start(ctx, "getCatalogEntry")
	defer span.End()

	collection := getCatalogEntryCollection(r.system, namespace, catalogName)
	var catalogEntry models.CatalogEntry
	err := collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&catalogEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCatalogEntryNotFound
		}

		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, ErrCatalogEntryNotFound
			}
		}

		err = errors.Join(errors.New("failed to get catalog entry from the database"), err)
		r.logger.Error("failed to get catalog entry from the database", slog.String("err", err.Error()))
		return nil, err
	}

	return &catalogEntry, nil
}

type catalogEntryStream struct {
	cursor *mongo.Cursor
}

func (s *catalogEntryStream) Next(ctx context.Context) (*models.CatalogEntry, error) {
	if !s.cursor.Next(ctx) {
		err := s.cursor.Err()
		if err == nil {
			return nil, io.EOF
		}
		return nil, err
	}

	var catalogEntry models.CatalogEntry
	err := s.cursor.Decode(&catalogEntry)
	if err != nil {
		err = errors.Join(errors.New("failed to decode catalog entry from the database"), err)
		return nil, err
	}

	return &catalogEntry, nil
}

func (s *catalogEntryStream) Close(ctx context.Context) error {
	return s.cursor.Close(ctx)
}

type emptyCatalogEntryStream struct{}

func (s *emptyCatalogEntryStream) Next(ctx context.Context) (*models.CatalogEntry, error) {
	return nil, io.EOF
}
func (s *emptyCatalogEntryStream) Close(ctx context.Context) error {
	return nil
}

func (r *catalogEntryRepository) List(ctx context.Context, namespace string, catalogName string, skip int64, limit int64) (CatalogEntryStream, error) {
	ctx, span := r.tracer.Start(ctx, "listCatalogEntries")
	defer span.End()

	collection := getCatalogEntryCollection(r.system, namespace, catalogName)
	options := options.Find()
	if skip > 0 {
		options = options.SetSkip(skip)
	}
	if limit > 0 {
		options = options.SetLimit(limit)
	}
	cursor, err := collection.Find(ctx, bson.M{}, options)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &emptyCatalogEntryStream{}, nil
			}
		}

		err = errors.Join(errors.New("failed to list catalog entries from the database"), err)
		r.logger.Error("failed to list catalog entries from the database", slog.String("err", err.Error()))
		return nil, err
	}

	return &catalogEntryStream{
		cursor: cursor,
	}, nil
}

func (r *catalogEntryRepository) Count(ctx context.Context, namespace string, catalogName string) (int64, error) {
	ctx, span := r.tracer.Start(ctx, "countCatalogEntries")
	defer span.End()

	collection := getCatalogEntryCollection(r.system, namespace, catalogName)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return 0, nil
			}
		}

		err = errors.Join(errors.New("failed to count catalog entries in the database"), err)
		r.logger.Error("failed to count catalog entries in the database", slog.String("err", err.Error()))
		return 0, err
	}

	return count, nil
}

func (r *catalogEntryRepository) Update(ctx context.Context, catalogEntry models.CatalogEntry) (*models.CatalogEntry, error) {
	ctx, span := r.tracer.Start(ctx, "updateCatalogEntry")
	defer span.End()

	collection := getCatalogEntryCollection(r.system, catalogEntry.Namespace, catalogEntry.Catalog)
	udpate := bson.M{
		"$set": bson.M{
			"data": catalogEntry.Data,
		},
		"$inc": bson.M{
			"version": 1,
		},
		"$currentDate": bson.M{
			"updated": bson.M{"$type": "timestamp"},
		},
	}
	var newCatalogEntry models.CatalogEntry
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": catalogEntry.UUID}, udpate, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&newCatalogEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCatalogEntryNotFound
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, ErrCatalogEntryNotFound
			}
		}

		err = errors.Join(errors.New("failed to update catalog entry in the database"), err)
		r.logger.Error("failed to update catalog entry in the database", slog.String("err", err.Error()))
		return nil, err
	}

	r.logger.Info("catalog entry updated", newCatalogEntry.ToLog("catalogEntry"))
	return &newCatalogEntry, nil
}

func (r *catalogEntryRepository) Delete(ctx context.Context, namespace string, catalogName string, uuid primitive.ObjectID) (*models.CatalogEntry, error) {
	ctx, span := r.tracer.Start(ctx, "deleteCatalogEntry")
	defer span.End()

	collection := getCatalogEntryCollection(r.system, namespace, catalogName)
	var catalogEntry models.CatalogEntry
	err := collection.FindOneAndDelete(ctx, bson.M{"_id": uuid}).Decode(&catalogEntry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCatalogEntryNotFound
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, ErrCatalogEntryNotFound
			}
		}

		err = errors.Join(errors.New("failed to delete catalog entry from the database"), err)
		r.logger.Error("failed to delete catalog entry from the database", slog.String("err", err.Error()))
		return nil, err
	}

	r.logger.Info("catalog entry deleted", catalogEntry.ToLog("catalogEntry"))
	return &catalogEntry, nil
}
