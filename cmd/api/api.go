package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"github.com/appu900/backend001/config"
)

type Application struct {
	cfg    *config.Config
	logger *slog.Logger
}

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all ok"))
	})
	return mux
}

// run the server bro !!!!!
func (app *Application) Run(ctx context.Context) error {
	srv := http.Server{
		Addr:         app.cfg.HttpAddr,
		Handler:      app.routes(),
		ReadTimeout:  app.cfg.ReadTimeout,
		WriteTimeout: app.cfg.WriteTimeout,
		IdleTimeout:  app.cfg.IdealTimeOut,
	}

	errCh := make(chan error, 1)
	go func() {
		app.logger.Info("server starting", "addr", srv.Addr, "env", app.cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		app.logger.Info("shutting down")
	}

	shutDownContext , cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	return srv.Shutdown(shutDownContext)
}
