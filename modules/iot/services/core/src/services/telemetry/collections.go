package telemetry

import (
	"context"
	"errors"
	"fmt"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const basicMetricsCollectionName = "iot_core_telemetry_basicmetric"
const logCollectionName = "iot_core_telemetry_log"
const eventCollectionName = "iot_core_telemetry_event"

const basicMetricExpirationTime = 60 * 60 * 24 * 3 // Basic metrics will expire after 3 days
const logExpirationTime = 60 * 60 * 24 * 14        // Log entries will expire after two weeks (14 days)

func telemetryDBByNamespace(systemStub *system.SystemStub, namespaceName string) *mongo.Database {
	if namespaceName == "" {
		return systemStub.DB.Database("openbp_global")
	} else {
		return systemStub.DB.Database(fmt.Sprintf("openbp_namespace_%s", namespaceName))
	}
}

func TelemetryBasicMetricCollectionByNamespace(systemStub *system.SystemStub, namespaceName string) *mongo.Collection {
	return telemetryDBByNamespace(systemStub, namespaceName).Collection(basicMetricsCollectionName)
}
func TelemetryLogCollectionByNamespace(systemStub *system.SystemStub, namespaceName string) *mongo.Collection {
	return telemetryDBByNamespace(systemStub, namespaceName).Collection(logCollectionName)
}
func TelemetryEventCollectionByNamespace(systemStub *system.SystemStub, namespaceName string) *mongo.Collection {
	return telemetryDBByNamespace(systemStub, namespaceName).Collection(eventCollectionName)
}

func CreateCollections(ctx context.Context, systemStub *system.SystemStub, namespace string) error {
	db := telemetryDBByNamespace(systemStub, namespace)

	basicMetricsOptions := options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().
		SetGranularity("minutes").
		SetMetaField("deviceUUID").
		SetTimeField("timestamp")).
		SetExpireAfterSeconds(basicMetricExpirationTime)
	err := db.CreateCollection(ctx, basicMetricsCollectionName, basicMetricsOptions)
	if err != nil {
		// If error is not "collection already exist"
		if cmd, ok := err.(mongo.CommandError); !ok || (cmd.Code != 17399 && cmd.Name != "NamespaceExists") {
			return errors.New("failed to create basic metrics collection: " + err.Error())
		}
	}

	logOptions := options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().
		SetGranularity("minutes").
		SetMetaField("deviceUUID").
		SetTimeField("timestamp")).
		SetExpireAfterSeconds(logExpirationTime) // Clear logs after
	err = db.CreateCollection(ctx, logCollectionName, logOptions)
	if err != nil {
		// If error is not "collection already exist"
		if cmd, ok := err.(mongo.CommandError); !ok || (cmd.Code != 17399 && cmd.Name != "NamespaceExists") {
			return errors.New("failed to create log collection: " + err.Error())
		}
	}

	eventOptions := options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().
		SetGranularity("hours").
		SetMetaField("deviceUUID").
		SetTimeField("timestamp"))
	err = db.CreateCollection(ctx, logCollectionName, eventOptions)
	if err != nil {
		// If error is not "collection already exist"
		if cmd, ok := err.(mongo.CommandError); !ok || (cmd.Code != 17399 && cmd.Name != "NamespaceExists") {
			return errors.New("failed to create event collection: " + err.Error())
		}
	}

	return nil
}
