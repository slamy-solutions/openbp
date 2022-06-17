package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func Connect(url string) (*mongo.Client, error) {
	dbOptions := &options.ClientOptions{}
	dbOptions.Monitor = otelmongo.NewMonitor(otelmongo.WithCommandAttributeDisabled(true))
	dbOptions.ApplyURI(url)
	err := dbOptions.Validate()
	if err != nil {
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(timeoutCtx, dbOptions)
	if err != nil {
		return nil, err
	}
	return dbClient, nil
}
