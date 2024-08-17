package app

import (
	"log/slog"
	grpcapp "sso-service/internal/app/grpc"
	"sso-service/internal/services/auth"
	eventsender "sso-service/internal/services/event-sender"
	"sso-service/internal/storage"
	"sso-service/internal/storage/repository"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
	Storage *storage.Storage
	Sender  *eventsender.Sender
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {

	a := &App{}

	a.Storage = storage.New(log)

	repo := repository.New(a.Storage.DB)
	authService := auth.New(log, repo, repo, repo, tokenTTL)

	a.Sender = eventsender.New(log, repo)
	a.GRPCSrv = grpcapp.New(log, authService, grpcPort)

	return a
}
