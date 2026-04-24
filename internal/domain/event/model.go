package event

import (
	"time"
	"github.com/google/uuid"
)

type Event struct {
	EventID string
	UserID string
	MovieID string
	EventType EventType
	Timestamp time.Time
	DeviceType DeviceType
	SessionID string
	ProgressSeconds int32
}

func NewEvent(userID string, movieID string, sessionID string, eventType EventType, deviceType DeviceType, progressSeconds int32) *Event {
	return &Event{
		EventID: uuid.New().String(),
		UserID: userID,
		MovieID: movieID,
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		DeviceType: deviceType,
		SessionID: sessionID,
		ProgressSeconds: progressSeconds,
	}
}
