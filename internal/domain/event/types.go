package event

type EventType string

const (
	EventViewStarted EventType = "VIEW_STARTED"
	EventViewFinished EventType = "VIEW_FINISHED"
	EventViewPaused EventType = "VIEW_PAUSED"
	EventViewResumed EventType = "VIEW_RESUMED"
	EventLiked EventType = "LIKED"
	EventSearched EventType = "SEARCHED"
)

func (e EventType) IsValid() bool {
	switch e {
	case EventViewStarted, EventViewFinished, EventViewPaused, EventViewResumed, EventLiked, EventSearched:
		return true
	}
	return false
}

type DeviceType string

const (
	DeviceMobile DeviceType = "MOBILE"
	DeviceDesktop DeviceType = "DESKTOP"
	DeviceTV DeviceType = "TV"
	DeviceTablet DeviceType = "TABLET"
)

func (d DeviceType) IsValid() bool {
	switch d {
	case DeviceMobile, DeviceDesktop, DeviceTV, DeviceTablet:
		return true
	}
	return false
}
