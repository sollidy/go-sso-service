package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sso-service/internal/app"
	"sso-service/internal/config"
	"sso-service/internal/lib/logger"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Debug("starting application", slog.Any("cfg", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()
	go application.Storage.MustConnect()
	application.Sender.StartProcessingEvents(ctx, 5*time.Second)

	shutdown(cancel, application, log) // graceful shutdown
}

func shutdown(cancel context.CancelFunc, app *app.App, log *slog.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	cancel()
	app.GRPCSrv.Stop()
	app.Storage.Close()

	log.Warn("STOPED application", slog.String("signal", sign.String()))
}
