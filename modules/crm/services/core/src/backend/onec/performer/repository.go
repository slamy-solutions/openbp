package performer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/connector"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type performer struct {
	UUID           string `json:"id"`
	Name           string `json:"name"`
	DepartmentUUID string `json:"departmentId"`
	Avatar         struct {
		Name string `json:"name"`
		UUID string `json:"id"`
	} `json:"avatar"`
	Login string `json:"login"`
	Email string `json:"email"`
	Tel1  string `json:"tel1"`
	Tel2  string `json:"tel2"`
}

type PerformerRepository struct {
	logger     *slog.Logger
	namespace  string
	systemStub *system.SystemStub
	nativeStub *native.NativeStub
	connector  *connector.OneCConnector
}

func NewPerformerRepository(logger *slog.Logger, namespace string, systemStub *system.SystemStub, nativeStub *native.NativeStub, connector *connector.OneCConnector) PerformerRepository {
	return PerformerRepository{
		logger:     logger,
		namespace:  namespace,
		systemStub: systemStub,
		nativeStub: nativeStub,
		connector:  connector,
	}
}

func (r *PerformerRepository) Create(ctx context.Context, departmentUUID string, userUUID string) (*models.Performer, error) {
	userData, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
		Namespace: r.namespace,
		Uuid:      userUUID,
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			err := errors.Join(models.ErrPerformerUserNotFound, err)
			r.logger.Error(err.Error())
			return nil, err
		}

		err := errors.Join(errors.New("failed to get user data"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	var requestData struct {
		Name           string `json:"name"`
		DepartmentUUID string `json:"departmentId"`
		Login          string `json:"login"`
	} = struct {
		Name           string "json:\"name\""
		DepartmentUUID string "json:\"departmentId\""
		Login          string "json:\"login\""
	}{
		Name:           userData.User.FullName,
		Login:          userData.User.Login,
		DepartmentUUID: departmentUUID,
	}

	createdPerformer, statusCode, err := connector.MakeRequest[performer](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/performer/%s", r.connector.ServerURL, r.connector.ServerToken),
		requestData,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create performer"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to create performer. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	// Add new performer to the database
	collection := getPerformerCollection(r.systemStub, r.namespace)
	_, err = collection.InsertOne(ctx, performerInMongo{
		UUID:           createdPerformer.UUID,
		DepartmentUUID: departmentUUID,
		UserUUID:       userUUID,
	})
	if err != nil {
		err := errors.Join(errors.New("failed to insert new performer to the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Performer{
		Namespace:       r.namespace,
		UUID:            createdPerformer.UUID,
		DeparatmentUUID: departmentUUID,
		UserUUID:        userUUID,
		Name:            userData.User.FullName,
		AvatarURL:       userData.User.Avatar,
	}, nil
}
func (r *PerformerRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Performer, error) {
	collection := getPerformerCollection(r.systemStub, r.namespace)
	var performer performerInMongo
	err := collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&performer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrPerformerNotFound
		}

		err := errors.Join(errors.New("failed to get performer from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
		Namespace: r.namespace,
		Uuid:      performer.UserUUID,
		UseCache:  useCache,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			err := errors.Join(models.ErrPerformerUserNotFound, err)
			r.logger.Error(err.Error())
			return nil, err
		}

		err := errors.Join(errors.New("failed to get user data"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Performer{
		Namespace:       r.namespace,
		UUID:            performer.UUID,
		DeparatmentUUID: performer.DepartmentUUID,
		UserUUID:        performer.UserUUID,
		Name:            userResponse.User.FullName,
		AvatarURL:       userResponse.User.Avatar,
	}, nil
}
func (r *PerformerRepository) GetAll(ctx context.Context, useCache bool) ([]models.Performer, error) {
	collection := getPerformerCollection(r.systemStub, r.namespace)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		err := errors.Join(errors.New("failed to get all performers from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(context.Background())

	var performers []models.Performer
	for cursor.Next(ctx) {
		var performer performerInMongo
		err := cursor.Decode(&performer)
		if err != nil {
			err := errors.Join(errors.New("failed to decode performer from the database"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
			Namespace: r.namespace,
			Uuid:      performer.UserUUID,
			UseCache:  useCache,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				err := errors.Join(models.ErrPerformerUserNotFound, err)
				r.logger.Error(err.Error())
				return nil, err
			}

			err := errors.Join(errors.New("failed to get user data"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		performers = append(performers, models.Performer{
			Namespace:       r.namespace,
			UUID:            performer.UUID,
			DeparatmentUUID: performer.DepartmentUUID,
			UserUUID:        performer.UserUUID,
			Name:            userResponse.User.FullName,
			AvatarURL:       userResponse.User.Avatar,
		})
	}

	return performers, nil
}
func (r *PerformerRepository) Update(ctx context.Context, uuid string, departmentUUID string) (*models.Performer, error) {
	collection := getPerformerCollection(r.systemStub, r.namespace)
	var existingPerformer performerInMongo
	err := collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&existingPerformer)
	if err != nil {
		err := errors.Join(errors.New("failed to get performer from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
		Namespace: r.namespace,
		Uuid:      existingPerformer.UserUUID,
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			err := errors.Join(models.ErrPerformerUserNotFound, err)
			r.logger.Error(err.Error())
			return nil, err
		}

		err := errors.Join(errors.New("failed to get user data"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	p := performer{
		UUID:           existingPerformer.UUID,
		Name:           userResponse.User.FullName,
		DepartmentUUID: departmentUUID,
		Login:          userResponse.User.Login,
		Email:          userResponse.User.Email,
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/performer/%s", r.connector.ServerURL, r.connector.ServerToken),
		p,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update performer"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to update performer. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	return &models.Performer{
		Namespace:       r.namespace,
		UUID:            existingPerformer.UUID,
		DeparatmentUUID: departmentUUID,
		UserUUID:        existingPerformer.UserUUID,
		Name:            userResponse.User.FullName,
		AvatarURL:       userResponse.User.Avatar,
	}, nil
}
func (r *PerformerRepository) Delete(ctx context.Context, uuid string) (*models.Performer, error) {
	collection := getPerformerCollection(r.systemStub, r.namespace)
	var existingPerformer performerInMongo
	err := collection.FindOneAndDelete(ctx, bson.M{"uuid": uuid}).Decode(&existingPerformer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrPerformerNotFound
		}

		err := errors.Join(errors.New("failed to get performer from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	userResponse, err := r.nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
		Namespace: r.namespace,
		Uuid:      existingPerformer.UserUUID,
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			err := errors.Join(models.ErrPerformerUserNotFound, err)
			r.logger.Error(err.Error())
			return nil, err
		}

		err := errors.Join(errors.New("failed to get user data"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"DELETE",
		fmt.Sprintf("%s/performer/%s/%s", r.connector.ServerURL, r.connector.ServerToken, uuid),
		bson.M{"id": uuid},
	)
	if err != nil {
		err := errors.Join(errors.New("failed to delete performer"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to delete performer. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	return &models.Performer{
		Namespace:       r.namespace,
		UUID:            existingPerformer.UUID,
		DeparatmentUUID: existingPerformer.DepartmentUUID,
		UserUUID:        existingPerformer.UserUUID,
		Name:            userResponse.User.FullName,
		AvatarURL:       userResponse.User.Avatar,
	}, nil
}
