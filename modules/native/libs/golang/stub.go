package native

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	iamActorUserGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
	iamAuthGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	iamAuthenticationOAuth2Grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/oauth2"
	iamAuthenticationPasswordGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	iamAuthenticationX509Grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/x509"
	iamIdentityGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	iamPolicyGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	iamRoleGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
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

type IamActorServices struct {
	User iamActorUserGrpc.ActorUserServiceClient
}

type IamAuthenticationServices struct {
	Password iamAuthenticationPasswordGrpc.IAMAuthenticationPasswordServiceClient
	X509     iamAuthenticationX509Grpc.IAMAuthenticationX509ServiceClient
	OAuth    IamAuthenticationOAuthServices
}

type IamAuthenticationOAuthServices struct {
	Config iamAuthenticationOAuth2Grpc.IAMAuthenticationOAuth2ConfigServiceClient
	OAuth2 iamAuthenticationOAuth2Grpc.IAMAuthenticationOAuth2ServiceClient
}

type IAMService struct {
	Actor          *IamActorServices
	Authentication *IamAuthenticationServices
	Identity       iamIdentityGrpc.IAMIdentityServiceClient
	Auth           iamAuthGrpc.IAMAuthServiceClient
	Policy         iamPolicyGrpc.IAMPolicyServiceClient
	Role           iamRoleGrpc.IAMRoleServiceClient
	Token          iamTokenGrpc.IAMTokenServiceClient
}

type services struct {
	IAM             *IAMService
	Keyvaluestorage keyvaluestorageGrpc.KeyValueStorageServiceClient
	Namespace       namespaceGrpc.NamespaceServiceClient
}

type GrpcServiceConfig struct {
	enabled bool
	url     string
}

type StubConfig struct {
	logger *log.Logger

	namespace       GrpcServiceConfig
	keyValueStorage GrpcServiceConfig

	iam GrpcServiceConfig
}

func NewStubConfig() *StubConfig {
	return &StubConfig{
		logger: log.StandardLogger(),
		namespace: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
		keyValueStorage: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
		iam: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
	}
}

func (sc *StubConfig) WithLogger(logger *log.Logger) *StubConfig {
	sc.logger = logger
	return sc
}

func (sc *StubConfig) WithNamespaceService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.namespace = conf[0]
	} else {
		sc.namespace = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_NAMESPACE_URL", "native_namespace:80"),
		}
	}
	return sc
}

func (sc *StubConfig) WithKeyValueStorageService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.keyValueStorage = conf[0]
	} else {
		sc.keyValueStorage = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_KEYVALUESTORAGE_URL", "native_keyvaluestorage:80"),
		}
	}
	return sc
}

func (sc *StubConfig) WithIAMService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.iam = conf[0]
	} else {
		sc.iam = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_IAM_URL", "native_iam:80"),
		}
	}
	return sc
}

type NativeStub struct {
	Services services

	log       *log.Logger
	config    *StubConfig
	mu        sync.Mutex
	connected bool
	dials     []*grpc.ClientConn
}

func NewNativeStub(config *StubConfig) *NativeStub {
	return &NativeStub{
		log:       config.logger,
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

	if n.config.namespace.enabled {
		conn, service, err := NewNamespaceConnection(n.config.namespace.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_namespace service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_namespace service")
		n.dials = append(n.dials, conn)
		n.Services.Namespace = service
	}

	if n.config.keyValueStorage.enabled {
		conn, service, err := NewKeyValueStorageConnection(n.config.keyValueStorage.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_iam_keyValueStorage service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_iam_keyValueStorage service")
		n.dials = append(n.dials, conn)
		n.Services.Keyvaluestorage = service
	}

	if n.config.iam.enabled {
		conn, services, err := NewIAMConnection(n.config.iam.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_iam service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_iam service")
		n.dials = append(n.dials, conn)
		n.Services.IAM = services
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
