package http

import (
	"net/http"
	"time"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewAPIClient() *APIClient {
	return &APIClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *APIClient) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
}
