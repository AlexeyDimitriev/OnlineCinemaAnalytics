package generator

import (
	"math/rand"
	"time"
	"online-cinema-analytics/internal/domain/event"
)

type Scenario struct {
	Events []*event.Event
}

func GenerateCompleteView(userID, movieID, sessionID string) Scenario {
	baseTime := time.Now().UTC()
	events := make([]*event.Event, 0, 4)

	// start of view
	started := event.NewEvent(userID, movieID, sessionID, event.EventViewStarted, event.DeviceDesktop, 0)
	started.Timestamp = baseTime
	events = append(events, started)

	// may pause
	if rand.Float32() > 0.75 {
		progress := int32(rand.Intn(300) + 60)
		paused := event.NewEvent(userID, movieID, sessionID, event.EventViewPaused, event.DeviceDesktop, progress)
		paused.Timestamp = baseTime.Add(time.Duration(progress) * time.Second)
		events = append(events, paused)

		// resume
		pauseDur := int32(rand.Intn(60) + 10)
		resumed := event.NewEvent(userID, movieID, sessionID, event.EventViewResumed, event.DeviceDesktop, progress)
		resumed.Timestamp = paused.Timestamp.Add(time.Duration(pauseDur) * time.Second)
		events = append(events, resumed)
		progress += 1
	}

	// end view
	finalProgress := int32(rand.Intn(1200) + 600)
	finished := event.NewEvent(userID, movieID, sessionID, event.EventViewFinished, event.DeviceDesktop, finalProgress)
	finished.Timestamp = baseTime.Add(time.Duration(finalProgress) * time.Second)
	events = append(events, finished)

	return Scenario{Events: events}
}
