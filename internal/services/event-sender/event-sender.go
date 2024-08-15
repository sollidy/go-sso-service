package eventsender

import (
	"context"
	"log/slog"
	"sso-service/internal/domain/models"
	"sso-service/internal/lib/logger/sl"
	"time"
)

type Sender struct {
	log           *slog.Logger
	eventProvider EventProvider
}

type EventProvider interface {
	GetNewEvent(ctx context.Context) (models.Event, error)
	SetEventDone(ctx context.Context, eventId int) error
}

func New(log *slog.Logger, eventProvider EventProvider) *Sender {
	return &Sender{log: log, eventProvider: eventProvider}
}

func (s *Sender) StartProcessingEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "services.event-sender.StartProcessingEvents"
	log := s.log.With(slog.String("op", op))
	ticker := time.NewTicker(handlePeriod)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("STOPED processing events")
				return
			case <-ticker.C:
			}
			event, err := s.eventProvider.GetNewEvent(ctx)
			if err != nil {
				log.Error("failed to get new event", sl.Err(err))
				continue
			}
			if event.Id == 0 {
				continue
			}

			s.SendMessage(ctx, event)

			if err := s.eventProvider.SetEventDone(ctx, event.Id); err != nil {

				log.Error("failed to get new event", sl.Err(err))
				continue
			}
			// log.Info("event processedlog.Error("failed to get new event",
		}
	}()
}

func (s *Sender) SendMessage(ctx context.Context, event models.Event) error {
	const op = "services.event-sender.SendMessage"
	log := s.log.With(slog.String("op", op))
	// TODO: send message to external service
	log.Info("sending message", slog.Any("event", event))
	return nil
}
