package system

import (
	"context"
	"errors"
	"os"

	"github.com/nats-io/nats.go"
	goredis "github.com/redis/go-redis/v9"
	cache "github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/db"
	system_nats "github.com/slamy-solutions/openbp/modules/system/libs/golang/nats"
	otel "github.com/slamy-solutions/openbp/modules/system/libs/golang/otel"
	redis "github.com/slamy-solutions/openbp/modules/system/libs/golang/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type CacheConfig struct {
	Enabled bool
	URL     string
}
type RedisConfig struct {
	Enabled bool
	URL     string
}
type NatsConfig struct {
	Enabled    bool
	URL        string
	ClientName string
}
type DBConfig struct {
	Enabled bool
	URL     string
}
type OTelConfig struct {
	Enabled bool
	URL     string

	ServiceModule     string
	ServiceName       string
	ServiceVersion    string
	ServiceInstanceID string
}

func getConfigEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewOTelConfig(module string, name string, version string, instanceID string) *OTelConfig {
	return &OTelConfig{
		Enabled:           true,
		URL:               getConfigEnv("SYSTEM_TELEMETRY_EXPORTER_ENDPOINT", "system_telemetry:55680"),
		ServiceModule:     module,
		ServiceName:       name,
		ServiceVersion:    version,
		ServiceInstanceID: instanceID,
	}
}

func (c *OTelConfig) WithURL(url string) *OTelConfig {
	c.URL = url
	return c
}

type SystemStubConfig struct {
	Redis RedisConfig
	Cache CacheConfig
	Nats  NatsConfig
	Db    DBConfig
	OTel  OTelConfig
}

func (s *SystemStubConfig) WithCache(config ...CacheConfig) *SystemStubConfig {
	cfg := &CacheConfig{
		Enabled: true,
		URL:     getConfigEnv("SYSTEM_CACHE_URL", "redis://system_cache"),
	}
	if len(config) > 0 {
		cfg = &config[0]
	}

	s.Cache = *cfg
	return s
}

func (s *SystemStubConfig) WithRedis(config ...RedisConfig) *SystemStubConfig {
	cfg := &RedisConfig{
		Enabled: true,
		URL:     getConfigEnv("SYSTEM_REDIS_URL", "redis://system_redis"),
	}
	if len(config) > 0 {
		cfg = &config[0]
	}

	s.Redis = *cfg
	return s
}

func (s *SystemStubConfig) WithNats(config ...NatsConfig) *SystemStubConfig {
	cfg := &NatsConfig{
		Enabled:    true,
		URL:        getConfigEnv("SYSTEM_NATS_URL", "nats://system_nats:4222"),
		ClientName: "",
	}
	if len(config) > 0 {
		cfg = &config[0]
	}

	s.Nats = *cfg
	return s
}

func (s *SystemStubConfig) WithDB(config ...DBConfig) *SystemStubConfig {
	cfg := &DBConfig{
		Enabled: true,
		URL:     getConfigEnv("SYSTEM_DB_URL", "mongodb://root:example@system_db/admin"),
	}
	if len(config) > 0 {
		cfg = &config[0]
	}

	s.Db = *cfg
	return s
}

func (s *SystemStubConfig) WithOTel(config *OTelConfig) *SystemStubConfig {
	s.OTel = *config
	return s
}

func NewSystemStubConfig() *SystemStubConfig {
	return &SystemStubConfig{
		Redis: RedisConfig{
			Enabled: false,
			URL:     "",
		},
		Cache: CacheConfig{
			Enabled: false,
			URL:     "",
		},
		Nats: NatsConfig{
			Enabled:    false,
			URL:        "",
			ClientName: "",
		},
		Db: DBConfig{
			Enabled: false,
			URL:     "",
		},
		OTel: OTelConfig{
			Enabled:           false,
			URL:               "",
			ServiceModule:     "",
			ServiceName:       "",
			ServiceVersion:    "",
			ServiceInstanceID: "",
		},
	}
}

type SystemStub struct {
	Redis *goredis.Client
	Cache cache.Cache
	DB    *mongo.Client
	OTel  otel.Telemetry
	Nats  *nats.Conn

	config *SystemStubConfig
}

func NewSystemStub(config *SystemStubConfig) *SystemStub {
	return &SystemStub{
		config: config,
	}
}

func (s *SystemStub) Connect(ctx context.Context) error {
	if s.config.OTel.Enabled {
		tel, err := otel.Register(ctx, s.config.OTel.URL, s.config.OTel.ServiceModule, s.config.OTel.ServiceName, s.config.OTel.ServiceVersion, s.config.OTel.ServiceInstanceID)
		if err != nil {
			return errors.New("failed to initialize connection to the otel: " + err.Error())
		}
		s.OTel = tel
	}

	if s.config.Cache.Enabled {
		cacheClient, err := cache.New(s.config.Cache.URL)
		if err != nil {
			//Close opened connections
			if s.config.OTel.Enabled {
				s.OTel.Shutdown(ctx)
			}

			return errors.New("failed to initialize connection to the cache: " + err.Error())
		}
		s.Cache = cacheClient
	}

	if s.config.Redis.Enabled {
		redisClient, err := redis.ConnectToRedis(s.config.Redis.URL)
		if err != nil {
			//Close opened connections
			if s.config.Cache.Enabled {
				s.Cache.Shutdown(ctx)
			}
			if s.config.OTel.Enabled {
				s.OTel.Shutdown(ctx)
			}

			return errors.New("failed to initialize connection to the redis: " + err.Error())
		}
		s.Redis = redisClient
	}

	if s.config.Db.Enabled {
		dbClient, err := db.Connect(s.config.Db.URL)
		if err != nil {
			//Close opened connections
			if s.config.Redis.Enabled {
				s.Redis.Close()
			}
			if s.config.Cache.Enabled {
				s.Cache.Shutdown(ctx)
			}
			if s.config.OTel.Enabled {
				s.OTel.Shutdown(ctx)
			}

			return errors.New("failed to initialize connection to the DB: " + err.Error())
		}

		s.DB = dbClient
	}

	if s.config.Nats.Enabled {
		natsClient, err := system_nats.Connect(s.config.Nats.URL, s.config.Nats.ClientName)
		if err != nil {
			//Close opened connections
			if s.config.Db.Enabled {
				s.DB.Disconnect(ctx)
			}
			if s.config.Redis.Enabled {
				s.Redis.Close()
			}
			if s.config.Cache.Enabled {
				s.Cache.Shutdown(ctx)
			}
			if s.config.OTel.Enabled {
				s.OTel.Shutdown(ctx)
			}

			return errors.New("failed to initialize connection to the Nats: " + err.Error())
		}

		s.Nats = natsClient
	}

	return nil
}

func (s *SystemStub) Close(ctx context.Context) {
	if s.config.Redis.Enabled {
		s.Redis.Close()
	}
	if s.config.Cache.Enabled {
		s.Cache.Shutdown(ctx)
	}
	if s.config.Db.Enabled {
		s.DB.Disconnect(ctx)
	}
	if s.config.Nats.Enabled {
		s.Nats.Close()
	}
	if s.config.OTel.Enabled {
		s.OTel.Shutdown(ctx)
	}
}
