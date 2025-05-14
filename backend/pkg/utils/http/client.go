package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Client interface defines the methods required for HTTP operations
type Client interface {
	Post(url string, body interface{}) (*Response, error)
}

// Response wraps the standard http.Response
type Response struct {
	StatusCode int
	Body       []byte
}

// HTTPClient implements the Client interface using standard http.Client
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTPClient instance
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{},
	}
}

// Post implements the Client interface
func (c *HTTPClient) Post(url string, body interface{}) (*Response, error) {
	// Convert body to JSON
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Return wrapped response
	return &Response{
		StatusCode: resp.StatusCode,
		Body:       bodyBytes,
	}, nil
}
