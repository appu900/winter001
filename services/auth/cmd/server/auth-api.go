package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/appu900/winter/services/auth/internal/config"
	"go.uber.org/zap"
)

type AuthApplication struct {
	cfg    *config.ApplicationConfig
	logger *zap.Logger
}

func NewAuthApplication(cfg *config.ApplicationConfig, logger *zap.Logger) *AuthApplication {
	return &AuthApplication{
		cfg:    cfg,
		logger: logger,
	}
}

func (ap *AuthApplication) HttpRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all ok"))
	})
	return mux
}

// start the server
func (ap *AuthApplication) Start(ctx context.Context) error {
	server := http.Server{
		Addr:         ap.cfg.ServerAddrs,
		Handler:      ap.HttpRouter(),
		WriteTimeout: ap.cfg.WriteTimeout,
		ReadTimeout:  ap.cfg.ReadTimeout,
	}

	errCh := make(chan error, 1)
	//start the server in a background go routine
	go func() {
		ap.logger.Info("starting the server", zap.String("addr", ap.cfg.ServerAddrs), zap.String("env", ap.cfg.Env))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		ap.logger.Info("shutting down the server")
	}

	// cancel all inflight request
	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Shutdown(shutdownContext)
}
