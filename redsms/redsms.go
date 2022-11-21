package redsms

import "net/http"

const (
	defaultBaseURL = "https://cp.redsms.ru/api/"
)

type Client struct {
	// httpClient communicates with the API.
	httpClient *http.Client
}

// NewClient returns a new RedSMS API client.
// If a nil httpClient is provided, a new http.Client will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		httpClient: httpClient,
	}
}
