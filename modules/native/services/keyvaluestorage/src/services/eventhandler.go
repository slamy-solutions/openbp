package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/proto"

	"github.com/nats-io/nats.go"
	system_nats "github.com/slamy-solutions/openbp/modules/system/libs/golang/nats"

	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

const (
	NAMESPACE_CREATION_EVENT_CONSUMER_NAME = "native_keyvaluestorage_namespacecreation"
)

type eventHandlerService struct {
	mongoClient                      *mongo.Client
	namespaceCreateEventSubscription *nats.Subscription
}

func NewEventHandlerService(mongoClient *mongo.Client, natsClient *nats.Conn) (*eventHandlerService, error) {
	service := &eventHandlerService{
		mongoClient:                      mongoClient,
		namespaceCreateEventSubscription: nil,
	}

	js, err := natsClient.JetStream()
	if err != nil {
		return nil, errors.New("Error while opening jetsteram context. " + err.Error())
	}
	_, err = js.AddConsumer("native_namespace_event", &nats.ConsumerConfig{
		Durable:        NAMESPACE_CREATION_EVENT_CONSUMER_NAME,
		Name:           NAMESPACE_CREATION_EVENT_CONSUMER_NAME,
		Description:    "Listens on native_namespace create events for native_keyvaluestorage",
		AckPolicy:      nats.AckExplicitPolicy,
		FilterSubject:  "native.namespace.event.created",
		DeliverSubject: "native.keyvaluestorage.deliver.namespace.create",
		DeliverGroup:   "native.keyvaluestorage.deliver.namespace.create",
	})
	if err != nil {
		return nil, errors.New("Error while creating consumer. " + err.Error())
	}
	subscribtion, err := js.QueueSubscribe("native.namespace.event.created", "native.keyvaluestorage.deliver.namespace.create", service.handleNamespaceCreationEvent, nats.Bind("native_namespace_event", NAMESPACE_CREATION_EVENT_CONSUMER_NAME))
	if err != nil {
		return nil, errors.New("Error while creating subscribtion. " + err.Error())
	}

	service.namespaceCreateEventSubscription = subscribtion

	return service, nil
}

func (s *eventHandlerService) Close() error {
	err := s.namespaceCreateEventSubscription.Unsubscribe()
	if err != nil {
		return errors.New("Error while unsubscribing from namespace events. " + err.Error())
	}
	return nil
}

func (s *eventHandlerService) handleNamespaceCreationEvent(msg *nats.Msg) {
	ctx, span := system_nats.StartTelemetrySpanFromMessage(context.Background(), msg, "Handle namespace creation event")
	defer span.End()

	var namespace namespaceGRPC.Namespace
	err := proto.Unmarshal(msg.Data, &namespace)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to unmarshal namespace from event: "+err.Error())
		span.RecordError(err)
		msg.NakWithDelay(time.Second * 5)
		return
	}

	collection := s.mongoClient.Database(fmt.Sprintf("openbp_namespace_%s", namespace.Name)).Collection("native_keyvaluestorage")
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Options: options.Index().SetName(KEY_INDEX_NAME),
		Keys:    bson.D{bson.E{Key: "key", Value: "hashed"}},
	})
	if err != nil {
		span.SetStatus(codes.Error, "Failed to create index: "+err.Error())
		span.RecordError(err)
		msg.Nak()
		return
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.KeyValue{
		Key:   "namespace",
		Value: attribute.StringValue(namespace.Name),
	})
	msg.Ack()
}
