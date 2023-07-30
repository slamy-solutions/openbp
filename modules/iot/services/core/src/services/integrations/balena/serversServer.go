package balena

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	balena "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServersServer struct {
	balena.UnimplementedBalenaServersServiceServer

	systemStub *system.SystemStub
	logger     *logrus.Entry
}

func NewServersServer(logger *logrus.Entry, systemStub *system.SystemStub) *ServersServer {
	return &ServersServer{
		systemStub: systemStub,
		logger:     logger,
	}
}

func (s *ServersServer) encryptAuthToken(ctx context.Context, plainToken string) ([]byte, error) {
	encryptStream, err := s.systemStub.Vault.EncryptStream(ctx)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
			return nil, status.Error(codes.FailedPrecondition, "auth token encryption precondition failed: "+err.Error())
		}

		err = errors.Join(errors.New("error while openning ecryption stream with system_vault service to encrypt auth token"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	dataToEncrypt := []byte(plainToken)
	if err = encryptStream.Send(&vault.EncryptStreamRequest{PlainData: dataToEncrypt}); err != nil {
		err = errors.Join(errors.New("error while sending auth token for ecryption"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err = encryptStream.CloseSend(); err != nil {
		err = errors.Join(errors.New("error while sending auth token for ecryption. Error closing send stream"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := make([]byte, 0, len(dataToEncrypt))
	for {
		chunk, err := encryptStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Join(errors.New("error while receiving encrypted token from system_vault"), err)
			s.logger.Error(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}

		response = append(response, chunk.EncryptedData...)
	}

	return response, nil
}

func (s *ServersServer) Create(ctx context.Context, in *balena.CreateServerRequest) (*balena.CreateServerResponse, error) {
	authToken, err := s.encryptAuthToken(ctx, in.AuthToken)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now().UTC()
	serverInMongo := BalenaServerInMongo{
		Namespace:   in.Namespace,
		Name:        in.Name,
		Description: in.Description,
		BaseURL:     in.BaseURL,
		AuthToken:   authToken,
		Enabled:     false,
		Created:     currentTime,
		Updated:     currentTime,
		Version:     0,
	}
	collection := getBalenaServerCollection(s.systemStub)
	updateResult, err := collection.UpdateOne(
		ctx,
		bson.M{"namespace": in.Namespace, "name": in.Name},
		bson.M{"$setOnInsert": serverInMongo},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		err = errors.Join(errors.New("error while inserting new balena server to the database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if updateResult.UpsertedCount == 0 {
		return nil, status.Error(codes.AlreadyExists, "server with same name already exist in provided namespace")
	}
	serverInMongo.UUID = updateResult.UpsertedID.(primitive.ObjectID)

	return &balena.CreateServerResponse{}, status.Error(codes.OK, "")
}
func (s *ServersServer) Get(ctx context.Context, in *balena.GetServerRequest) (*balena.GetServerResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "server not found. UUID has bad format")
	}

	collection := getBalenaServerCollection(s.systemStub)
	var server BalenaServerInMongo
	err = collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&server)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		err = errors.Join(errors.New("error while getting balena server information from the database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.GetServerResponse{Server: server.ToGRPCServer()}, status.Error(codes.OK, "")
}
func (s *ServersServer) List(in *balena.ListServersRequest, out balena.BalenaServersService_ListServer) error {
	ctx := out.Context()
	collection := getBalenaServerCollection(s.systemStub)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		err = errors.Join(errors.New("error while getting balena servers list from the database: failed to open database cursor"), err)
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	defer cur.Close(context.Background())

	for cur.Next(ctx) {
		var server BalenaServerInMongo
		err = cur.Decode(&server)
		if err != nil {
			err = errors.Join(errors.New("error while getting balena servers list from the database: error while decoding server information"), err)
			s.logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		if err = out.Send(&balena.ListServersResponse{Server: server.ToGRPCServer()}); err != nil {
			err = errors.Join(errors.New("error while sending server information to the output stream"), err)
			s.logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}

	if err = cur.Err(); err != nil {
		err = errors.Join(errors.New("error while getting balena servers list from the database"), err)
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	return status.Error(codes.OK, "")
}
func (s *ServersServer) Count(ctx context.Context, in *balena.CountServersRequest) (*balena.CountServersResponse, error) {
	collection := getBalenaServerCollection(s.systemStub)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		err = errors.Join(errors.New("error while counting balena servers in database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.CountServersResponse{TotalCount: uint64(count)}, status.Error(codes.OK, "")
}
func (s *ServersServer) ListInNamespace(in *balena.ListServersInNamespaceRequest, out balena.BalenaServersService_ListInNamespaceServer) error {
	ctx := out.Context()
	collection := getBalenaServerCollection(s.systemStub)
	cur, err := collection.Find(ctx, bson.M{"namespace": in.Namespace})
	if err != nil {
		err = errors.Join(errors.New("error while getting balena servers list from the database: failed to open database cursor"), err)
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	defer cur.Close(context.Background())

	for cur.Next(ctx) {
		var server BalenaServerInMongo
		err = cur.Decode(&server)
		if err != nil {
			err = errors.Join(errors.New("error while getting balena servers list from the database: error while decoding server information"), err)
			s.logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		if err = out.Send(&balena.ListServersInNamespaceResponse{Server: server.ToGRPCServer()}); err != nil {
			err = errors.Join(errors.New("error while sending server information to the output stream"), err)
			s.logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}
	}

	if err = cur.Err(); err != nil {
		err = errors.Join(errors.New("error while getting balena servers list from the database"), err)
		s.logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}

	return status.Error(codes.OK, "")
}
func (s *ServersServer) CountInNamespace(ctx context.Context, in *balena.CountServersInNamespaceRequest) (*balena.CountServersInNamespaceResponse, error) {
	collection := getBalenaServerCollection(s.systemStub)
	count, err := collection.CountDocuments(ctx, bson.M{"namespace": in.Namespace})
	if err != nil {
		err = errors.Join(errors.New("error while counting balena servers in database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.CountServersInNamespaceResponse{TotalCount: uint64(count)}, status.Error(codes.OK, "")
}
func (s *ServersServer) SetEnabled(ctx context.Context, in *balena.SetServerEnabledRequest) (*balena.SetServerEnabledResponse, error) {
	serverUUID, err := primitive.ObjectIDFromHex(in.ServerUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "server not found. UUID has bad format")
	}

	collection := getBalenaServerCollection(s.systemStub)
	var updatedServer BalenaServerInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": serverUUID},
		bson.M{
			"$set":         bson.M{"enabled": in.Enabled},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedServer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		err = errors.Join(errors.New("error while setting balena server enabled in database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.SetServerEnabledResponse{
		Server: updatedServer.ToGRPCServer(),
	}, status.Error(codes.OK, "")
}
func (s *ServersServer) Update(ctx context.Context, in *balena.UpdateServerRequest) (*balena.UpdateServerResponse, error) {
	serverUUID, err := primitive.ObjectIDFromHex(in.ServerUUID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "server not found. UUID has bad format")
	}

	collection := getBalenaServerCollection(s.systemStub)
	updateQuery := bson.M{
		"$set": bson.M{
			"description": in.NewDescription,
		},
		"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
		"$inc":         bson.M{"version": 1},
	}
	if in.NewBaseURL != "" {
		setData := updateQuery["$set"].(bson.M)
		setData["baseURL"] = in.NewBaseURL
		setData["enabled"] = false
	}
	if in.NewAuthToken != "" {
		encryptedToken, err := s.encryptAuthToken(ctx, in.NewAuthToken)
		if err != nil {
			return nil, err
		}
		setData := updateQuery["$set"].(bson.M)
		setData["authToken"] = encryptedToken
	}

	var updatedServer BalenaServerInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": serverUUID},
		updateQuery,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedServer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		err = errors.Join(errors.New("error while updating balena server in database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.UpdateServerResponse{
		Server: updatedServer.ToGRPCServer(),
	}, status.Error(codes.OK, "")
}
func (s *ServersServer) Delete(ctx context.Context, in *balena.DeleteServerRequest) (*balena.DeleteServerResponse, error) {
	serverUUID, err := primitive.ObjectIDFromHex(in.ServerUUID)
	if err != nil {
		return &balena.DeleteServerResponse{Existed: false}, status.Error(codes.OK, "server not found. UUID has bad format")
	}

	collection := getBalenaServerCollection(s.systemStub)
	response, err := collection.DeleteOne(ctx, bson.M{"_id": serverUUID})
	if err != nil {
		err = errors.Join(errors.New("error while deleting balena server from the database"), err)
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &balena.DeleteServerResponse{
		Existed: response.DeletedCount != 0,
	}, status.Error(codes.OK, "")
}
