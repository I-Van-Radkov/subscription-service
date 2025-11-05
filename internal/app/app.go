package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/I-Van-Radkov/subscription-service/internal/config"
	v1 "github.com/I-Van-Radkov/subscription-service/internal/controller/http/v1"
	postgres "github.com/I-Van-Radkov/subscription-service/pkg/db"
	"github.com/I-Van-Radkov/subscription-service/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	httpServer *v1.Server
	postgresDb *postgres.Database
	logger     logger.Logger
}

func NewApp(cfg *config.Config, lg logger.Logger) (*App, error) {
	db, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	server := v1.NewServer(cfg.Port, cfg.ReadTimeout, cfg.WriteTimeout, db.Pool)
	err = server.RegisterHandlers()
	if err != nil {
		return nil, fmt.Errorf("failed to register handlers: %w", err)
	}

	return &App{
		httpServer: server,
		postgresDb: db,
		logger:     lg,
	}, nil
}

func (a *App) MustRun(ctx context.Context, port int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.logger.Info(ctx, fmt.Sprintf("HTTP server listening on port %d", port))
		if err := a.httpServer.Start(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error(ctx, "server error", zap.Error(err))
		}
	}()

	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt, syscall.SIGTERM)
	<-graceSh

	a.logger.Info(ctx, "Shutdown signal received, starting graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := a.httpServer.Stop(shutdownCtx); err != nil {
		a.logger.Error(ctx, "server shutdown error", zap.Error(err))

	}

	a.postgresDb.Close()
	a.logger.Info(ctx, "Database connection pool closed")

	wg.Wait()
	a.logger.Info(ctx, "Server stopped gracefully")
}
