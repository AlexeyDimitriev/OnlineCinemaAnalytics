package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"online-cinema-analytics/internal/domain/event"
	"online-cinema-analytics/internal/infrastructure/logger"

	"github.com/IBM/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topic string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0
	cfg.Producer.RequiredAcks = sarama.WaitForAll // acks=all
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 0

	sp, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}
	return &Producer{syncProducer: sp, topic: topic}, nil
}

func (p *Producer) Publish(ctx context.Context, evt *event.Event) error {
   msgData := map[string]interface{}{
       "event_id": evt.EventID,
       "user_id": evt.UserID,
       "movie_id": evt.MovieID,
       "event_type": string(evt.EventType),
       "timestamp": evt.Timestamp.Format(time.RFC3339),
       "device_type": string(evt.DeviceType),
       "session_id": evt.SessionID,
       "progress_seconds": evt.ProgressSeconds,
   }
   data, err := json.Marshal(msgData)
   if err != nil {
       return fmt.Errorf("serialization failed: %w", err)
   }

   msg := &sarama.ProducerMessage{       
       Topic: p.topic,        
       Key: sarama.StringEncoder(evt.UserID), // partitioning        
       Value: sarama.ByteEncoder(data),    
   }

	var lastErr error
	baseDelay := 100 * time.Millisecond
	maxAttempts := 5

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		_, _, err := p.syncProducer.SendMessage(msg)
		if err == nil {
			lgr := logger.NewLogger()
			lgr.Info(
				"Event published",
				"ID", evt.EventID,
				"Type", evt.EventType,
				"Timestamp", evt.Timestamp,
			)
			return nil
		}

		lastErr = err
		if attempt == maxAttempts {
			break
		}

		delay := baseDelay * time.Duration(1 << (attempt - 1))
		select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
		}
	}

	return fmt.Errorf("Failed to publish after %v attempts: %v", maxAttempts, lastErr)
}

func (p *Producer) Close() {
	p.syncProducer.Close()
}
