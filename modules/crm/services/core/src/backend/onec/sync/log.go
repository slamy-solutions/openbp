package sync

import (
	"context"
	"time"

	onecSync "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/onecsync"
	"google.golang.org/protobuf/types/known/timestamppb"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SyncEvent struct {
	UUDI      primitive.ObjectID `json:"uuid" bson:"_id,omitempty"`
	Namespace string             `json:"namespace" bson:"-"`
	Status    bool               `json:"status" bson:"status"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Error     string             `json:"error" bson:"error"`
	Log       string             `json:"log" bson:"log"`
}

func (e *SyncEvent) ToGRPC() *onecSync.OneCSyncEvent {
	return &onecSync.OneCSyncEvent{
		Uuid:      e.UUDI.Hex(),
		Namespace: e.Namespace,
		Success:   e.Status,
		Error:     e.Error,
		Log:       e.Log,
		Timestamp: timestamppb.New(e.Timestamp),
	}
}

const syncLogCollectionName = "crm_onec_sync_log"

func initSyncLogCollection(ctx context.Context, systemStub *system.SystemStub, namespace string) error {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	db := systemStub.DB.Database(dbName)
	tsOptions := options.TimeSeries().SetTimeField("timestamp").SetGranularity("hours")

	return db.CreateCollection(ctx, syncLogCollectionName, options.CreateCollection().SetTimeSeriesOptions(tsOptions).SetExpireAfterSeconds(60*60*24*7))
}

func GetSyncLogCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(syncLogCollectionName)
}

func addSyncLog(ctx context.Context, systemStub *system.SystemStub, event SyncEvent) error {
	syncLogCollection := GetSyncLogCollection(systemStub, event.Namespace)

	_, err := syncLogCollection.InsertOne(ctx, event)

	return err
}

func getSyncLog(ctx context.Context, systemStub *system.SystemStub, namespace string, skip int, limit int) ([]SyncEvent, int, error) {
	syncLogCollection := GetSyncLogCollection(systemStub, namespace)

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	if skip > 0 {
		opts.SetSkip(int64(skip))
	}
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := syncLogCollection.Find(ctx, nil, opts)
	if err != nil {
		return nil, 0, err
	}

	var events []SyncEvent
	err = cursor.All(ctx, &events)
	if err != nil {
		return nil, 0, err
	}

	total, err := syncLogCollection.CountDocuments(ctx, nil)
	if err != nil {
		return nil, 0, err
	}

	for i := range events {
		events[i].Namespace = namespace
	}

	return events, int(total), nil
}
