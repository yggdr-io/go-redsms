package redsms

import (
	"encoding/json"
	"io"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.BaseURL.String(), defaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %s, want %s", got, want)
	}

	c2 := NewClient(nil)
	if c.httpClient == c2.httpClient {
		t.Error("NewClient returned the same http.Client, but they should differ")
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "foo", defaultBaseURL+"foo"
	inBody := &struct {
		Foo string `json:"foo"`
	}{
		Foo: "bar",
	}
	outBody := `{"foo":"bar"}` + "\n"

	req, _ := c.NewRequest("GET", inURL, inBody)

	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %s, want %s", inURL, got, want)
	}

	body, _ := io.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest() Body is %s, want %s", got, want)
	}
	if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("NewRequest() Content-Type is %s, want %s", got, want)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)

	_, err := c.NewRequest("GET", ":", nil)

	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestNewRequest_badMethod(t *testing.T) {
	c := NewClient(nil)

	_, err := c.NewRequest("FOO\nBAR", "foo", nil)

	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestNewRequest_badBody(t *testing.T) {
	c := NewClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.NewRequest("GET", ".", &T{})

	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a json unsupported type error, got %+v", err)
	}
}
