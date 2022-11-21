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

// NewRequest creates an API request.
// If specified, the body is JSON encoded.
func (c *Client) NewRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
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
