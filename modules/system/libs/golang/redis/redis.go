package redis

import (
	"errors"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func ConnectToRedis(url string) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.New("failed to parse redis connection string: " + err.Error())
	}
	rdb := redis.NewClient(options)

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, errors.New("failed to instrument redis tracing: " + err.Error())
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		return nil, errors.New("failed to instrument redis metrics: " + err.Error())
	}

	return rdb, nil
}
