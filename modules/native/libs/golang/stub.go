package native

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	actorUserGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	iamAuthGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	iamAuthenticationPasswordGrpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
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

type services struct {
	ActorUser                 actorUserGrpc.ActorUserServiceClient
	IamAuthenticationPassword iamAuthenticationPasswordGrpc.IAMAuthenticationPasswordServiceClient
	IamIdentity               iamIdentityGrpc.IAMIdentityServiceClient
	IamAuth                   iamAuthGrpc.IAMAuthServiceClient
	IamPolicy                 iamPolicyGrpc.IAMPolicyServiceClient
	IamRole                   iamRoleGrpc.IAMRoleServiceClient
	IamToken                  iamTokenGrpc.IAMTokenServiceClient
	Keyvaluestorage           keyvaluestorageGrpc.KeyValueStorageServiceClient
	Namespace                 namespaceGrpc.NamespaceServiceClient
}

type GrpcServiceConfig struct {
	enabled bool
	url     string
}

type StubConfig struct {
	logger *log.Logger

	namespace       GrpcServiceConfig
	keyValueStorage GrpcServiceConfig
	actorUser       GrpcServiceConfig

	iamAuthenticationPassword GrpcServiceConfig
	iamIdentity               GrpcServiceConfig
	iamAuth                   GrpcServiceConfig
	iamPolicy                 GrpcServiceConfig
	iamRole                   GrpcServiceConfig
	iamToken                  GrpcServiceConfig
}

func NewStubConfig() *StubConfig {
	return &StubConfig{
		logger: log.StandardLogger(),
		actorUser: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
		namespace: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
	}
}

func (sc *StubConfig) WithLogger(logger *log.Logger) *StubConfig {
	sc.logger = logger
	return sc
}

func (sc *StubConfig) WithActorUserService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.actorUser = conf[0]
	} else {
		sc.actorUser = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_ACTOR_USER_URL", "native_actor_user:80"),
		}
	}
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

func (sc *StubConfig) WithIAMAuthenticationPasswordService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.iamAuthenticationPassword = conf[0]
	} else {
		sc.iamAuthenticationPassword = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_IAM_AUTHENTICATION_PASSWORD_URL", "native_iam_authentication_password:80"),
		}
	}
	return sc
}

func (sc *StubConfig) WithIAMIdentityService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.iamIdentity = conf[0]
	} else {
		sc.iamIdentity = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_IAM_IDENTITY_URL", "native_iam_identity:80"),
		}
	}
	return sc
}

func (sc *StubConfig) WithIAMPolicyService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.iamPolicy = conf[0]
	} else {
		sc.iamPolicy = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_IAM_POLICY_URL", "native_iam_policy:80"),
		}
	}
	return sc
}

func (sc *StubConfig) WithIAMRoleService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.iamRole = conf[0]
	} else {
		sc.iamRole = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_IAM_ROLE_URL", "native_iam_role:80"),
		}
	}
	return sc
}

func (sc *StubConfig) WithIAMTokenService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.iamToken = conf[0]
	} else {
		sc.iamToken = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_IAM_TOKEN_URL", "native_iam_token:80"),
		}
	}
	return sc
}

func (sc *StubConfig) WithIAMAuthService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.iamAuth = conf[0]
	} else {
		sc.iamAuth = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("NATIVE_IAM_AUTH_URL", "native_iam_auth:80"),
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

	if n.config.actorUser.enabled {
		conn, service, err := NewActorUserConnection(n.config.actorUser.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_actor_user service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_actor_user service")
		n.dials = append(n.dials, conn)
		n.Services.ActorUser = service
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

	if n.config.iamAuthenticationPassword.enabled {
		conn, service, err := NewIAMAuthenticationPasswordConnection(n.config.iamAuthenticationPassword.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_iam_authentication_password service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_iam_authentication_password service")
		n.dials = append(n.dials, conn)
		n.Services.IamAuthenticationPassword = service
	}

	if n.config.iamIdentity.enabled {
		conn, service, err := NewIAMIdentityConnection(n.config.iamIdentity.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_iam_identity service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_iam_identity service")
		n.dials = append(n.dials, conn)
		n.Services.IamIdentity = service
	}

	if n.config.iamPolicy.enabled {
		conn, service, err := NewIAMPolicyConnection(n.config.iamPolicy.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_iam_policy service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_iam_policy service")
		n.dials = append(n.dials, conn)
		n.Services.IamPolicy = service
	}

	if n.config.iamRole.enabled {
		conn, service, err := NewIAMRoleConnection(n.config.iamRole.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_iam_role service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_iam_role service")
		n.dials = append(n.dials, conn)
		n.Services.IamRole = service
	}

	if n.config.iamAuth.enabled {
		conn, service, err := NewIAMAuthConnection(n.config.iamAuth.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_iam_auth service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_iam_auth service")
		n.dials = append(n.dials, conn)
		n.Services.IamAuth = service
	}

	if n.config.iamToken.enabled {
		conn, service, err := NewIAMTokenConnection(n.config.iamToken.url)
		if err != nil {
			n.log.Error("Error while connecting to the native_iam_token service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the native_iam_token service")
		n.dials = append(n.dials, conn)
		n.Services.IamToken = service
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
