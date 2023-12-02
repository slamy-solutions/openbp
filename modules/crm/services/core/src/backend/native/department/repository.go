package department

import (
	"context"
	"errors"
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DepartmentRepository struct {
	logger     *slog.Logger
	namespace  string
	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

func NewDepartmentRepository(logger *slog.Logger, namespace string, systemStub *system.SystemStub, nativeStub *native.NativeStub) DepartmentRepository {
	return DepartmentRepository{
		logger:     logger,
		namespace:  namespace,
		systemStub: systemStub,
		nativeStub: nativeStub,
	}
}

func (r *DepartmentRepository) Create(ctx context.Context, name string) (*models.Department, error) {
	collection := GetDepartmentCollection(r.systemStub, r.namespace)
	department := DepartmentInMongo{
		Name: name,
	}
	insertResult, err := collection.InsertOne(ctx, department)
	if err != nil {
		err = errors.Join(errors.New("failed to insert department to the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	r.logger.Info("Created department", slog.Group("department", "name", name, "uuid", insertResult.InsertedID))
	return &models.Department{
		Namespace: r.namespace,
		UUID:      insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Name:      name,
	}, nil
}
func (r *DepartmentRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Department, error) {
	collection := GetDepartmentCollection(r.systemStub, r.namespace)
	departmentUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrDepartmentNotFound
	}

	var department DepartmentInMongo
	err = collection.FindOne(ctx, bson.M{"_id": departmentUUID}).Decode(&department)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrDepartmentNotFound
		}

		err = errors.Join(errors.New("failed to get department from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Department{
		Namespace: r.namespace,
		UUID:      department.UUID.Hex(),
		Name:      department.Name,
	}, nil
}
func (r *DepartmentRepository) GetAll(ctx context.Context, useCache bool) ([]models.Department, error) {
	collection := GetDepartmentCollection(r.systemStub, r.namespace)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		err = errors.Join(errors.New("failed to get departments from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	var departments []models.Department
	for cursor.Next(ctx) {
		var department DepartmentInMongo
		err = cursor.Decode(&department)
		if err != nil {
			err = errors.Join(errors.New("failed to decode department from the database"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		departments = append(departments, models.Department{
			Namespace: r.namespace,
			UUID:      department.UUID.Hex(),
			Name:      department.Name,
		})
	}

	return departments, nil
}
func (r *DepartmentRepository) Update(ctx context.Context, uuid string, name string) (*models.Department, error) {
	collection := GetDepartmentCollection(r.systemStub, r.namespace)
	departmentUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrDepartmentNotFound
	}

	updateResult, err := collection.UpdateOne(ctx, bson.M{"_id": departmentUUID}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		err = errors.Join(errors.New("failed to update department in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		return nil, models.ErrDepartmentNotFound
	}

	r.logger.Info("Updated department", slog.Group("department", "name", name, "uuid", uuid))
	return &models.Department{
		Namespace: r.namespace,
		UUID:      uuid,
		Name:      name,
	}, nil
}
func (r *DepartmentRepository) Delete(ctx context.Context, uuid string) (*models.Department, error) {
	collection := GetDepartmentCollection(r.systemStub, r.namespace)
	departmentUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrDepartmentNotFound
	}

	var department DepartmentInMongo
	err = collection.FindOneAndDelete(ctx, bson.M{"_id": departmentUUID}).Decode(&department)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrDepartmentNotFound
		}

		err = errors.Join(errors.New("failed to delete department from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	r.logger.Info("Deleted department", slog.Group("department", "uuid", uuid, "name", department.Name))
	return &models.Department{
		Namespace: r.namespace,
		UUID:      uuid,
		Name:      department.Name,
	}, nil
}
