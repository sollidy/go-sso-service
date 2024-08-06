package app

import (
	"log/slog"
	grpcapp "sso-service/internal/app/grpc"
	"sso-service/internal/services/auth"
	"sso-service/internal/storage/pg"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
	Storage *pg.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {

	storage, err := pg.New(log)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
		Storage: storage,
	}
}
