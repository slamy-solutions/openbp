package services

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/proto"

	"github.com/nats-io/nats.go"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	system_nats "github.com/slamy-solutions/openbp/modules/system/libs/golang/nats"

	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

const (
	NAMESPACE_CREATION_EVENT_CONSUMER_NAME = "native_actor_user_namespacecreation"
)

type eventHandlerService struct {
	systemStub                       *system.SystemStub
	namespaceCreateEventSubscription *nats.Subscription
}

func NewEventHandlerService(systemStub *system.SystemStub) (*eventHandlerService, error) {
	service := &eventHandlerService{
		systemStub:                       systemStub,
		namespaceCreateEventSubscription: nil,
	}

	js, err := systemStub.Nats.JetStream()
	if err != nil {
		return nil, errors.New("Error while opening jetsteram context. " + err.Error())
	}
	_, err = js.AddConsumer("native_namespace_event", &nats.ConsumerConfig{
		Durable:        NAMESPACE_CREATION_EVENT_CONSUMER_NAME,
		Name:           NAMESPACE_CREATION_EVENT_CONSUMER_NAME,
		Description:    "Listens on native_namespace create events for native_actor_user",
		AckPolicy:      nats.AckExplicitPolicy,
		FilterSubject:  "native.namespace.event.created",
		DeliverSubject: "native.actor.user.deliver.namespace.create",
		DeliverGroup:   "native.actor.user.deliver.namespace.create",
	})
	if err != nil {
		return nil, errors.New("Error while creating consumer. " + err.Error())
	}
	subscribtion, err := js.QueueSubscribe("native.namespace.event.created", "native.actor.user.deliver.namespace.create", service.handleNamespaceCreationEvent, nats.Bind("native_namespace_event", NAMESPACE_CREATION_EVENT_CONSUMER_NAME))
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
		// TODO: Dead leter queue
		msg.Ack()
		return
	}

	// Create indexes for user collection inside namespace
	err = ensureIndexesForNamespace(ctx, namespace.Name, s.systemStub)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to create indexes. "+err.Error())
		span.RecordError(err)
		// TODO: Dead leter queue
		msg.NakWithDelay(time.Second * 5)
		return
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.KeyValue{
		Key:   "namespace",
		Value: attribute.StringValue(namespace.Name),
	})
	msg.Ack()
}