package event

import (
	"errors"
)

var (
	ErrInvalidEventType = errors.New("Invalid event_type")
	ErrInvalidDeviceType = errors.New("Invalid device_type")
	ErrProgressMustBePositive = errors.New("Progress_seconds must be > 0 for view events")
	ErrProgressMustBeZero = errors.New("Progress_seconds must be 0 for non-view events")
)

func Validate(e *Event) error {
	if !e.EventType.IsValid() {
		return ErrInvalidEventType
	}
	if !e.DeviceType.IsValid() {
		return ErrInvalidDeviceType
	}

	isViewEvent := e.EventType == EventViewStarted || e.EventType == EventViewFinished ||
		e.EventType == EventViewPaused || e.EventType == EventViewResumed

	if isViewEvent && e.ProgressSeconds < 0 {
		return ErrProgressMustBePositive
	}

	if !isViewEvent && e.ProgressSeconds != 0 {
		return ErrProgressMustBeZero
	}

	return nil
}
