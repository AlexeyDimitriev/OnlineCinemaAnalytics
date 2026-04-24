package clients

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"errors"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseClient struct {
 	Conn clickhouse.Conn
}

type EventRecord struct {
	EventID string `ch:"event_id"`
	UserID string `ch:"user_id"`
	MovieID string `ch:"movie_id"`
	EventType string `ch:"event_type"`
	Timestamp time.Time `ch:"timestamp"`
	DeviceType string `ch:"device_type"`
	SessionID string `ch:"session_id"`
	ProgressSeconds int32 `ch:"progress_seconds"`
}

func NewClickHouseClient(host, port, database, user, password string) (*ClickHouseClient, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
		Auth: clickhouse.Auth{
			Database: database,
			Username: user,
			Password: password,
		},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ClickHouse: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("Failed to ping ClickHouse: %v", err)
	}

	return &ClickHouseClient{Conn: conn}, nil
}

func (c *ClickHouseClient) FindEventByID(ctx context.Context, eventID string, maxWait time.Duration) (*EventRecord, error) {
	start := time.Now()
	pollInterval := 500 * time.Millisecond

	for time.Since(start) < maxWait {
		var record EventRecord
		
		err := c.Conn.QueryRow(
			ctx,
			`SELECT event_id, user_id, movie_id, event_type, 
			timestamp, device_type, session_id, progress_seconds
			FROM movie_events
			WHERE event_id = ?
			LIMIT 1`,
			eventID,
		).Scan(
			&record.EventID, &record.UserID, &record.MovieID, &record.EventType,
			&record.Timestamp, &record.DeviceType, &record.SessionID, &record.ProgressSeconds,
		)

		if err == nil {
			return &record, nil
		}

		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Query error: %v", err)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}
	}

	return nil, fmt.Errorf("event not found after %v", maxWait)
}

func (c *ClickHouseClient) DeleteEventsBySessionID(ctx context.Context, sessionID string) error {
	err := c.Conn.Exec(
		ctx,
		`ALTER TABLE movie_events DELETE WHERE session_id = ?`,
		sessionID,
	)
	return err
}

func (c *ClickHouseClient) CountEventsBySessionID(ctx context.Context, sessionID string) (int, error) {
	var count uint64
	err := c.Conn.QueryRow(
		ctx,
		`SELECT count() FROM movie_events WHERE session_id = ?`,
		sessionID,
	).Scan(&count)
	return int(count), err
}

func (c *ClickHouseClient) Close() error {
 	return c.Conn.Close()
}
