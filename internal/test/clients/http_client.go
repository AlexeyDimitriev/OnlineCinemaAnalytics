package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ProducerClient struct {
	baseURL string
	client *http.Client
}

type CreateEventRequest struct {
	UserID string `json:"user_id"`
	MovieID string `json:"movie_id"`
	EventType string `json:"event_type"`
	DeviceType string `json:"device_type"`
	SessionID string `json:"session_id"`
	ProgressSeconds int32 `json:"progress_seconds"`
}

type CreateEventResponse struct {
 	EventID string `json:"event_id"`
}

func NewProducerClient(baseURL string) *ProducerClient {
	return &ProducerClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ProducerClient) PublishEvent(req CreateEventRequest) (string, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("%s/api/v1/events", c.baseURL)
	resp, err := c.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Unexpected status %v: %v", resp.StatusCode, string(respBody))
	}

	var result CreateEventResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("Failed to unmarshal response: %v", err)
	}

	return result.EventID, nil
}
