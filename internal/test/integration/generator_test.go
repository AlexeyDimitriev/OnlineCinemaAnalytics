package integration_test

import (
	"context"
	"testing"
	"time"

	"online-cinema-analytics/internal/test/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// for test generator should be already launched with
// GENERATOR_ENABLED=true
// GENERATOR_USERS_COUNT=1
// GENERATOR_MOVIES_COUNT=1
// and with test- prefixes
func TestGeneratorFlow(t *testing.T) {
	chHost := getEnv("CLICKHOUSE_HOST", "localhost")
	chPort := getEnv("CLICKHOUSE_PORT", "9000")
	chDB := getEnv("CLICKHOUSE_DB", "movie_analytics")
	chUser := getEnv("CLICKHOUSE_USER", "analytics")
	chPass := getEnv("CLICKHOUSE_PASSWORD", "analytics")

	ctx := context.Background()
	chClient, err := clients.NewClickHouseClient(chHost, chPort, chDB, chUser, chPass)
	require.NoError(t, err)
	defer chClient.Close()

	time.Sleep(5 * time.Second)

	t.Log("Fetching generated events for validation...")

	rows, err := chClient.Conn.Query(
		ctx,
		`SELECT event_id, user_id, movie_id, event_type, 
			timestamp, device_type, session_id, progress_seconds
		FROM movie_analytics.movie_events
		WHERE timestamp >= now() - INTERVAL 2 MINUTE
			AND user_id LIKE 'user-%' 
			AND movie_id LIKE 'movie-%'
		ORDER BY session_id, timestamp`,
	)
	require.NoError(t, err, "Failed to query events")
	defer rows.Close()

	sessions := make(map[string][]clients.EventRecord)
	for rows.Next() {
		var e clients.EventRecord
		err := rows.Scan(
			&e.EventID, &e.UserID, &e.MovieID, &e.EventType,
			&e.Timestamp, &e.DeviceType, &e.SessionID, &e.ProgressSeconds,
		)
		require.NoError(t, err, "Failed to scan row")
		sessions[e.SessionID] = append(sessions[e.SessionID], e)
	}

	require.Greater(t, len(sessions), 0, "No generated sessions found")
	t.Logf("Found %d sessions to validate", len(sessions))

	validatedSessions := 0
	for sessionID, events := range sessions {
		if len(events) < 2 {
			continue // too short
		}

		for _, e := range events {
			assert.Equal(t, sessionID, e.SessionID, "Event session_id mismatch")
		}

		first := events[0]
		if first.EventType != "VIEW_STARTED" {
			t.Logf("Session %s doesn't start with VIEW_STARTED, skipping", sessionID)
			continue
		}
		assert.Equal(t, int32(0), first.ProgressSeconds, "VIEW_STARTED should have progress=0")

		for i := 1; i < len(events); i++ {
			assert.GreaterOrEqual(t, events[i].ProgressSeconds, events[i-1].ProgressSeconds,
				"Progress decreased in session %s: event %d (%d) -> %d (%d)", 
				sessionID, i-1, events[i-1].ProgressSeconds, i, events[i].ProgressSeconds)
		}

		for i := 1; i < len(events); i++ {
			assert.True(t, events[i].Timestamp.After(events[i-1].Timestamp) || 
						events[i].Timestamp.Equal(events[i-1].Timestamp),
				"Timestamp decreased in session %s", sessionID)
		}

		validatedSessions++
	}

	require.Greater(t, validatedSessions, 0, "No valid sessions found to assert")
	t.Logf("Successfully validated %d sessions", validatedSessions)
}
