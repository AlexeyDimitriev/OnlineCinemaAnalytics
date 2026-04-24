package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"online-cinema-analytics/internal/application/dto"
	"online-cinema-analytics/internal/application/event_service"
	"online-cinema-analytics/internal/infrastructure/logger"
)

type Engine struct {
	svc *eventservice.EventService
	cfg struct {
		Users int
		Movies int
		Interval time.Duration
	}
}

func NewEngine(svc *eventservice.EventService, users, movies int, interval time.Duration) *Engine {
	return &Engine{
		svc: svc,
		cfg: struct {
			Users    int
			Movies   int
			Interval time.Duration
		}{
			Users: users,
			Movies: movies,
			Interval: interval,
		},
	}
}

func (e *Engine) Run(ctx context.Context) {
	ticker := time.NewTicker(e.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.emitSession()
		}
	}
}

func (e *Engine) emitSession() {
	userID := fmt.Sprintf("user-%d", rand.Intn(e.cfg.Users)+1)
	movieID := fmt.Sprintf("movie-%d", rand.Intn(e.cfg.Movies)+1)
	sessionID := fmt.Sprintf("sess-%d", time.Now().UnixNano())

	scenario := GenerateCompleteView(userID, movieID, sessionID)

	for _, evt := range scenario.Events {
		req := &dto.CreateEventRequest{
			UserID: evt.UserID,
			MovieID: evt.MovieID,
			EventType: string(evt.EventType),
			DeviceType: string(evt.DeviceType),
			SessionID: evt.SessionID,
			ProgressSeconds: evt.ProgressSeconds,
		}

		if _, err := e.svc.Publish(context.Background(), req); err != nil {
			lgr := logger.NewLogger()
			lgr.Info("Error while publishing: %v", err.Error())
		}
		time.Sleep(200 * time.Millisecond)
	}
}
