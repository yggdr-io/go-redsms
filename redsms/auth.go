package redsms

import (
	"crypto/md5"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// SimpleAuthTransport is an http.RoundTripper that authenticates all requests
// using simple authentication algorithm.
type SimpleAuthTransport struct {
	Login  string
	APIKey string

	// Transport is the underlying HTTP transport to use when making requests.
	Transport http.RoundTripper
}

func (t *SimpleAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *SimpleAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

// RoundTrip implements the RoundTripper interface.
func (t *SimpleAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())

	ts := uuid.New().String()
	s := fmt.Sprintf("%x", md5.Sum([]byte(ts+t.APIKey)))

	req2.Header.Set("login", t.Login)
	req2.Header.Set("ts", ts)
	req2.Header.Set("secret", s)

	return t.transport().RoundTrip(req2)
}
