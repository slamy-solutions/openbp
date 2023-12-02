package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/proto"

	"github.com/nats-io/nats.go"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	system_nats "github.com/slamy-solutions/openbp/modules/system/libs/golang/nats"

	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/services/storage/src/services/bucket"
	file "github.com/slamy-solutions/openbp/modules/native/services/storage/src/services/fs"
)

const (
	NAMESPACE_CREATION_EVENT_CONSUMER_NAME = "native_file_namespacecreation"
)

type eventHandlerService struct {
	systemStub                       *system.SystemStub
	namespaceCreateEventSubscription *nats.Subscription
	logger                           *slog.Logger
}

func NewEventHandlerService(systemStub *system.SystemStub, logger *slog.Logger) (*eventHandlerService, error) {
	service := &eventHandlerService{
		systemStub:                       systemStub,
		namespaceCreateEventSubscription: nil,
		logger:                           logger.With("type", "event_handler"),
	}

	js, err := systemStub.Nats.JetStream()
	if err != nil {
		return nil, errors.New("Error while opening jetsteram context. " + err.Error())
	}
	_, err = js.AddConsumer("native_namespace_event", &nats.ConsumerConfig{
		Durable:        NAMESPACE_CREATION_EVENT_CONSUMER_NAME,
		Name:           NAMESPACE_CREATION_EVENT_CONSUMER_NAME,
		Description:    "Listens on native_namespace create events for native_file",
		AckPolicy:      nats.AckExplicitPolicy,
		FilterSubject:  "native.namespace.event.created",
		DeliverSubject: "native.file.deliver.namespace.create",
		DeliverGroup:   "native.file.deliver.namespace.create",
	})
	if err != nil {
		return nil, errors.New("Error while creating consumer. " + err.Error())
	}
	subscribtion, err := js.QueueSubscribe("native.namespace.event.created", "native.file.deliver.namespace.create", service.handleNamespaceCreationEvent, nats.Bind("native_namespace_event", NAMESPACE_CREATION_EVENT_CONSUMER_NAME))
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
		s.logger.Error("Failed to unmarshal namespace from event.", "err", err.Error())
		// TODO: Dead leter queue
		msg.Ack()
		return
	}
	span.SetAttributes(attribute.KeyValue{
		Key:   "namespace",
		Value: attribute.StringValue(namespace.Name),
	})

	badAnswer := func(err error) {
		s.logger.Error("Failed to handle namespace creation event.", "namespace", namespace.Name, "err", err.Error())
		span.SetStatus(codes.Error, err.Error())
		//TODO: Dead letter queue
		msg.NakWithDelay(time.Second * 5)
	}

	err = bucket.HandleNamespaceCreationEvent(ctx, &namespace, s.systemStub)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for bucket service: " + err.Error()))
		return
	}
	err = file.HandleNamespaceCreationEvent(ctx, &namespace, s.systemStub)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for file service: " + err.Error()))
		return
	}

	s.logger.Info("Namespace creation event handled successfully.", "namespace", namespace.Name)
	span.SetStatus(codes.Ok, "")
	msg.Ack()
}
