package repositories

import (
	"context"
	"errors"
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type CatalogIndexRepository interface {
	Ensure(ctx context.Context, index models.CatalogIndex) error
	GetAllForCatalog(ctx context.Context, namespace string, catalogName string) ([]models.CatalogIndex, error)
	Delete(ctx context.Context, namespace string, catalogName string, indexName string) error
}

type catalogIndexRepository struct {
	logger *slog.Logger
	system *system.SystemStub
	tracer trace.Tracer
}

func NewCatalogIndexRepository(logger *slog.Logger, system *system.SystemStub) CatalogIndexRepository {
	return &catalogIndexRepository{
		logger: logger,
		system: system,
		tracer: otel.GetTracerProvider().Tracer("catalog/repository/index"),
	}
}

func (r *catalogIndexRepository) Ensure(ctx context.Context, index models.CatalogIndex) error {
	ctx, span := r.tracer.Start(ctx, "ensureCatalogIndex")
	defer span.End()

	collection := getCatalogEntryCollection(r.system, index.Namespace, index.Catalog)
	_, err := collection.Indexes().CreateOne(ctx, index.ToMongoIndex())
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return ErrCatalogNotFound
			}
		}

		err = errors.Join(errors.New("failed to ensure catalog index"), err)
		r.logger.Error(err.Error())
		return err
	}

	r.logger.Info("catalog index ensured", index.ToLog("index"))
	return nil
}

func (r *catalogIndexRepository) GetAllForCatalog(ctx context.Context, namespace string, catalogName string) ([]models.CatalogIndex, error) {
	ctx, span := r.tracer.Start(ctx, "getAllCatalogIndexes")
	defer span.End()

	collection := getCatalogEntryCollection(r.system, namespace, catalogName)
	indexes, err := collection.Indexes().List(ctx)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, ErrCatalogNotFound
			}
		}

		err = errors.Join(errors.New("failed to get catalog indexes"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	var result []models.CatalogIndex
	for indexes.Next(ctx) {
		var index mongo.IndexModel
		err := indexes.Decode(index)
		if err != nil {
			err = errors.Join(errors.New("failed to decode catalog index"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		catalogIndex := models.CatalogIndexFromMongoModel(&index, namespace, catalogName)
		result = append(result, *catalogIndex)
	}

	return result, nil
}

func (r *catalogIndexRepository) Delete(ctx context.Context, namespace string, catalogName string, indexName string) error {
	ctx, span := r.tracer.Start(ctx, "deleteCatalogIndex")
	defer span.End()

	collection := getCatalogEntryCollection(r.system, namespace, catalogName)
	_, err := collection.Indexes().DropOne(ctx, indexName)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return ErrCatalogNotFound
			}
		}

		err = errors.Join(errors.New("failed to delete catalog index"), err)
		r.logger.Error(err.Error())
		return err
	}

	r.logger.Info("catalog index deleted", slog.String("indexName", indexName), slog.String("catalogName", catalogName), slog.String("namespace", namespace))
	return nil
}
