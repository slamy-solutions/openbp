package main

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/proto"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	system_nats "github.com/slamy-solutions/openbp/modules/system/libs/golang/nats"

	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/oauth"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/x509"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/identity"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/policy"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/role"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

const (
	NAMESPACE_CREATION_EVENT_CONSUMER_NAME = "native_iam_namespacecreation"
)

type eventHandlerService struct {
	systemStub                       *system.SystemStub
	nativeStub                       *native.NativeStub
	namespaceCreateEventSubscription *nats.Subscription

	policyServer *policy.IAMPolicyServer
	roleServer   *role.IAMRoleServer
}

func NewEventHandlerService(systemStub *system.SystemStub, nativeStub *native.NativeStub, policyServer *policy.IAMPolicyServer, roleServer *role.IAMRoleServer) (*eventHandlerService, error) {
	service := &eventHandlerService{
		systemStub:                       systemStub,
		nativeStub:                       nativeStub,
		policyServer:                     policyServer,
		roleServer:                       roleServer,
		namespaceCreateEventSubscription: nil,
	}

	js, err := systemStub.Nats.JetStream()
	if err != nil {
		return nil, errors.New("Error while opening jetsteram context. " + err.Error())
	}
	_, err = js.AddConsumer("native_namespace_event", &nats.ConsumerConfig{
		Durable:        NAMESPACE_CREATION_EVENT_CONSUMER_NAME,
		Name:           NAMESPACE_CREATION_EVENT_CONSUMER_NAME,
		Description:    "Listens on native_namespace create events for native_iam",
		AckPolicy:      nats.AckExplicitPolicy,
		FilterSubject:  "native.namespace.event.created",
		DeliverSubject: "native.iam.deliver.namespace.create",
		DeliverGroup:   "native.iam.deliver.namespace.create",
	})
	if err != nil {
		return nil, errors.New("Error while creating consumer. " + err.Error())
	}
	subscribtion, err := js.QueueSubscribe("native.namespace.event.created", "native.iam.deliver.namespace.create", service.handleNamespaceCreationEvent, nats.Bind("native_namespace_event", NAMESPACE_CREATION_EVENT_CONSUMER_NAME))
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
	span.SetAttributes(attribute.KeyValue{
		Key:   "namespace",
		Value: attribute.StringValue(namespace.Name),
	})

	// Handle event in the services
	logger := logrus.StandardLogger()

	badAnswer := func(err error) {
		span.SetStatus(codes.Error, err.Error())
		//TODO: Dead letter queue
		msg.NakWithDelay(time.Second * 5)
	}

	err = policy.HandleNamespaceCreationEvent(ctx, logger.WithField("service", "policy"), &namespace, s.systemStub)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for policy service: " + err.Error()))
		return
	}
	err = role.HandleNamespaceCreationEvent(ctx, logger.WithField("service", "role"), &namespace, s.systemStub, s.nativeStub, s.policyServer)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for role service: " + err.Error()))
		return
	}
	err = identity.HandleNamespaceCreationEvent(ctx, logger.WithField("service", "identity"), &namespace, s.systemStub)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for identity service: " + err.Error()))
		return
	}
	err = password.HandleNamespaceCreationEvent(ctx, logger.WithField("service", "authentication_password"), &namespace, s.systemStub)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for authentication_password service: " + err.Error()))
		return
	}
	err = x509.HandleNamespaceCreationEvent(ctx, logger.WithField("service", "authentication_x509"), &namespace, s.systemStub)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for authentication_x509 service: " + err.Error()))
		return
	}
	err = oauth.HandleNamespaceCreationEvent(ctx, logger.WithField("service", "authentication_oauth"), &namespace, s.systemStub)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for authentication_oauth service: " + err.Error()))
		return
	}
	err = user.HandleNamespaceCreationEvent(ctx, logger.WithField("service", "actor_user"), &namespace, s.systemStub)
	if err != nil {
		badAnswer(errors.New("failed to handle creation event for actor_user service: " + err.Error()))
		return
	}

	span.SetStatus(codes.Ok, "")
	msg.Ack()
}
