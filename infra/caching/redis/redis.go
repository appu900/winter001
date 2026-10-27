package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/appu900/backend001/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
	logger *slog.Logger
}

const (
	RedisPoolsize      = 20
	MinIdealConnection = 5
	MaxRetries         = 3
	DialTimeout        = 5 * time.Second
	ReadTimeout        = 5 * time.Second
	WriteTimeout       = 3 * time.Second
	Pooltimeout        = 5 * time.Second
	ConnIdeallifeTime  = 5 * time.Minute
	ConnMaxLifetime    = 30 * time.Minute
)

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*RedisClient, error) {
	var opts *redis.Options
	opts = &redis.Options{
		Addr:     cfg.RedisAddr,
		Username: cfg.RedisUsername,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	opts.PoolSize = RedisPoolsize
	opts.MinIdleConns = MinIdealConnection
	opts.MaxRetries = MaxRetries
	opts.WriteTimeout = WriteTimeout
	opts.ReadTimeout = ReadTimeout
	opts.PoolTimeout = Pooltimeout
	opts.ConnMaxIdleTime = ConnIdeallifeTime

	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		logger.Error("error", "redis", "connection failed")
		return nil, fmt.Errorf("connect to redis at %s: %w", cfg.RedisAddr, err)
	}
    logger.Info("Redis connected",cfg.Env)
	return &RedisClient{rdb: rdb, logger: logger}, nil
}
