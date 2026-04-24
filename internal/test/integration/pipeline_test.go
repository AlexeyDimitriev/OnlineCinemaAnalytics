package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"online-cinema-analytics/internal/test/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullPipeline(t *testing.T) {
	producerURL := getEnv("PRODUCER_URL", "http://localhost:8080")
	chHost := getEnv("CLICKHOUSE_HOST", "localhost")
	chPort := getEnv("CLICKHOUSE_PORT", "9000")
	chDB := getEnv("CLICKHOUSE_DB", "movie_analytics")
	chUser := getEnv("CLICKHOUSE_USER", "analytics")
	chPass := getEnv("CLICKHOUSE_PASSWORD", "analytics")

	ctx := context.Background()

	producer := clients.NewProducerClient(producerURL)
	chClient, err := clients.NewClickHouseClient(chHost, chPort, chDB, chUser, chPass)
	require.NoError(t, err, "Failed to connect to ClickHouse")
	defer chClient.Close()

	// unique session id for isolation
	sessionID := fmt.Sprintf("test-pipeline-%d-%d", time.Now().UnixNano(), os.Getpid())
	userID := "test-user-pipeline"
	movieID := "test-movie-pipeline"
	eventID := ""

	req := clients.CreateEventRequest{
		UserID: userID,
		MovieID: movieID,
		EventType: "VIEW_STARTED",
		DeviceType: "DESKTOP",
		SessionID: sessionID,
		ProgressSeconds: 0,
	}

	eventID, err = producer.PublishEvent(req)
	require.NoError(t, err, "Failed to publish event via HTTP")
	require.NotEmpty(t, eventID, "event_id should not be empty")

	record, err := chClient.FindEventByID(ctx, eventID, 30*time.Second)
	require.NoError(t, err, "Event not found in ClickHouse within timeout")

	assert.Equal(t, eventID, record.EventID, "event_id mismatch")
	assert.Equal(t, userID, record.UserID, "user_id mismatch")
	assert.Equal(t, movieID, record.MovieID, "movie_id mismatch")
	assert.Equal(t, "VIEW_STARTED", record.EventType, "event_type mismatch")
	assert.Equal(t, "DESKTOP", record.DeviceType, "device_type mismatch")
	assert.Equal(t, sessionID, record.SessionID, "session_id mismatch")
	assert.Equal(t, int32(0), record.ProgressSeconds, "progress_seconds mismatch")

	now := time.Now()
	assert.WithinDuration(t, now, record.Timestamp, 5*time.Minute, "Timestamp is too far from now")

	// cleaning
	_ = chClient.DeleteEventsBySessionID(ctx, sessionID)
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}