package repository

import (
	"context"
	"fmt"
	"sso-service/internal/domain/models"
	"sso-service/internal/storage"
	"sso-service/prisma/db"
	"strings"
)

type Repository struct {
	db *db.PrismaClient
}

var (
	errUniqueConstraint  = "Unique constraint failed"
	userCreatedEventType = "UserCreated"
)

func New(db *db.PrismaClient) *Repository {
	return &Repository{db: db}
}

func (s *Repository) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "repository.SaveUser"
	event := models.Event{
		Type:    userCreatedEventType,
		Payload: fmt.Sprintf("User %s created", email),
	}
	addUser := s.db.User.CreateOne(
		db.User.Email.Set(email),
		db.User.PassHash.Set(passHash),
	).Tx()

	err := s.db.Prisma.Transaction(addUser, s.saveEvent(event)).Exec(ctx)

	if err != nil {
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), errUniqueConstraint) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int64(addUser.Result().ID), nil
}

func (s *Repository) saveEvent(event models.Event) db.EventUniqueTxResult {
	addEvent := s.db.Event.CreateOne(
		db.Event.EventType.Set(event.Type),
		db.Event.Payload.Set(event.Payload),
	).Tx()
	return addEvent
}

func (s *Repository) User(ctx context.Context, email string) (models.User, error) {
	const op = "repository.User"
	user, err := s.db.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	userData := models.User{}
	if err != nil {
		if err := db.IsErrNotFound(err); err {
			return userData, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return userData, fmt.Errorf("%s: %w", op, err)
	}
	userData = models.User{
		ID:       int64(user.ID),
		Email:    user.Email,
		PassHash: user.PassHash,
	}
	return userData, nil
}

func (s *Repository) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "repository.IsAdmin"
	user, err := s.db.User.FindUnique(
		db.User.ID.Equals(int(userID)),
	).Exec(ctx)
	if err != nil {
		if err := db.IsErrNotFound(err); err {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return user.IsAdmin, nil
}

func (s *Repository) App(ctx context.Context, appID int) (models.App, error) {
	const op = "repository.App"
	app, err := s.db.App.FindUnique(
		db.App.ID.Equals(appID),
	).Exec(ctx)
	if err != nil {
		if err := db.IsErrNotFound(err); err {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return models.App{
		ID:     int(app.ID),
		Name:   app.Name,
		Secret: app.Secret,
	}, nil
}

func (s *Repository) GetNewEvent(ctx context.Context) (models.Event, error) {
	const op = "repository.GetNewEvent"
	event, err := s.db.Event.FindFirst(
		db.Event.Status.Equals("NEW"),
	).Exec(ctx)
	if err != nil {
		if err := db.IsErrNotFound(err); err {
			return models.Event{}, nil // events with status "NEW" not found
		}
		return models.Event{}, fmt.Errorf("%s: %w", op, err)
	}
	return models.Event{
		Id:      int(event.ID),
		Type:    event.EventType,
		Payload: event.Payload,
	}, nil
}

func (s *Repository) SetEventDone(ctx context.Context, eventId int) error {
	const op = "repository.UpdateEvent"
	_, err := s.db.Event.FindUnique(
		db.Event.ID.Equals(eventId),
	).Update(
		db.Event.Status.Set("DONE"),
	).Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
