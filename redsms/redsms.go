package redsms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	c.Message = (*MessageService)(&c.common)

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

// BareDo sends an API request and lets you handle the api response.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			// the context's error is probably more useful
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	err = CheckResponse(resp)
	if err != nil {
		defer resp.Body.Close()
	}

	return &Response{Response: resp}, err
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v,
// or returned as an error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	decErr := json.NewDecoder(resp.Body).Decode(v)
	if decErr == io.EOF {
		// ignore EOF errors caused by empty response body
		decErr = nil
	}
	if decErr != nil {
		return resp, err
	}

	return resp, nil
}

// CheckResponse checks the API response for errors,
// and returns them if present.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	body, err := io.ReadAll(r.Body)
	if err == nil && body != nil {
		json.Unmarshal(body, errorResponse)
	}

	return errorResponse
}

// Response is a RedSMS API response.
type Response struct {
	*http.Response

	// FIXME: It should provide convenient access to pagination links.
}

// ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response `json:"-"`

	ErrorMessage string `json:"error_message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.ErrorMessage)
}
