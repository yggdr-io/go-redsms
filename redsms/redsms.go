package redsms

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://cp.redsms.ru/api/"
)

type Client struct {
	// httpClient communicates with the API.
	httpClient *http.Client

	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

	// Base URL for API requests.
	BaseURL *url.URL

	// Services used for talking to different parts of the RedSMS API.
	Client *ClientService
}

type service struct {
	client *Client
}

// NewClient returns a new RedSMS API client.
// If a nil httpClient is provided, a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		httpClient: httpClient,
		BaseURL:    baseURL,
	}
	c.common.client = c
	c.Client = (*ClientService)(&c.common)

	return c
}

// NewRequest creates an API request.
// If specified, the body is JSON encoded.
func (c *Client) NewRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
