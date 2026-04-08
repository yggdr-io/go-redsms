package redsms

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

	ctl := "\x7f" // ASCII control character
	_, err := c.NewRequest("GET", ctl, nil)

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
		A map[any]any
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

// newTestServer creates a test HTTP server and a Client configured to talk to it.
func newTestServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	srv := httptest.NewServer(handler)
	c := NewClient(srv.Client())
	u, _ := url.Parse(srv.URL + "/api/")
	c.BaseURL = u
	return srv, c
}

func TestBareDo_nilContext(t *testing.T) {
	c := NewClient(nil)
	req, _ := http.NewRequest("GET", "/", nil)
	//nolint:staticcheck // SA1012: intentionally passing nil context to test guard
	_, err := c.BareDo(nil, req)
	if err == nil {
		t.Fatal("expected error for nil context")
	}
	if err != errNonNilContext {
		t.Errorf("expected errNonNilContext, got %v", err)
	}
}

func TestBareDo_success(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	req, _ := c.NewRequest("GET", "test", nil)
	resp, err := c.BareDo(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestBareDo_httpError(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error_message":"bad"}`)
	})
	defer srv.Close()

	req, _ := c.NewRequest("GET", "test", nil)
	_, err := c.BareDo(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for 400 response")
	}
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("expected *ErrorResponse, got %T", err)
	}
	if errResp.ErrorMessage != "bad" {
		t.Errorf("expected error_message 'bad', got %q", errResp.ErrorMessage)
	}
}

func TestBareDo_canceledContext(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	req, _ := c.NewRequest("GET", "test", nil)
	_, err := c.BareDo(ctx, req)
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestDo_success(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name":"test"}`)
	})
	defer srv.Close()

	req, _ := c.NewRequest("GET", "test", nil)
	var result struct {
		Name string `json:"name"`
	}
	resp, err := c.Do(context.Background(), req, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
	if result.Name != "test" {
		t.Errorf("expected name 'test', got %q", result.Name)
	}
}

func TestDo_errorResponse(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"error_message":"fail"}`)
	})
	defer srv.Close()

	req, _ := c.NewRequest("GET", "test", nil)
	var result struct{}
	_, err := c.Do(context.Background(), req, &result)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckResponse_success(t *testing.T) {
	resp := &http.Response{StatusCode: http.StatusOK}
	if err := CheckResponse(resp); err != nil {
		t.Errorf("expected no error for 200, got %v", err)
	}
}

func TestCheckResponse_errorWithBadJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(newStringReader("not json")),
		Request:    &http.Request{},
	}
	err := CheckResponse(resp)
	if err == nil {
		t.Fatal("expected error for 400")
	}
	errResp, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("expected *ErrorResponse, got %T", err)
	}
	if errResp.ErrorMessage != "unknown error (json unmarshaling)" {
		t.Errorf("expected fallback error message, got %q", errResp.ErrorMessage)
	}
}

func TestSimpleAuthTransport_RoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("login") == "" {
			t.Error("expected login header to be set")
		}
		if r.Header.Get("ts") == "" {
			t.Error("expected ts header to be set")
		}
		if r.Header.Get("secret") == "" {
			t.Error("expected secret header to be set")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	tp := &SimpleAuthTransport{
		Login:     "testlogin",
		APIKey:    "testapikey",
		Transport: srv.Client().Transport,
	}
	client := tp.Client()
	resp, err := client.Get(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestSimpleAuthTransport_defaultTransport(t *testing.T) {
	tp := &SimpleAuthTransport{Login: "l", APIKey: "k"}
	if tp.transport() != http.DefaultTransport {
		t.Error("expected default transport when Transport is nil")
	}
}

func TestClientService_GetInfo(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"info":{"login":"user1","balance":100.5,"active":true}}`)
	})
	defer srv.Close()

	info, resp, err := c.Client.GetInfo(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
	if info.Login != "user1" {
		t.Errorf("expected login 'user1', got %q", info.Login)
	}
	if info.Balance != 100.5 {
		t.Errorf("expected balance 100.5, got %f", info.Balance)
	}
	if !info.Active {
		t.Error("expected active to be true")
	}
}

func TestClientService_GetInfo_error(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error_message":"unauthorized"}`)
	})
	defer srv.Close()

	_, _, err := c.Client.GetInfo(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestMessageService_Send(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"items":[{"uuid":"abc","to":"+123"}],"errors":[],"success":true}`)
	})
	defer srv.Close()

	msg := &Message{
		From:  "sender",
		To:    "+123",
		Text:  "hello",
		Route: MessageRouteSMS,
	}
	report, resp, err := c.Message.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
	if !report.Success {
		t.Error("expected success to be true")
	}
	if len(report.Items) != 1 || report.Items[0].UUID != "abc" {
		t.Errorf("unexpected report items: %+v", report.Items)
	}
}

func TestMessageService_Send_error(t *testing.T) {
	srv, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"error_message":"forbidden"}`)
	})
	defer srv.Close()

	msg := &Message{From: "s", To: "+1", Text: "t", Route: MessageRouteSMS}
	_, _, err := c.Message.Send(context.Background(), msg)
	if err == nil {
		t.Fatal("expected error for 403 response")
	}
}

// newStringReader is a helper to create an io.Reader from a string.
func newStringReader(s string) io.Reader {
	return io.NopCloser(io.LimitReader(
		readerFunc(func(p []byte) (int, error) {
			n := copy(p, s)
			return n, io.EOF
		}), int64(len(s)),
	))
}

type readerFunc func(p []byte) (int, error)

func (f readerFunc) Read(p []byte) (int, error) { return f(p) }
