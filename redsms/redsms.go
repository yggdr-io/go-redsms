package redsms

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://cp.redsms.ru/api/"
)

type Client struct {
	// httpClient communicates with the API.
	httpClient *http.Client

	// Base URL for API requests.
	BaseURL *url.URL
}

// NewClient returns a new RedSMS API client.
// If a nil httpClient is provided, a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		httpClient: httpClient,
		BaseURL:    baseURL,
	}
}
