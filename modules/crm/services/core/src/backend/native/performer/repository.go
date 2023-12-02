package performer

import (
	"context"
	"errors"
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PerformerRepository struct {
	logger     *slog.Logger
	namespace  string
	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

func NewPerformerRepository(logger *slog.Logger, namespace string, systemStub *system.SystemStub, nativeStub *native.NativeStub) PerformerRepository {
	return PerformerRepository{
		logger:     logger,
		namespace:  namespace,
		systemStub: systemStub,
		nativeStub: nativeStub,
	}
}

func (r *PerformerRepository) Create(ctx context.Context, departmentUUID string, userUUID string) (*models.Performer, error) {
	departmentUUIDObject, err := primitive.ObjectIDFromHex(departmentUUID)
	if err != nil {
		return nil, models.ErrPerformerBadDepartmentUUID
	}

	userUUIDObject, err := primitive.ObjectIDFromHex(userUUID)
	if err != nil {
		return nil, models.ErrPerformerBadUserUUID
	}

	userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
		Namespace: r.namespace,
		Uuid:      userUUID,
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, models.ErrPerformerUserNotFound
		}

		err = errors.Join(errors.New("failed to get user"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	performer := PerformerInMongo{
		DepartmentUUID: departmentUUIDObject,
		UserUUID:       userUUIDObject,
	}
	insertResult, err := GetPerformerCollection(r.systemStub, r.namespace).InsertOne(ctx, performer)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, models.ErrPerformerAlreadyExists
		}

		err = errors.Join(errors.New("failed to insert performer"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Performer{
		Namespace:      r.namespace,
		UUID:           insertResult.InsertedID.(primitive.ObjectID).Hex(),
		DepartmentUUID: departmentUUID,
		UserUUID:       userUUID,
		Name:           userResponse.User.FullName,
		AvatarURL:      userResponse.User.Avatar,
	}, nil
}
func (r *PerformerRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Performer, error) {
	performerUUIDObject, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrPerformerNotFound
	}

	performer := PerformerInMongo{}
	err = GetPerformerCollection(r.systemStub, r.namespace).FindOne(ctx, bson.M{"_id": performerUUIDObject}).Decode(&performer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrPerformerNotFound
		}

		err = errors.Join(errors.New("failed to get performer"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
		Namespace: r.namespace,
		Uuid:      performer.UserUUID.Hex(),
		UseCache:  useCache,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, models.ErrPerformerUserNotFound
		}

		err = errors.Join(errors.New("failed to get user"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Performer{
		Namespace:      r.namespace,
		UUID:           performer.UUID.Hex(),
		DepartmentUUID: performer.DepartmentUUID.Hex(),
		UserUUID:       performer.UserUUID.Hex(),
		Name:           userResponse.User.FullName,
		AvatarURL:      userResponse.User.Avatar,
	}, nil

}
func (r *PerformerRepository) GetAll(ctx context.Context, useCache bool) ([]models.Performer, error) {
	performers := []PerformerInMongo{}
	cursor, err := GetPerformerCollection(r.systemStub, r.namespace).Find(ctx, bson.M{})
	if err != nil {
		err = errors.Join(errors.New("failed to get performers"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	err = cursor.All(ctx, &performers)
	if err != nil {
		err = errors.Join(errors.New("failed to get performers"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	// TODO: move to goroutine

	performersResponse := make([]models.Performer, len(performers))
	for i, performer := range performers {
		userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
			Namespace: r.namespace,
			Uuid:      performer.UserUUID.Hex(),
			UseCache:  useCache,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				return nil, models.ErrPerformerUserNotFound
			}

			err = errors.Join(errors.New("failed to get user"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		performersResponse[i] = models.Performer{
			Namespace:      r.namespace,
			UUID:           performer.UUID.Hex(),
			DepartmentUUID: performer.DepartmentUUID.Hex(),
			UserUUID:       performer.UserUUID.Hex(),
			Name:           userResponse.User.FullName,
			AvatarURL:      userResponse.User.Avatar,
		}
	}

	return performersResponse, nil
}
func (r *PerformerRepository) Update(ctx context.Context, uuid string, departmentUUID string) (*models.Performer, error) {
	performerUUIDObject, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrPerformerNotFound
	}

	departmentUUIDObject, err := primitive.ObjectIDFromHex(departmentUUID)
	if err != nil {
		return nil, models.ErrPerformerBadDepartmentUUID
	}

	performer := PerformerInMongo{}
	err = GetPerformerCollection(r.systemStub, r.namespace).FindOneAndUpdate(ctx, bson.M{"_id": performerUUIDObject}, bson.M{"$set": bson.M{"departmentUUID": departmentUUIDObject}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&performer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrPerformerNotFound
		}

		err = errors.Join(errors.New("failed to update performer"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
		Namespace: r.namespace,
		Uuid:      performer.UserUUID.Hex(),
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, models.ErrPerformerUserNotFound
		}

		err = errors.Join(errors.New("failed to get user"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Performer{
		Namespace:      r.namespace,
		UUID:           performer.UUID.Hex(),
		DepartmentUUID: performer.DepartmentUUID.Hex(),
		UserUUID:       performer.UserUUID.Hex(),
		Name:           userResponse.User.FullName,
		AvatarURL:      userResponse.User.Avatar,
	}, nil

}
func (r *PerformerRepository) Delete(ctx context.Context, uuid string) (*models.Performer, error) {
	performerUUIDObject, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrPerformerNotFound
	}

	performer := PerformerInMongo{}
	err = GetPerformerCollection(r.systemStub, r.namespace).FindOneAndDelete(ctx, bson.M{"_id": performerUUIDObject}).Decode(&performer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrPerformerNotFound
		}

		err = errors.Join(errors.New("failed to delete performer"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
		Namespace: r.namespace,
		Uuid:      performer.UserUUID.Hex(),
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, models.ErrPerformerUserNotFound
		}

		err = errors.Join(errors.New("failed to get user"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Performer{
		Namespace:      r.namespace,
		UUID:           performer.UUID.Hex(),
		DepartmentUUID: performer.DepartmentUUID.Hex(),
		UserUUID:       performer.UserUUID.Hex(),
		Name:           userResponse.User.FullName,
		AvatarURL:      userResponse.User.Avatar,
	}, nil
}
