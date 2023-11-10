package services

import (
	"context"
	"errors"
	"log/slog"

	client "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/client"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientService struct {
	client.UnimplementedClientServiceServer

	backend backend.BackendFactory
	logger  *slog.Logger
}

func NewClientServer(backend backend.BackendFactory, logger *slog.Logger) *ClientService {
	return &ClientService{
		backend: backend,
		logger:  logger,
	}
}

func (s *ClientService) GetAll(ctx context.Context, in *client.GetAllRequest) (*client.GetAllResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	clients, err := bkd.ClientRepository().GetAll(ctx, in.UseCache)
	if err != nil {
		err := errors.Join(errors.New("failed to get clients"), err)
		s.logger.With(slog.String("route", "GetAll")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get clients: %s", err.Error())
	}

	var responseClients []*client.Client = make([]*client.Client, len(clients))
	for i, client := range clients {
		responseClients[i] = client.ToGRPC()
	}

	return &client.GetAllResponse{
		Clients: responseClients,
	}, status.Error(codes.OK, "")
}
func (s *ClientService) Get(ctx context.Context, in *client.GetRequest) (*client.GetResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "Get")))
	if err != nil {
		return nil, err
	}

	c, err := bkd.ClientRepository().Get(ctx, in.Uuid, in.UseCache)
	if err != nil {
		err := errors.Join(errors.New("failed to get client"), err)
		s.logger.With(slog.String("route", "Get")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get client: %s", err.Error())
	}

	return &client.GetResponse{
		Client: c.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ClientService) Create(ctx context.Context, in *client.CreateRequest) (*client.CreateResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "Create")))
	if err != nil {
		return nil, err
	}

	c, err := bkd.ClientRepository().Create(ctx, in.Name)
	if err != nil {
		err := errors.Join(errors.New("failed to create client"), err)
		s.logger.With(slog.String("route", "Create")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create client: %s", err.Error())
	}

	return &client.CreateResponse{
		Client: c.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ClientService) Update(ctx context.Context, in *client.UpdateRequest) (*client.UpdateResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "Update")))
	if err != nil {
		return nil, err
	}

	c, err := bkd.ClientRepository().Update(ctx, in.Uuid, in.Name)
	if err != nil {
		err := errors.Join(errors.New("failed to update client"), err)
		s.logger.With(slog.String("route", "Update")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update client: %s", err.Error())
	}

	return &client.UpdateResponse{
		Client: c.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ClientService) Delete(ctx context.Context, in *client.DeleteRequest) (*client.DeleteResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "Delete")))
	if err != nil {
		return nil, err
	}

	_, err = bkd.ClientRepository().Delete(ctx, in.Uuid)
	if err != nil {
		err := errors.Join(errors.New("failed to delete client"), err)
		s.logger.With(slog.String("route", "Delete")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to delete client: %s", err.Error())
	}

	return &client.DeleteResponse{}, status.Error(codes.OK, "")
}

func (s *ClientService) AddContactPerson(ctx context.Context, in *client.AddContactPersonRequest) (*client.AddContactPersonResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "AddContactPerson")))
	if err != nil {
		return nil, err
	}

	contact, err := bkd.ClientRepository().AddContactPerson(ctx, in.ClientUUID, in.Name, in.Email, in.Phone, in.Comment)
	if err != nil {
		err := errors.Join(errors.New("failed to add contact person"), err)
		s.logger.With(slog.String("route", "AddContactPerson")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to add contact person: %s", err.Error())
	}

	return &client.AddContactPersonResponse{
		ContactPerson: contact.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ClientService) UpdateContactPerson(ctx context.Context, in *client.UpdateContactPersonRequest) (*client.UpdateContactPersonResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "UpdateContactPerson")))
	if err != nil {
		return nil, err
	}

	contact, err := bkd.ClientRepository().UpdateContactPerson(ctx, in.ClientUUID, in.ContactPersonUUID, in.Name, in.Email, in.Phone, in.NotRelevant, in.Comment)
	if err != nil {
		if err == models.ErrClientContactPersonNotFound {
			return nil, status.Errorf(codes.NotFound, "client contact person not found")
		}

		err := errors.Join(errors.New("failed to update contact person"), err)
		s.logger.With(slog.String("route", "UpdateContactPerson")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update contact person: %s", err.Error())
	}

	return &client.UpdateContactPersonResponse{
		ContactPerson: contact.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ClientService) DeleteContactPerson(ctx context.Context, in *client.DeleteContactPersonRequest) (*client.DeleteContactPersonResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "DeleteContactPerson")))
	if err != nil {
		return nil, err
	}

	contact, err := bkd.ClientRepository().DeleteContactPerson(ctx, in.ContactPersonUUID)
	if err != nil {
		if err == models.ErrClientContactPersonNotFound {
			return nil, status.Errorf(codes.NotFound, "client contact person not found")
		}

		err := errors.Join(errors.New("failed to delete contact person"), err)
		s.logger.With(slog.String("route", "DeleteContactPerson")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to delete contact person: %s", err.Error())
	}

	return &client.DeleteContactPersonResponse{
		ContactPerson: contact.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ClientService) GetContactPersonsForClient(ctx context.Context, in *client.GetContactPersonsForClientRequest) (*client.GetContactPersonsForClientResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetContactPersonsForClient")))
	if err != nil {
		return nil, err
	}

	contacts, err := bkd.ClientRepository().GetContactPersonsForClient(ctx, in.ClientUUID, in.UseCache)
	if err != nil {
		err := errors.Join(errors.New("failed to get contact persons for client"), err)
		s.logger.With(slog.String("route", "GetContactPersonsForClient")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get contact persons for client: %s", err.Error())
	}

	var responseContacts []*client.ContactPerson = make([]*client.ContactPerson, len(contacts))
	for i, contact := range contacts {
		responseContacts[i] = contact.ToGRPC()
	}

	return &client.GetContactPersonsForClientResponse{
		ContactPersons: responseContacts,
	}, status.Error(codes.OK, "")
}
