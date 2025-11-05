package main

import (
	"context"
	"fmt"
	"os"

	"github.com/I-Van-Radkov/subscription-service/internal/app"
	"github.com/I-Van-Radkov/subscription-service/internal/config"
	"github.com/I-Van-Radkov/subscription-service/pkg/logger"
)

func main() {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "./config/.env"
	}

	cfg, err := config.ParseConfigFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	lg := logger.NewLogger(cfg.Env)

	ctx := logger.WithRequestID(context.Background(), "12345678")
	lg.Info(ctx, "starting server")

	app, err := app.NewApp(cfg, lg)
	if err != nil {
		panic(fmt.Errorf("failed to creating the app structure: %w", err))
	}

	app.MustRun(ctx, cfg.Port, cfg.ReadTimeout)
}
