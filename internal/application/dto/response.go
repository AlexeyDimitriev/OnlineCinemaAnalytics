package dto

type CreateEventResponse struct {
	EventID string `json:"event_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
