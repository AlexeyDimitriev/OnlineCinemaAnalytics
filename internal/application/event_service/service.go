package eventservice

import (
	"context"
	"fmt"
	"online-cinema-analytics/internal/application/dto"
	"online-cinema-analytics/internal/application/port"
	"online-cinema-analytics/internal/domain/event"
)

type EventService struct {
	publisher port.KafkaPublisher
}

func NewEventService(pub port.KafkaPublisher) *EventService {
	return &EventService{
		publisher: pub,
	}
}

func (s *EventService) Publish(ctx context.Context, req *dto.CreateEventRequest) (*dto.CreateEventResponse, error) {
	domainEvt := event.NewEvent(
		req.UserID, req.MovieID, req.SessionID,
		event.EventType(req.EventType),
		event.DeviceType(req.DeviceType),
		req.ProgressSeconds,
	)

	if err := event.Validate(domainEvt); err != nil {
		return nil, fmt.Errorf("Domain validation failed: %v", err)
	}

	if err := s.publisher.Publish(ctx, domainEvt); err != nil {
		return nil, fmt.Errorf("Publish failed: %v", err)
	}

	return &dto.CreateEventResponse{EventID: domainEvt.EventID}, nil
}
