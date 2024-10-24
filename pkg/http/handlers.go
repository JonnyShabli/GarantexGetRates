package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

const (
	defaultReadHeaderTimeout = 30 * time.Minute
	shutdownTimeout          = 5 * time.Second
)

func NewHandler(basePath string, opts ...RouterOption) http.Handler {
	baseRouter := chi.NewRouter()
	baseRouter.Use(middleware.Logger)
	baseRouter.Route(basePath, func(r chi.Router) {
		for _, opt := range opts {
			opt(r)
		}
	})
	return baseRouter
}

func NewServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ErrorLog:          log.Default(),
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}
}

func RunServer(ctx context.Context, addr string, logger *zap.Logger, handler http.Handler) error {
	logger.With(zap.String("address", addr)).Info("Starting http server")
	server := NewServer(addr, handler)
	errListen := make(chan error, 1)
	go func() {
		errListen <- server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctxShutdown)
		if err != nil {
			return fmt.Errorf("can't shutdown server: %w", err)
		}
		return nil
	case err := <-errListen:
		return fmt.Errorf("can't run server: %w", err)
	}
}
