package redsms

import (
	"encoding/json"
	"io"
	"net/http"
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

	req, err := c.NewRequest("GET", inURL, inBody)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}

	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %s, want %s", inURL, got, want)
	}

	body, _ := io.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest Body is %s, want %s", got, want)
	}
	if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("NewRequest Content-Type is %s, want %s", got, want)
	}
}

func TestNewRequest_badMethod(t *testing.T) {
	c := NewClient(nil)

	_, err := c.NewRequest("FOO\nBAR", ".", nil)

	if err == nil {
		t.Errorf("Expected error to be returned")
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)

	_, err := c.NewRequest("GET", ":", nil)

	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %v", err)
	}
}

func TestNewRequest_urlLeadingSlash(t *testing.T) {
	testcases := map[string]struct {
		inURL string
		want  string
	}{
		"with leading slash": {
			inURL: "/foo",
			want:  defaultBaseURL + "foo",
		},
		"without leading slash": {
			inURL: "foo",
			want:  defaultBaseURL + "foo",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			c := NewClient(nil)
			req, _ := c.NewRequest("GET", tc.inURL, nil)
			if got := req.URL.String(); got != tc.want {
				t.Errorf("NewRequest(%q) URL is %s, want %s",
					tc.inURL, got, tc.want)
			}
		})
	}
}

func TestNewRequest_badBody(t *testing.T) {
	c := NewClient(nil)

	type T struct {
		A map[interface{}]interface{}
	}
	_, err := c.NewRequest("GET", ".", &T{})

	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON unsupported type error, got %v", err)
	}
}

func TestNewRequest_noBody(t *testing.T) {
	c := NewClient(nil)

	req, err := c.NewRequest("GET", ".", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}

	if req.Body != nil {
		t.Errorf("Constructed request contains a non-nil Body")
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{ErrorMessage: "m", Response: res}

	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}
