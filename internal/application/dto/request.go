package dto

type CreateEventRequest struct {
	UserID string `json:"user_id" binding:"required"`
	MovieID string `json:"movie_id" binding:"required"`
	EventType string `json:"event_type" binding:"required,oneof=VIEW_STARTED VIEW_FINISHED VIEW_PAUSED VIEW_RESUMED LIKED SEARCHED"`
	DeviceType string `json:"device_type" binding:"required,oneof=MOBILE DESKTOP TV TABLET"`
	SessionID string `json:"session_id" binding:"required"`
	ProgressSeconds int32 `json:"progress_seconds" binding:"gte=0"`
}
