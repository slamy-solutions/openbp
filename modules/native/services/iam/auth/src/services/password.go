package services

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type passwordIdentificationService struct {
	mongoClient           *mongo.Client
	mongoPrefix           string
	globalMongoCollection *mongo.Collection
}

type PasswordIdentificationService interface {
	// Searches for identity and validates its password. Returns true only if identity was found and password is valid.
	Authorize(ctx context.Context, namespace string, identity string, password string) (bool, error)
	// Creates or updates password for identity
	CreateOrUpdate(ctx context.Context, namespace string, identity string, password string) error
	// Deletes identity password identification method
	Delete(ctx context.Context, namespace string, identity string) error
}

type passwordInMongo struct {
	Identity string
	Password []byte
}

func collectionByNamespace(s *passwordIdentificationService, namespace string) *mongo.Collection {
	if namespace == "" {
		return s.globalMongoCollection
	} else {
		dbName := fmt.Sprintf("%snamespace_%s", s.mongoPrefix, namespace)
		return s.mongoClient.Database(dbName).Collection("native_iam_auth_authentication_password")
	}
}

func NewPasswordIdentificationService(mongoClient *mongo.Client, mongoPrefix string) PasswordIdentificationService {
	return &passwordIdentificationService{
		mongoClient:           mongoClient,
		mongoPrefix:           mongoPrefix,
		globalMongoCollection: mongoClient.Database(fmt.Sprintf("%sglobal", mongoPrefix)).Collection("native_iam_auth_authentication_password"),
	}
}

func (s *passwordIdentificationService) Authorize(ctx context.Context, namespace string, identity string, password string) (bool, error) {
	collection := collectionByNamespace(s, namespace)
	var entry passwordInMongo
	err := collection.FindOne(ctx, bson.M{"identity": identity}).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(entry.Password, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, errors.New("Failed to compare passwords: " + err.Error())
	}

	return true, nil
}

func (s *passwordIdentificationService) CreateOrUpdate(ctx context.Context, namespace string, identity string, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return errors.New("Failed to hash password: " + err.Error())
	}

	entry := &passwordInMongo{
		Identity: identity,
		Password: passwordHash,
	}
	collection := collectionByNamespace(s, namespace)
	_, err = collection.UpdateOne(ctx, bson.M{"identity": identity}, bson.M{"$set": entry, "$setOnInsert": entry}, options.Update().SetUpsert(true))
	if err != nil {
		return errors.New("Error on updating password in database: " + err.Error())
	}

	return nil
}

func (s *passwordIdentificationService) Delete(ctx context.Context, namespace string, identity string) error {
	collection := collectionByNamespace(s, namespace)
	_, err := collection.DeleteOne(ctx, bson.M{"identity": identity})
	if err != nil {
		return errors.New("Error on deleting password in database: " + err.Error())
	}
	return nil
}
