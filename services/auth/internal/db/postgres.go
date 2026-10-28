package db

import (
	"context"
	"fmt"
	"time"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgressInstance(ctx context.Context, connstring string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connstring)
	if err != nil {
		return nil, fmt.Errorf("error in parsing the database %w", err)
	}
	cfg.MaxConns = 10
	cfg.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(ctx,cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool %w", err)
	}

	pingctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("error in connecting database %w", err)
	}
	return pool, nil
}
