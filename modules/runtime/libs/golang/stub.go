package iot

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	rpc "github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/rpc"
	runtime "github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/runtime"
)

func getConfigEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type ManagerService struct {
	Runtime runtime.RuntimeServiceClient
	RPC     rpc.RPCServiceClient
}

type GrpcServiceConfig struct {
	enabled bool
	url     string
}

type StubConfig struct {
	logger *log.Logger

	manager GrpcServiceConfig
}

func NewStubConfig() *StubConfig {
	return &StubConfig{
		logger: log.StandardLogger(),
		manager: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
	}
}

func (sc *StubConfig) WithLogger(logger *log.Logger) *StubConfig {
	sc.logger = logger
	return sc
}

func (sc *StubConfig) WithManagerService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.manager = conf[0]
	} else {
		sc.manager = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("RUNTIME_MANAGER_URL", "runtime_manager:80"),
		}
	}
	return sc
}

type RuntimeStub struct {
	Manager *ManagerService

	log       *log.Logger
	config    *StubConfig
	mu        sync.Mutex
	connected bool
	dials     []*grpc.ClientConn
}

func NewRuntimeStub(config *StubConfig) *RuntimeStub {
	return &RuntimeStub{
		log:       config.logger,
		config:    config,
		mu:        sync.Mutex{},
		connected: false,
		dials:     make([]*grpc.ClientConn, 0),
	}
}

func (n *RuntimeStub) Connect() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.connected {
		return nil
	}

	if n.config.manager.enabled {
		conn, services, err := NewManagerConnection(n.config.manager.url)
		if err != nil {
			n.log.Error("Error while connecting to the runtime_manager service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the runtime_manager service")
		n.dials = append(n.dials, conn)
		n.Manager = services
	}

	n.connected = true
	return nil
}

func (n *RuntimeStub) Close() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.closeConnections()
}

func (n *RuntimeStub) closeConnections() {
	for _, dial := range n.dials {
		dial.Close()
	}
	n.dials = make([]*grpc.ClientConn, 0)
	n.connected = false
}
