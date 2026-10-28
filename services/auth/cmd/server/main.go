package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/appu900/winter/pkg/logger"
	"github.com/appu900/winter/services/auth/internal/config"
)

func main() {
	if err := Bootstarp(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}
func Bootstarp() error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config error %v", err)
	}
	logger, err := logger.New("auth", cfg.Env, "debug")
	if err != nil {
		log.Fatalf("init logger %v", err)
	}
	defer logger.Sync()

	authApplication := NewAuthApplication(cfg, logger)

	//context for graceful shudown server
	graceShtdownContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return authApplication.Start(graceShtdownContext)

}
