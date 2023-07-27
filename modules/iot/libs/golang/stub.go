package iot

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
)

func getConfigEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type CoreService struct {
	Device device.DeviceServiceClient
	Fleet  fleet.FleetServiceClient
}

type GrpcServiceConfig struct {
	enabled bool
	url     string
}

type StubConfig struct {
	logger *log.Logger

	core GrpcServiceConfig
}

func NewStubConfig() *StubConfig {
	return &StubConfig{
		logger: log.StandardLogger(),
		core: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
	}
}

func (sc *StubConfig) WithLogger(logger *log.Logger) *StubConfig {
	sc.logger = logger
	return sc
}

func (sc *StubConfig) WithCoreService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.core = conf[0]
	} else {
		sc.core = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("IOT_CORE_URL", "iot_core:80"),
		}
	}
	return sc
}

type IOTStub struct {
	Core *CoreService

	log       *log.Logger
	config    *StubConfig
	mu        sync.Mutex
	connected bool
	dials     []*grpc.ClientConn
}

func NewIOTStub(config *StubConfig) *IOTStub {
	return &IOTStub{
		log:       config.logger,
		config:    config,
		mu:        sync.Mutex{},
		connected: false,
		dials:     make([]*grpc.ClientConn, 0),
	}
}

func (n *IOTStub) Connect() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.connected {
		return nil
	}

	if n.config.core.enabled {
		conn, services, err := NewCoreConnection(n.config.core.url)
		if err != nil {
			n.log.Error("Error while connecting to the iot_core service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the iot_core service")
		n.dials = append(n.dials, conn)
		n.Core = services
	}

	n.connected = true
	return nil
}

func (n *IOTStub) Close() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.closeConnections()
}

func (n *IOTStub) closeConnections() {
	for _, dial := range n.dials {
		dial.Close()
	}
	n.dials = make([]*grpc.ClientConn, 0)
	n.connected = false
}
