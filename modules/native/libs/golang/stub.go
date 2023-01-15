package native

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	actorUserGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	iamAuthenticationPasswordGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	iamIdentityGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	iamOAuthGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/oauth"
	iamPolicyGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	iamTokenGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	keyvaluestorageGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	namespaceGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

func getConfigEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type services struct {
	actorUser                 *actorUserGrpc.ActorUserServiceClient
	iamAuthenticationPassword *iamAuthenticationPasswordGrpc.IAMAuthenticationPasswordServiceClient
	iamIdentity               *iamIdentityGrpc.IAMIdentityServiceClient
	iamOAuth                  *iamOAuthGrpc.IAMOAuthServiceClient
	iamPolicy                 *iamPolicyGrpc.IAMPolicyServiceClient
	iamToken                  *iamTokenGrpc.IAMTokenServiceClient
	keyvaluestorage           *keyvaluestorageGrpc.KeyValueStorageServiceClient
	namespace                 *namespaceGrpc.NamespaceServiceClient
}

type GrpcServiceConfig struct {
	enabled bool
	url     string
}

type StubConfig struct {
	actorUser GrpcServiceConfig
}

func NewStubConfig() *StubConfig {
	return &StubConfig{
		actorUser: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
	}
}

func (sc *StubConfig) WithActorUserService(conf ...GrpcServiceConfig) {
	if len(conf) != 0 {
		sc.actorUser = conf[0]
	} else {
		sc.actorUser = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_NAMESPACE_URL", "native_actor_user:80"),
		}
	}
}

type NativeStub struct {
	Services services

	log       log.Logger
	config    StubConfig
	mu        sync.Mutex
	connected bool
	dials     []*grpc.ClientConn
}

func NewNativeStub(config StubConfig, logger log.Logger) *NativeStub {
	return &NativeStub{
		log:       logger,
		config:    config,
		mu:        sync.Mutex{},
		connected: false,
		dials:     make([]*grpc.ClientConn, 0),
	}
}

func (n *NativeStub) Connect() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.connected {
		return nil
	}

	if n.config.actorUser.enabled {
		conn, service, err := NewActorUserConnection(n.config.actorUser.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_actor_user service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_actor_user service")
		n.dials = append(n.dials, conn)
		n.Services.actorUser = &service
	}

	n.connected = true
	return nil
}

func (n *NativeStub) Close() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.closeConnections()
}

func (n *NativeStub) closeConnections() {
	for _, dial := range n.dials {
		dial.Close()
	}
	n.dials = make([]*grpc.ClientConn, 0)
	n.connected = false
}
