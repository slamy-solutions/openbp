package services

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"

	system_cache "github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"
	system_db "github.com/slamy-solutions/openbp/modules/system/libs/golang/db"
	system_nats "github.com/slamy-solutions/openbp/modules/system/libs/golang/nats"
	system_otel "github.com/slamy-solutions/openbp/modules/system/libs/golang/otel"
)

type SystemConnectionHandler struct {
	Telemetry system_otel.Telemetry
	Cache     system_cache.Cache
	DB        *mongo.Client
	Nats      *nats.Conn
}

func ConnectToSystemServices(ctx context.Context) (*SystemConnectionHandler, error) {
	SYSTEM_DB_URL := getConfigEnv("SYSTEM_DB_URL", "mongodb://root:example@system_db/admin")
	SYSTEM_CACHE_URL := getConfigEnv("SYSTEM_CACHE_URL", "redis://system_cache")
	SYSTEM_NATS_URL := getConfigEnv("SYSTEM_NATS_URL", "nats://system_nats")
	SYSTEM_TELEMETRY_EXPORTER_ENDPOINT := getConfigEnv("SYSTEM_TELEMETRY_EXPORTER_ENDPOINT", "system_telemetry:55680")

	telemetryProvider, err := system_otel.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "tools", "rest", "1.0.0", "1")
	if err != nil {
		panic(err)
	}

	cache, err := system_cache.New(SYSTEM_CACHE_URL)
	if err != nil {
		telemetryProvider.Shutdown(ctx)
		return nil, err
	}

	dbConn, err := system_db.Connect(SYSTEM_DB_URL)
	if err != nil {
		telemetryProvider.Shutdown(ctx)
		cache.Shutdown(ctx)
		return nil, err
	}

	natsConn, err := system_nats.Connect(SYSTEM_NATS_URL, "native_rest")
	if err != nil {
		telemetryProvider.Shutdown(ctx)
		cache.Shutdown(ctx)
		dbConn.Disconnect(ctx)
		return nil, err
	}

	return &SystemConnectionHandler{
		Telemetry: telemetryProvider,
		Cache:     cache,
		DB:        dbConn,
		Nats:      natsConn,
	}, nil
}

func (h *SystemConnectionHandler) Shutdown(ctx context.Context) {
	h.Telemetry.Shutdown(ctx)
	h.Cache.Shutdown(ctx)
	h.DB.Disconnect(ctx)
	h.Nats.Drain()
}
