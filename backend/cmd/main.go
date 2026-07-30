package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"qvarkk/kv2/internal/buildinfo"
	"qvarkk/kv2/internal/config"
	"qvarkk/kv2/internal/httpapi"
	"qvarkk/kv2/internal/observability"
	"syscall"
	"time"
)

func main() {
	config.LoadEnv()
	info := buildinfo.Read()

	logging, err := observability.NewLogger(
		os.Stderr,
		observability.LogConfig{
			Format:    config.LogFormat.GetValue(),
			Level:     config.LogLevel.GetValue(),
			AddSource: false,
		},
		info,
	)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	logger := logging.Logger
	address := ":" + config.APIPort.GetValue()

	logger.Info(
		"server starting",
		slog.String("address", address),
	)

	if err := run(logger, address); err != nil {
		logger.Error("application failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger, address string) error {
	router := httpapi.NewRouter()

	server := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrs := make(chan error, 1)

	go func() {
		serverErrs <- server.ListenAndServe()
	}()

	signalCtx, stopSignals := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignals()

	select {
	case err := <-serverErrs:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return fmt.Errorf("serve HTTP: %w", err)
	case <-signalCtx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	logger.Info("server stopped")
	return nil
}
