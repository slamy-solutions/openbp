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
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (s *ServersServer) List(in *balena.ListServersRequest, out balena.BalenaServersService_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (s *ServersServer) Count(ctx context.Context, in *balena.CountServersRequest) (*balena.CountServersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Count not implemented")
}
func (s *ServersServer) ListInNamespace(in *balena.ListServersInNamespaceRequest, out balena.BalenaServersService_ListInNamespaceServer) error {
	return status.Errorf(codes.Unimplemented, "method ListInNamespace not implemented")
}
func (s *ServersServer) CountInNamespace(ctx context.Context, in *balena.CountServersInNamespaceRequest) (*balena.CountServersInNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CountInNamespace not implemented")
}
func (s *ServersServer) SetEnabled(ctx context.Context, in *balena.SetServerEnabledRequest) (*balena.SetServerEnabledResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetEnabled not implemented")
}
func (s *ServersServer) Update(ctx context.Context, in *balena.UpdateServerRequest) (*balena.UpdateServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (s *ServersServer) Delete(ctx context.Context, in *balena.DeleteServerRequest) (*balena.DeleteServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
