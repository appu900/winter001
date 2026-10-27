package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/appu900/backend001/config"
	"github.com/appu900/backend001/infra/caching/redis"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// contexts for graceful shutfown and another for redis and database connection
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	withTimeoutctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	defer stop()

	
	//starting the redis connection
	redisClient, err  := redis.New(withTimeoutctx,cfg,logger)
	if err != nil {
		panic(err)
	}
	log.Println(redisClient)
	
	
	app := &Application{cfg: cfg, logger: logger}
	return app.Run(ctx)
}
