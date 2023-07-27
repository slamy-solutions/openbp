package main

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/device"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/fleet"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/telemetry"
	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	systemNATS "github.com/slamy-solutions/openbp/modules/system/libs/golang/nats"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/proto"
)

const namespaceCreationEventConsumerName = "iot_core_namespacecreation"

type EventHandler struct {
	systemStub                       *system.SystemStub
	namespaceCreateEventSubscription *nats.Subscription

	logger *logrus.Entry
}

func NewEventHandler(logger *logrus.Entry, systemStub *system.SystemStub) (*EventHandler, error) {
	handler := &EventHandler{
		systemStub: systemStub,
		logger:     logger,
	}

	js, err := systemStub.Nats.JetStream()
	if err != nil {
		return nil, errors.New("Error while opening jetsteram context. " + err.Error())
	}
	_, err = js.AddConsumer("native_namespace_event", &nats.ConsumerConfig{
		Durable:        namespaceCreationEventConsumerName,
		Name:           namespaceCreationEventConsumerName,
		Description:    "Listens on native_namespace create events for native_iam",
		AckPolicy:      nats.AckExplicitPolicy,
		FilterSubject:  "native.namespace.event.created",
		DeliverSubject: "iot.deliver.namespace.create",
		DeliverGroup:   "iot.deliver.namespace.create",
	})
	if err != nil {
		return nil, errors.New("Error while creating consumer. " + err.Error())
	}
	subscribtion, err := js.QueueSubscribe("native.namespace.event.created", "iot.deliver.namespace.create", handler.handleNamespaceCreationEvent, nats.Bind("native_namespace_event", namespaceCreationEventConsumerName))
	if err != nil {
		return nil, errors.New("Error while creating subscribtion. " + err.Error())
	}

	handler.namespaceCreateEventSubscription = subscribtion

	return handler, nil
}

func (h *EventHandler) Close() error {
	err := h.namespaceCreateEventSubscription.Unsubscribe()
	if err != nil {
		return errors.New("Error while unsubscribing from namespace events. " + err.Error())
	}
	return nil
}

func (h *EventHandler) handleNamespaceCreationEvent(msg *nats.Msg) {
	ctx, span := systemNATS.StartTelemetrySpanFromMessage(context.Background(), msg, "Handle namespace creation event")
	defer span.End()

	logger := h.logger.WithField("event", "namespace_created")

	var namespace namespaceGRPC.Namespace
	err := proto.Unmarshal(msg.Data, &namespace)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to unmarshal namespace from event: "+err.Error())
		logger.Error("Failed to unmarshal namespace from event: " + err.Error())
		span.RecordError(err)
		// TODO: Dead leter queue
		msg.Ack()
		return
	}
	span.SetAttributes(attribute.KeyValue{
		Key:   "namespace",
		Value: attribute.StringValue(namespace.Name),
	})
	logger = logger.WithField("namespace", namespace.Name)

	err = device.CreateIndexesForNamespace(ctx, h.systemStub, namespace.Name)
	if err != nil {
		err = errors.New("error while creating indexes in namespace for device service: " + err.Error())
		logger.Error(err.Error())
		span.SetStatus(codes.Error, err.Error())
		msg.NakWithDelay(time.Second * 5) //TODO: Dead letter queue
		return
	}
	logger.Info("Created indexes for device service database collections.")

	err = fleet.CreateIndexesForNamespace(ctx, h.systemStub, namespace.Name)
	if err != nil {
		err = errors.New("error while creating indexes in namespace for fleet service: " + err.Error())
		logger.Error(err.Error())
		span.SetStatus(codes.Error, err.Error())
		msg.NakWithDelay(time.Second * 5) //TODO: Dead letter queue
		return
	}
	logger.Info("Created indexes for fleet service database collections.")

	err = telemetry.CreateCollections(ctx, h.systemStub, namespace.Name)
	if err != nil {
		err = errors.New("error while creating collections in namespace for telemetry service: " + err.Error())
		logger.Error(err.Error())
		span.SetStatus(codes.Error, err.Error())
		msg.NakWithDelay(time.Second * 5) //TODO: Dead letter queue
		return
	}
	logger.Info("Created telemetry collections in the database.")

	msg.Ack()
	span.SetStatus(codes.Ok, "")
	logger.Info("Successfully handled namespace creation event.")
}
