package project

import (
	"context"
	"errors"
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectRepository struct {
	logger     *slog.Logger
	namespace  string
	systemStub *system.SystemStub
}

func NewProjectRepository(logger *slog.Logger, namespace string, systemStub *system.SystemStub) ProjectRepository {
	return ProjectRepository{
		logger:     logger,
		namespace:  namespace,
		systemStub: systemStub,
	}
}

func (r *ProjectRepository) Create(ctx context.Context, name string, clientUUID string, contactUUID string, departmentUUID string) (*models.Project, error) {
	clientUUIDObject, err := primitive.ObjectIDFromHex(clientUUID)
	if err != nil {
		return nil, models.ErrProjectBadClientUUID
	}

	contactUUIDObject, err := primitive.ObjectIDFromHex(contactUUID)
	if err != nil {
		return nil, models.ErrProjectBadContactUUID
	}

	departmentUUIDObject, err := primitive.ObjectIDFromHex(departmentUUID)
	if err != nil {
		return nil, models.ErrProjectBadDepartmentUUID
	}

	collection := GetProjectCollection(r.systemStub, r.namespace)

	project := ProjectInMongo{
		Name:           name,
		ClientUUID:     clientUUIDObject,
		ContactUUID:    contactUUIDObject,
		DepartmentUUID: departmentUUIDObject,
		NotRelevant:    false,
	}
	insertResponse, err := collection.InsertOne(ctx, project)
	if err != nil {
		err = errors.Join(errors.New("failed to insert project"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	projectUUID := insertResponse.InsertedID.(primitive.ObjectID).Hex()
	return &models.Project{
		Namespace:      r.namespace,
		UUID:           projectUUID,
		Name:           project.Name,
		ClientUUID:     project.ClientUUID.Hex(),
		ContactUUID:    project.ContactUUID.Hex(),
		DepartmentUUID: project.DepartmentUUID.Hex(),
		NotRelevant:    project.NotRelevant,
	}, nil
}
func (r *ProjectRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Project, error) {
	uuidObject, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrProjectNotFound
	}

	collection := GetProjectCollection(r.systemStub, r.namespace)
	var project ProjectInMongo
	err = collection.FindOne(ctx, bson.M{"_id": uuidObject}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrProjectNotFound
		}

		err = errors.Join(errors.New("failed to get project from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Project{
		Namespace:      r.namespace,
		UUID:           project.UUID.Hex(),
		Name:           project.Name,
		ClientUUID:     project.ClientUUID.Hex(),
		ContactUUID:    project.ContactUUID.Hex(),
		DepartmentUUID: project.DepartmentUUID.Hex(),
		NotRelevant:    project.NotRelevant,
	}, nil
}
func (r *ProjectRepository) GetAll(ctx context.Context, useCache bool, clientUUID string, departmentUUID string) ([]models.Project, error) {
	collection := GetProjectCollection(r.systemStub, r.namespace)

	cursor, err := collection.Find(ctx, bson.M{"clientUUID": clientUUID, "departmentUUID": departmentUUID, "notRelevant": false})
	if err != nil {
		err = errors.Join(errors.New("failed to get projects from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	projects := []models.Project{}
	for cursor.Next(ctx) {
		var project ProjectInMongo
		err = cursor.Decode(&project)
		if err != nil {
			err = errors.Join(errors.New("failed to decode project"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		projects = append(projects, models.Project{
			Namespace:      r.namespace,
			UUID:           project.UUID.Hex(),
			Name:           project.Name,
			ClientUUID:     project.ClientUUID.Hex(),
			ContactUUID:    project.ContactUUID.Hex(),
			DepartmentUUID: project.DepartmentUUID.Hex(),
			NotRelevant:    project.NotRelevant,
		})
	}

	return projects, nil
}
func (r *ProjectRepository) Update(ctx context.Context, uuid string, name string, clientUUID string, contactUUID string, departmentUUID string, notRelevant bool) (*models.Project, error) {
	uuidObject, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrProjectNotFound
	}

	clientUUIDObject, err := primitive.ObjectIDFromHex(clientUUID)
	if err != nil {
		return nil, models.ErrProjectBadClientUUID
	}

	contactUUIDObject, err := primitive.ObjectIDFromHex(contactUUID)
	if err != nil {
		return nil, models.ErrProjectBadContactUUID
	}

	departmentUUIDObject, err := primitive.ObjectIDFromHex(departmentUUID)
	if err != nil {
		return nil, models.ErrProjectBadDepartmentUUID
	}

	collection := GetProjectCollection(r.systemStub, r.namespace)
	updateResponse, err := collection.UpdateOne(ctx, bson.M{"_id": uuidObject}, bson.M{
		"$set": bson.M{
			"name":           name,
			"clientUUID":     clientUUIDObject,
			"contactUUID":    contactUUIDObject,
			"departmentUUID": departmentUUIDObject,
			"notRelevant":    notRelevant,
		},
	})
	if err != nil {
		err = errors.Join(errors.New("failed to update project"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	if updateResponse.MatchedCount == 0 {
		return nil, models.ErrProjectNotFound
	}

	return &models.Project{
		Namespace:      r.namespace,
		UUID:           uuid,
		Name:           name,
		ClientUUID:     clientUUID,
		ContactUUID:    contactUUID,
		DepartmentUUID: departmentUUID,
		NotRelevant:    notRelevant,
	}, nil

}
func (r *ProjectRepository) Delete(ctx context.Context, uuid string) (*models.Project, error) {
	uuidObject, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrProjectNotFound
	}

	collection := GetProjectCollection(r.systemStub, r.namespace)
	var project ProjectInMongo
	err = collection.FindOneAndDelete(ctx, bson.M{"_id": uuidObject}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrProjectNotFound
		}

		err = errors.Join(errors.New("failed to delete project"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Project{
		Namespace:      r.namespace,
		UUID:           project.UUID.Hex(),
		Name:           project.Name,
		ClientUUID:     project.ClientUUID.Hex(),
		ContactUUID:    project.ContactUUID.Hex(),
		DepartmentUUID: project.DepartmentUUID.Hex(),
		NotRelevant:    project.NotRelevant,
	}, nil

}
