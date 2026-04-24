package port

import (
	"context"
	"online-cinema-analytics/internal/domain/event"
)

type KafkaPublisher interface {
	Publish(ctx context.Context, event *event.Event) error
}
