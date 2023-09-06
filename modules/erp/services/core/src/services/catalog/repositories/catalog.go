package repositories

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrCatalogAlreadyExists = errors.New("catalog already exists")
var ErrCatalogNotFound = errors.New("catalog not found")

type CatalogRepository interface {
	Create(ctx context.Context, catalog models.Catalog) (*models.Catalog, error)
	Get(ctx context.Context, namespace string, catalogName string) (*models.Catalog, error)
	GetAllInNamespace(ctx context.Context, namespace string) ([]models.Catalog, error)
	Update(ctx context.Context, catalog models.Catalog) (*models.Catalog, error)
	Delete(ctx context.Context, namespace string, catalogName string) (*models.Catalog, error)
}

type catalogRepository struct {
	logger *slog.Logger
	system *system.SystemStub
	tracer trace.Tracer
}

func NewCatalogRepository(logger *slog.Logger, system *system.SystemStub) CatalogRepository {

	return &catalogRepository{
		logger: logger,
		system: system,
		tracer: otel.GetTracerProvider().Tracer("catalog/repository/catalog"),
	}
}

func (r *catalogRepository) Create(ctx context.Context, catalog models.Catalog) (*models.Catalog, error) {
	ctx, span := r.tracer.Start(ctx, "createCatalog")
	defer span.End()

	creationTime := time.Now().UTC()
	catalog.Created = creationTime
	catalog.Updated = creationTime
	catalog.Version = 0

	collection := getCatalogCollection(r.system, catalog.Namespace)
	_, err := collection.InsertOne(ctx, catalog)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, ErrCatalogAlreadyExists
		}

		err = errors.Join(errors.New("failed to insert catalog to the database"), err)
		r.logger.Error("failed to insert catalog to the database", slog.String("err", err.Error()))
		return nil, err
	}

	err = createEntryCollection(ctx, r.system, catalog.Namespace, catalog.Name)
	if err != nil {
		_, err := collection.DeleteOne(ctx, bson.M{"name": catalog.Name})
		if err != nil {
			err = errors.Join(errors.New("failed to delete catalog from the database"), err)
			r.logger.Error("failed to delete catalog from the database", slog.String("err", err.Error()))
		}

		err = errors.Join(errors.New("failed to create catalog entry collection"), err)
		r.logger.Error("failed to create catalog entry collection", slog.String("err", err.Error()))
		return nil, err
	}

	r.logger.Info("Catalog created", catalog.ToLog("catalog"))
	return &catalog, nil
}

func (r *catalogRepository) Get(ctx context.Context, namespace string, catalogName string) (*models.Catalog, error) {
	ctx, span := r.tracer.Start(ctx, "getCatalog")
	defer span.End()

	collection := getCatalogCollection(r.system, namespace)

	var catalog models.Catalog
	err := collection.FindOne(ctx, bson.M{"name": catalogName}).Decode(&catalog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCatalogNotFound
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, ErrCatalogNotFound
			}
		}

		err = errors.Join(errors.New("failed to get catalog from the database"), err)
		r.logger.Error("failed to get catalog from the database", slog.String("err", err.Error()))
		return nil, err
	}
	catalog.Namespace = namespace

	return &catalog, nil
}

func (r *catalogRepository) GetAllInNamespace(ctx context.Context, namespace string) ([]models.Catalog, error) {
	ctx, span := r.tracer.Start(ctx, "getAllCatalogsInNamespace")
	defer span.End()

	collection := getCatalogCollection(r.system, namespace)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return []models.Catalog{}, nil
			}
		}

		err = errors.Join(errors.New("failed to get catalogs from the database"), err)
		r.logger.Error("failed to get catalogs from the database", slog.String("err", err.Error()))
		return nil, err
	}

	var catalogs []models.Catalog
	err = cursor.All(ctx, &catalogs)
	if err != nil {
		err = errors.Join(errors.New("failed to get catalogs from the database"), err)
		r.logger.Error("failed to get catalogs from the database", slog.String("err", err.Error()))
		return nil, err
	}

	for i := range catalogs {
		catalogs[i].Namespace = namespace
	}

	return catalogs, nil
}

func (r *catalogRepository) Update(ctx context.Context, catalog models.Catalog) (*models.Catalog, error) {
	ctx, span := r.tracer.Start(ctx, "updateCatalog")
	defer span.End()

	catalog.Updated = time.Now().UTC()
	catalog.Version++

	collection := getCatalogCollection(r.system, catalog.Namespace)
	var newCatalog models.Catalog
	err := collection.FindOneAndUpdate(ctx, bson.M{"name": catalog.Name}, bson.M{
		"$set": bson.M{
			"publicName": catalog.PublicName,
			"fields":     catalog.Fields,
		},
		"$inc": bson.M{
			"version": 1,
		},
		"$currentDate": bson.M{
			"updated": bson.M{"$type": "timestamp"},
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&newCatalog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCatalogNotFound
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, ErrCatalogNotFound
			}
		}

		err = errors.Join(errors.New("failed to update catalog in the database"), err)
		r.logger.Error("failed to update catalog in the database", slog.String("err", err.Error()))
		return nil, err
	}
	newCatalog.Namespace = catalog.Namespace

	r.logger.Info("Catalog updated", newCatalog.ToLog("catalog"))
	return &newCatalog, nil
}

func (r *catalogRepository) Delete(ctx context.Context, namespace string, catalogName string) (*models.Catalog, error) {
	ctx, span := r.tracer.Start(ctx, "deleteCatalog")
	defer span.End()

	collection := getCatalogCollection(r.system, namespace)
	var catalog models.Catalog
	err := collection.FindOneAndDelete(ctx, bson.M{"name": catalogName}).Decode(&catalog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCatalogNotFound
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, ErrCatalogNotFound
			}
		}

		err = errors.Join(errors.New("failed to delete catalog from the database"), err)
		r.logger.Error("failed to delete catalog from the database", slog.String("err", err.Error()))
		return nil, err
	}

	err = dropEntryCollection(ctx, r.system, namespace, catalogName)
	if err != nil {
		err = errors.Join(errors.New("failed to drop catalog entry collection"), err)
		r.logger.Error("failed to drop catalog entry collection", slog.String("err", err.Error()))
		return nil, err
	}

	r.logger.Info("Catalog deleted", catalog.ToLog("catalog"))
	return &catalog, nil
}
