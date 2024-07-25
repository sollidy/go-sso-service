package auth

import (
	"context"
	"fmt"
	"log/slog"
	"sso-service/internal/domain/models"
	"sso-service/internal/lib/logger/sl"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New returns a new instance of the Auth service.
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		userSaver:    userSaver,
		userProvider: userProvider,
		log:          log,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	panic("not implemented")
}

// RegisterNewUser registers a new user with the given email and password.
// It returns the generated user ID and any error encountered.
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {

	const op = "auth.RegisterNewUser"
	log := a.log.With(slog.String("op", op))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed t generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save new user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

	// TODO: Implement user registration logic
	panic("not implemented")
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}
