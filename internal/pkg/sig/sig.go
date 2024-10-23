package sig

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var ErrSignalReceived = errors.New("operating system signal")

func ListenSignal(ctx context.Context, logger *zap.Logger) error {
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return nil
	case sig := <-sigquit:
		logger.With(zap.String("signal", sig.String()))
		logger.Info("Gracefully shutting down server...")
		return ErrSignalReceived
	}
}
