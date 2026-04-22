package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SchemaRegistryClient struct {
	baseURL string
	httpClient *http.Client
}

type RegisterSchemaRequest struct {
	Schema string `json:"schema"`
	SchemaType string `json:"schemaType,omitempty"`
}

type RegisterSchemaResponse struct {
	ID int `json:"id"`
	Version int `json:"version"`
	Schema string `json:"schema"`
}

type GetSchemaResponse struct {
	Subject string `json:"subject"`
	Version int `json:"version"`
	ID int `json:"id"`
	Schema string `json:"schema"`
}

func NewSchemaRegistryClient(baseURL string) *SchemaRegistryClient {
	return &SchemaRegistryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *SchemaRegistryClient) RegisterSchema(subject string, schema string) (int, error) {
	reqBody := RegisterSchemaRequest{
		Schema: schema,
		SchemaType: "AVRO",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("Failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("%s/subjects/%s/versions", c.baseURL, subject)

	resp, err := c.httpClient.Post(
		url,
		"application/vnd.schemaregistry.v1+json",
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return 0, fmt.Errorf("Failed to register schema: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp := make(map[string]interface{}, 0)
		json.NewDecoder(resp.Body).Decode(&errResp)
		return 0, fmt.Errorf(
			"Failed to register schema, code: %v, body: %v",
			resp.StatusCode, errResp,
		)
	}

	var res RegisterSchemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, fmt.Errorf("Failed to decode response: %v", err)
	}
	return res.Version, nil
}

func (c *SchemaRegistryClient) GetSchema(subject string, version int) (string, error) {
	url := fmt.Sprintf("%s/subjects/%s/versions", c.baseURL, subject)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("Failed to get schema: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to get schema, status code: %v", resp.StatusCode)
	}

	var res GetSchemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("Failed to decode response: %v", err)
	}
	return res.Schema, nil
}

func (c *SchemaRegistryClient) GetLatestSchema(subject string) (string, error) {
	return c.GetSchema(subject, -1)
}
