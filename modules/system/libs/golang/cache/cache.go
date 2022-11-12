package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type BackendConfig struct {
	// Base redis connection options
	RedisOptions *redis.Options
}

type Cache interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	BulkSet(ctx context.Context, entries ...CacheEntry) error
	Remove(ctx context.Context, keys ...string) error
	Get(ctx context.Context, key string) ([]byte, error)
	BulkGet(ctx context.Context, keys ...string) ([][]byte, error)
	Has(ctx context.Context, key string) (bool, error)

	Shutdown(ctx context.Context) error
}

type cache struct {
	rdb    *redis.Client
	tracer trace.Tracer
}

type CacheEntry struct {
	// Key under wich to store data
	Key string
	// Stored data
	Value []byte
	// After what time in miliseconds, cache entry should be deleted. Use 0 if you want this entry last as long as possible
	Timeout time.Duration
}

// Creates new cache instance and initializes tracing
func New(url string) (Cache, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(options)
	tracer := otel.GetTracerProvider().Tracer("github.com/slamy-solutions/openbp/modules/system/libs/go/cache")
	return &cache{rdb: rdb, tracer: tracer}, nil
}

// Set cache under specified key
func (c cache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	ctx, span := c.tracer.Start(ctx, "cache.set")
	defer span.End()
	span.SetAttributes(attribute.String("db.type", "redis"))
	err := c.rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	return err
}

// Set multiple cache values at once
func (c cache) BulkSet(ctx context.Context, entries ...CacheEntry) error {
	ctx, span := c.tracer.Start(ctx, "cache.bulkSet")
	defer span.End()
	span.SetAttributes(attribute.String("db.type", "redis"), attribute.Int("cache.bulkCount", len(entries)))

	pipe := c.rdb.Pipeline()
	for _, entry := range entries {
		_, err := pipe.Set(ctx, entry.Key, entry.Value, entry.Timeout).Result()
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			pipe.Discard()
			return err
		}
	}
	_, err := pipe.Exec(ctx)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}

// Remove cache for keys. If cache doesnt exist - does nothing
func (c cache) Remove(ctx context.Context, keys ...string) error {
	ctx, span := c.tracer.Start(ctx, "cache.remove")
	defer span.End()
	span.SetAttributes(attribute.String("db.type", "redis"))

	err := c.rdb.Del(ctx, keys...).Err()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}

// Tries to get cache for key
func (c cache) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, span := c.tracer.Start(ctx, "cache.get")
	defer span.End()
	span.SetAttributes(attribute.String("db.type", "redis"))

	val, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			span.SetAttributes(attribute.Bool("cache.hit", false))
			span.SetStatus(codes.Ok, "")
			return nil, nil
		}
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.Bool("cache.hit", true))
	span.SetStatus(codes.Ok, "")

	return val, nil
}

// Get multiple cache values at once
func (c cache) BulkGet(ctx context.Context, keys ...string) ([][]byte, error) {
	ctx, span := c.tracer.Start(ctx, "cache.bulkGet")
	defer span.End()
	span.SetAttributes(attribute.String("db.type", "redis"), attribute.Int("cache.bulkCount", len(keys)))

	result, err := c.rdb.MGet(ctx, keys...).Result()

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	} else {
		span.SetStatus(codes.Ok, "")
	}

	cachedData := make([][]byte, len(keys))
	hits := 0

	for cmdIndex, cmd := range result {
		data, _ := cmd.(*redis.StringCmd).Bytes()
		cachedData[cmdIndex] = data
		if data != nil {
			hits += 1
		}
	}

	span.SetAttributes(attribute.Int("cache.hits", hits))

	return cachedData, nil
}

// Checks if there is cache for specified key. Returns false on error
func (c cache) Has(ctx context.Context, key string) (bool, error) {
	ctx, span := c.tracer.Start(ctx, "cache.has")
	defer span.End()
	span.SetAttributes(attribute.String("db.type", "redis"))

	val, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return false, err
	} else {
		has := val == 1
		span.SetAttributes(attribute.Bool("cache.hit", has))
		span.SetStatus(codes.Ok, "")
		return has, nil
	}
}

func (c cache) Shutdown(ctx context.Context) error {
	return c.rdb.Close()
}
