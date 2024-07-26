package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso-service/internal/app"
	"sso-service/internal/config"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Debug("starting application", slog.Any("cfg", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	application.GRPCSrv.Stop()

	log.Info("application stopped", slog.String("signal", sign.String()))
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	tintHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})

	switch env {
	case envLocal:
		log = slog.New(tintHandler)
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
