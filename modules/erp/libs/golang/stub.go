package erp

import (
	"log/slog"
	"os"
	"sync"

	catalogGRPC "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"
	"google.golang.org/grpc"
)

func getConfigEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type CatalogServices struct {
	Catalog catalogGRPC.CatalogServiceClient
	Entry   catalogGRPC.CatalogEntryServiceClient
}

type CoreService struct {
	Catalog CatalogServices
}

type GrpcServiceConfig struct {
	enabled bool
	url     string
}

type StubConfig struct {
	logger *slog.Logger

	core GrpcServiceConfig
}

func NewStubConfig() *StubConfig {
	return &StubConfig{
		logger: slog.Default(),
		core: GrpcServiceConfig{
			enabled: false,
			url:     "",
		},
	}
}

func (sc *StubConfig) WithLogger(logger *slog.Logger) *StubConfig {
	sc.logger = logger
	return sc
}

func (sc *StubConfig) WithCoreService(conf ...GrpcServiceConfig) *StubConfig {
	if len(conf) != 0 {
		sc.core = conf[0]
	} else {
		sc.core = GrpcServiceConfig{
			enabled: true,
			url:     getConfigEnv("ERP_CORE_URL", "erp_core:80"),
		}
	}
	return sc
}

type ERPStub struct {
	Core *CoreService

	log       *slog.Logger
	config    *StubConfig
	mu        sync.Mutex
	connected bool
	dials     []*grpc.ClientConn
}

func NewERPStub(config *StubConfig) *ERPStub {
	return &ERPStub{
		log:       config.logger,
		config:    config,
		mu:        sync.Mutex{},
		connected: false,
		dials:     make([]*grpc.ClientConn, 0),
	}
}

func (n *ERPStub) Connect() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.connected {
		return nil
	}

	if n.config.core.enabled {
		conn, services, err := NewCoreConnection(n.config.core.url)
		if err != nil {
			n.log.Error("Error while connecting to the erp_core service: " + err.Error())
			n.closeConnections()
			return err
		}
		n.log.Info("Successfully connected to the erp_core service")
		n.dials = append(n.dials, conn)
		n.Core = services
	}

	n.connected = true
	return nil
}

func (n *ERPStub) Close() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.closeConnections()
}

func (n *ERPStub) closeConnections() {
	for _, dial := range n.dials {
		dial.Close()
	}
	n.dials = make([]*grpc.ClientConn, 0)
	n.connected = false
}
