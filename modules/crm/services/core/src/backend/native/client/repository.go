package client

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientRepository struct {
	logger     *slog.Logger
	namespace  string
	systemStub *system.SystemStub
}

func NewClientRepository(logger *slog.Logger, namespace string, systemStub *system.SystemStub) ClientRepository {
	return ClientRepository{
		logger:     logger,
		namespace:  namespace,
		systemStub: systemStub,
	}
}

func (r *ClientRepository) Create(ctx context.Context, name string) (*models.Client, error) {
	collection := GetClientCollection(r.systemStub, r.namespace)
	currentTile := time.Now().UTC()
	client := &ClientInMongo{
		Name:           name,
		LastUpdateTime: currentTile,
		CreationTime:   currentTile,
		Version:        1,
	}

	result, err := collection.InsertOne(ctx, client)
	if err != nil {
		err = errors.Join(errors.New("failed to insert new client to the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	clientUUID := result.InsertedID.(primitive.ObjectID).Hex()
	return &models.Client{
		Namespace:      r.namespace,
		UUID:           clientUUID,
		Name:           name,
		LastUpdateTime: currentTile,
		CreationTime:   currentTile,
		Version:        1,
	}, nil
}
func (r *ClientRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Client, error) {
	collection := GetClientCollection(r.systemStub, r.namespace)

	clientUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrClientNotFound
	}

	var client ClientInMongo
	err = collection.FindOne(ctx, primitive.M{"_id": clientUUID}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrClientNotFound
		}

		err = errors.Join(errors.New("failed to get client from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	contacts, err := r.GetContactPersonsForClient(ctx, uuid, useCache)
	if err != nil {
		return nil, err
	}

	return &models.Client{
		Namespace:      r.namespace,
		UUID:           clientUUID.Hex(),
		Name:           client.Name,
		ContactPersons: contacts,
		LastUpdateTime: client.LastUpdateTime,
		CreationTime:   client.CreationTime,
		Version:        client.Version,
	}, nil
}
func (r *ClientRepository) GetAll(ctx context.Context, useCache bool) ([]models.Client, error) {
	collection := GetClientCollection(r.systemStub, r.namespace)

	cursor, err := collection.Find(ctx, primitive.M{})
	if err != nil {
		err = errors.Join(errors.New("failed to get clients from the database. failed to open cursor"), err)
		r.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	// TODO: use aggregation and lookup to get contact persons

	var clients []models.Client
	for cursor.Next(ctx) {
		var client ClientInMongo
		err := cursor.Decode(&client)
		if err != nil {
			err = errors.Join(errors.New("failed to get clients from the database. failed to decode client"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		contacts, err := r.GetContactPersonsForClient(ctx, client.UUID.Hex(), useCache)
		if err != nil {
			return nil, err
		}

		clients = append(clients, models.Client{
			Namespace:      r.namespace,
			UUID:           client.UUID.Hex(),
			Name:           client.Name,
			ContactPersons: contacts,
			LastUpdateTime: client.LastUpdateTime,
			CreationTime:   client.CreationTime,
			Version:        client.Version,
		})
	}

	return clients, nil
}
func (r *ClientRepository) Update(ctx context.Context, uuid string, name string) (*models.Client, error) {
	collection := GetClientCollection(r.systemStub, r.namespace)

	clientUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrClientNotFound
	}

	var client ClientInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": clientUUID}, primitive.M{
		"$set":         bson.M{"name": name},
		"$inc":         bson.M{"version": 1},
		"$currentDate": bson.M{"lastUpdateTime": bson.M{"$type": "timestamp"}},
	}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrClientNotFound
		}

		err = errors.Join(errors.New("failed to update client in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	contacts, err := r.GetContactPersonsForClient(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	return &models.Client{
		Namespace:      r.namespace,
		UUID:           clientUUID.Hex(),
		Name:           name,
		ContactPersons: contacts,
		LastUpdateTime: client.LastUpdateTime,
		CreationTime:   client.CreationTime,
		Version:        client.Version,
	}, nil
}
func (r *ClientRepository) Delete(ctx context.Context, uuid string) (*models.Client, error) {
	collection := GetClientCollection(r.systemStub, r.namespace)

	clientUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return nil, models.ErrClientNotFound
	}

	contacts, err := r.GetContactPersonsForClient(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	var client ClientInMongo
	err = collection.FindOneAndDelete(ctx, bson.M{"_id": clientUUID}).Decode(&client)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrClientNotFound
		}

		err = errors.Join(errors.New("failed to delete client from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return &models.Client{
		Namespace:      r.namespace,
		UUID:           clientUUID.Hex(),
		Name:           client.Name,
		ContactPersons: contacts,
		LastUpdateTime: client.LastUpdateTime,
		CreationTime:   client.CreationTime,
		Version:        client.Version,
	}, nil
}

func (r *ClientRepository) AddContactPerson(ctx context.Context, clientUUID string, name string, email string, phone []string, comment string) (*models.ContactPerson, error) {
	collection := GetClientContactPersonCollection(r.systemStub, r.namespace)

	clientUUIDBson, err := primitive.ObjectIDFromHex(clientUUID)
	if err != nil {
		return nil, models.ErrClientNotFound
	}

	contactPerson := &ContactPersonInMongo{
		ClientUUID:  clientUUIDBson,
		Name:        name,
		Email:       email,
		Phone:       phone,
		NotRelevant: false,
		Comment:     comment,
	}

	result, err := collection.InsertOne(ctx, contactPerson)
	if err != nil {
		err = errors.Join(errors.New("failed to insert new contact person to the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	contactPerson.UUID = result.InsertedID.(primitive.ObjectID)
	return contactPerson.ToBackendModel(r.namespace), nil
}
func (r *ClientRepository) UpdateContactPerson(ctx context.Context, clientUUID string, contactPersonUUID string, name string, email string, phone []string, notRelevant bool, comment string) (*models.ContactPerson, error) {
	collection := GetClientContactPersonCollection(r.systemStub, r.namespace)

	clientUUIDBson, err := primitive.ObjectIDFromHex(clientUUID)
	if err != nil {
		return nil, models.ErrClientNotFound
	}

	contactPersonID, err := primitive.ObjectIDFromHex(contactPersonUUID)
	if err != nil {
		return nil, models.ErrClientContactPersonNotFound
	}

	var contactPerson ContactPersonInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": contactPersonID}, primitive.M{
		"$set": bson.M{"name": name, "email": email, "phone": phone, "notRelevant": notRelevant, "comment": comment, "clientUUID": clientUUIDBson},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&contactPerson)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrClientContactPersonNotFound
		}

		err = errors.Join(errors.New("failed to update contact person in the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return contactPerson.ToBackendModel(r.namespace), nil
}
func (r *ClientRepository) DeleteContactPerson(ctx context.Context, contactPersonUUID string) (*models.ContactPerson, error) {
	collection := GetClientContactPersonCollection(r.systemStub, r.namespace)

	contactPersonID, err := primitive.ObjectIDFromHex(contactPersonUUID)
	if err != nil {
		return nil, models.ErrClientContactPersonNotFound
	}

	var contactPerson ContactPersonInMongo
	err = collection.FindOneAndDelete(ctx, bson.M{"_id": contactPersonID}).Decode(&contactPerson)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, models.ErrClientContactPersonNotFound
		}

		err = errors.Join(errors.New("failed to delete contact person from the database"), err)
		r.logger.Error(err.Error())
		return nil, err
	}

	return contactPerson.ToBackendModel(r.namespace), nil
}
func (r *ClientRepository) GetContactPersonsForClient(ctx context.Context, clientUUID string, useCache bool) ([]models.ContactPerson, error) {
	collection := GetClientContactPersonCollection(r.systemStub, r.namespace)

	clientUUIDBson, err := primitive.ObjectIDFromHex(clientUUID)
	if err != nil {
		return nil, models.ErrClientNotFound
	}

	cursor, err := collection.Find(ctx, primitive.M{"clientUUID": clientUUIDBson})
	if err != nil {
		err = errors.Join(errors.New("failed to get contact persons from the database. failed to open cursor"), err)
		r.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var contactPersons []models.ContactPerson
	for cursor.Next(ctx) {
		var contactPerson ContactPersonInMongo
		err := cursor.Decode(&contactPerson)
		if err != nil {
			err = errors.Join(errors.New("failed to get contact persons from the database. failed to decode contact person"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		contactPersons = append(contactPersons, *contactPerson.ToBackendModel(r.namespace))
	}

	return contactPersons, nil
}

func (r *ClientRepository) GetContactPersonsForClients(ctx context.Context, clientUUIDs []string) (map[string][]models.ContactPerson, error) {
	collection := GetClientContactPersonCollection(r.systemStub, r.namespace)

	clientUUIDsBson := make([]primitive.ObjectID, 0, len(clientUUIDs))
	for _, clientUUID := range clientUUIDs {
		clientUUIDBson, err := primitive.ObjectIDFromHex(clientUUID)
		if err != nil {
			clientUUIDsBson = append(clientUUIDsBson, clientUUIDBson)
		}
	}

	cursor, err := collection.Find(ctx, primitive.M{"clientUUID": bson.M{"$in": clientUUIDsBson}})
	if err != nil {
		err = errors.Join(errors.New("failed to get contact persons from the database. failed to open cursor"), err)
		r.logger.Error(err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	contactPersonsPerClient := make(map[string][]models.ContactPerson, len(clientUUIDs))
	for _, clientUUID := range clientUUIDs {
		contactPersonsPerClient[clientUUID] = make([]models.ContactPerson, 0, 1)
	}

	for cursor.Next(ctx) {
		var contactPerson ContactPersonInMongo
		err := cursor.Decode(&contactPerson)
		if err != nil {
			err = errors.Join(errors.New("failed to get contact persons from the database. failed to decode contact person"), err)
			r.logger.Error(err.Error())
			return nil, err
		}

		contactPersonsPerClient[contactPerson.ClientUUID.Hex()] = append(contactPersonsPerClient[contactPerson.ClientUUID.Hex()], *contactPerson.ToBackendModel(r.namespace))
	}

	return contactPersonsPerClient, nil
}
