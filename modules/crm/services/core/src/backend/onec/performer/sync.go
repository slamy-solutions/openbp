package performer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

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

func SyncWithOneCServer(ctx context.Context, logger *slog.Logger, namespace string, systemStub *system.SystemStub, nativeStub *native.NativeStub, c *connector.OneCConnector) error {
	performers, statusCode, err := connector.MakeRequest[[]performer](
		ctx,
		c,
		"GET",
		fmt.Sprintf("%s/performer/%s", c.ServerURL, c.ServerToken),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to get performers from backend"), err)
		logger.Error(err.Error())
		return err
	}

	if statusCode != 200 {
		err := errors.New(fmt.Sprintf("failed to get performers from backend. Bad status code: %d", statusCode))
		logger.Error(err.Error(), slog.Int("status_code", statusCode))
		return err
	}

	logger.Info("got performers from backend", slog.Int("count", len(*performers)))

	for _, p := range *performers {
		pLogger := logger.With(slog.Group("performer", slog.String("id", p.UUID), slog.String("name", p.Name), slog.String("login", p.Login), slog.String("email", p.Email), slog.String("tel1", p.Tel1), slog.String("tel2", p.Tel2)))

		collection := getPerformerCollection(systemStub, namespace)

		// Check if already exist
		exist := true
		var performerStoredInMongo performerInMongo
		err := collection.FindOne(ctx, bson.M{"uuid": p.UUID}).Decode(&performerStoredInMongo)
		if err != nil {
			if err != mongo.ErrNoDocuments {
				err := errors.Join(errors.New("failed to check performer existance"), err)
				pLogger.Error(err.Error())
				return err
			}

			exist = false
		}

		if exist {
			pLogger = pLogger.With(slog.String("performer_user_uuid", performerStoredInMongo.UserUUID))
			userResponse, err := nativeStub.Services.IAM.Actor.User.Get(ctx, &user.GetRequest{
				Namespace: namespace,
				Uuid:      performerStoredInMongo.UserUUID,
				UseCache:  true,
			})
			if err != nil {
				if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
					pLogger.Error(err.Error())
					return models.ErrPerformerUserNotFound
				}

				err := errors.Join(errors.New("failed to get user"), err)
				pLogger.Error(err.Error())
				return err
			}

			if p.Name != userResponse.User.FullName || p.Email != userResponse.User.Email || p.Login != userResponse.User.Login {
				_, err := nativeStub.Services.IAM.Actor.User.Update(ctx, &user.UpdateRequest{
					Namespace: namespace,
					Uuid:      performerStoredInMongo.UserUUID,
					Login:     p.Login,
					FullName:  p.Name,
					Avatar:    userResponse.User.Avatar, //TODO: change avatar
					Email:     p.Email,
				})
				if err != nil {
					err := errors.Join(errors.New("failed to update user"), err)
					pLogger.Error(err.Error())
					return err
				}
				pLogger.Info("user updated")
			}
		} else {
			userCreateResponse, err := nativeStub.Services.IAM.Actor.User.Create(ctx, &user.CreateRequest{
				Namespace: namespace,
				Login:     p.Login,
				FullName:  p.Name,
				Avatar:    "",
				Email:     p.Email,
			})
			if err != nil {
				err := errors.Join(errors.New("failed to create user"), err)
				pLogger.Error(err.Error())
				return err
			}
			pLogger.Info("user created", slog.String("user_uuid", userCreateResponse.User.Uuid))

			performerStoredInMongo = performerInMongo{
				UUID:           p.UUID,
				UserUUID:       userCreateResponse.User.Uuid,
				DepartmentUUID: p.DepartmentUUID,
			}

			_, err = collection.InsertOne(ctx, performerStoredInMongo)
			if err != nil {
				_, userRemoveErr := nativeStub.Services.IAM.Actor.User.Delete(ctx, &user.DeleteRequest{Namespace: namespace, Uuid: userCreateResponse.User.Uuid})
				if userRemoveErr != nil {
					pLogger.Warn("failed to remove user", slog.String("user_uuid", userCreateResponse.User.Uuid))
				}

				err := errors.Join(errors.New("failed to insert performer"), err)
				pLogger.Error(err.Error())
				return err
			}
			pLogger.Info("performer inserted to the database")
		}
	}

	logger.Info("performers synced successfully")

	return nil
}
