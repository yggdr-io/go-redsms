package redsms

import "testing"

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseUtl is %s, want %s", got, want)
	}

	c2 := NewClient(nil)
	if c.httpClient == c2.httpClient {
		t.Error("NewClient returned the same http.Client, but they should differ")
	}
}
