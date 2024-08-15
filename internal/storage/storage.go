package storage

import (
	"errors"
	"fmt"
	"log/slog"
	"sso-service/prisma/db"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)

type Storage struct {
	log *slog.Logger
	DB  *db.PrismaClient
}

func New(log *slog.Logger) *Storage {
	return &Storage{log: log, DB: db.NewClient()}
}

func (s *Storage) MustConnect() {
	if err := s.Connect(); err != nil {
		panic(err)
	}
}

func (s *Storage) Connect() error {
	const op = "storage.Connect"
	if err := s.DB.Prisma.Connect(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	s.log.With(slog.String("op", op)).Info("connected to database")
	return nil
}

func (s *Storage) Close() error {
	const op = "storage.Close"
	if err := s.DB.Prisma.Disconnect(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	s.log.With(slog.String("op", op)).Info("DISCONNECTED from database")
	return nil
}
