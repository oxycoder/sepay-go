package einvoice

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestServer returns a client wired to a fake API. The token endpoint is
// handled automatically and hands out "tok-1", "tok-2", ...
func newTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *atomic.Int32) {
	t.Helper()
	var tokens atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/token" {
			id, secret, ok := r.BasicAuth()
			if !ok || id != "id" || secret != "secret" {
				w.WriteHeader(http.StatusUnauthorized)
				io.WriteString(w, `{"success":false,"error":{"code":"UNAUTHORIZED","message":"bad creds"}}`)
				return
			}
			n := tokens.Add(1)
			writeJSON(w, map[string]any{"success": true, "data": map[string]any{
				"access_token": "tok-" + string(rune('0'+n)), "token_type": "Bearer", "expires_in": 86400,
			}})
			return
		}
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	return New("id", "secret", WithBaseURL(srv.URL), WithPollInterval(time.Millisecond)), &tokens
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func TestTokenIsCachedAndSent(t *testing.T) {
	client, tokens := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer tok-1" {
			t.Errorf("Authorization = %q", got)
		}
		writeJSON(w, map[string]any{"data": map[string]any{"quota_remaining": "534"}})
	})

	for range 3 {
		u, err := client.GetUsage(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if u.QuotaRemaining != 534 {
			t.Errorf("QuotaRemaining = %v", u.QuotaRemaining)
		}
	}
	if n := tokens.Load(); n != 1 {
		t.Errorf("token fetched %d times, want 1", n)
	}
}

func TestTokenRefreshedOn401(t *testing.T) {
	client, tokens := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer tok-1" {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, `{"success":false,"error":{"code":"TOKEN_EXPIRED","message":"expired"}}`)
			return
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"reference_code":"ref-1"`) {
			t.Errorf("body not resent on retry: %s", body)
		}
		writeJSON(w, map[string]any{"success": true, "data": map[string]any{"tracking_code": "trk"}})
	})

	res, err := client.IssueInvoice(context.Background(), "ref-1")
	if err != nil {
		t.Fatal(err)
	}
	if res.TrackingCode != "trk" {
		t.Errorf("TrackingCode = %q", res.TrackingCode)
	}
	if n := tokens.Load(); n != 2 {
		t.Errorf("token fetched %d times, want 2", n)
	}
}

func TestPersistent401IsReturned(t *testing.T) {
	client, tokens := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"success":false,"error":{"code":"UNAUTHORIZED","message":"no"}}`)
	})
	_, err := client.GetUsage(context.Background())
	if !HasCode(err, CodeUnauthorized) {
		t.Fatalf("err = %v", err)
	}
	if n := tokens.Load(); n != 2 {
		t.Errorf("token fetched %d times, want 2 (one retry only)", n)
	}
}

func TestBadCredentials(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("API must not be called without a token")
	})
	client.clientSecret = "wrong"
	_, err := client.GetUsage(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 401 || apiErr.Code != CodeUnauthorized {
		t.Fatalf("err = %v", err)
	}
}

func TestAPIErrorDecoding(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/invoices/missing":
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, `{"success":false,"error":{"code":"NOT_FOUND","message":"không tìm thấy"}}`)
		case "/v1/invoices/html":
			w.WriteHeader(http.StatusBadGateway)
			io.WriteString(w, "<html>bad gateway</html>")
		case "/v1/invoices/quota":
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"success":false,"error":{"code":"QUOTA_HAS_BEEN_USERD_UP","message":"hết"}}`)
		}
	})
	ctx := context.Background()

	_, err := client.GetInvoice(ctx, "missing")
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(%v) = false", err)
	}

	_, err = client.GetInvoice(ctx, "html")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 502 || apiErr.Message != "<html>bad gateway</html>" {
		t.Errorf("err = %#v", err)
	}

	_, err = client.GetInvoice(ctx, "quota")
	if !HasCode(err, CodeQuotaUsedUp) {
		t.Errorf("err = %v", err)
	}
}

func TestSuccessFalseWith200IsError(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success":false,"message":"nope"}`)
	})
	err := client.DeleteBusinessLocation(context.Background(), "x")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Message != "nope" {
		t.Fatalf("err = %v", err)
	}
}

func TestOptions(t *testing.T) {
	c := New("id", "secret")
	if c.baseURL != ProductionURL || c.userAgent != defaultUserAgent || c.pollInterval != defaultPollInterval {
		t.Errorf("defaults = %+v", c)
	}
	if c.httpClient == nil || c.httpClient.Timeout != 30*time.Second {
		t.Errorf("default http client = %+v", c.httpClient)
	}

	hc := &http.Client{}
	c = New("id", "secret", WithSandbox(), WithHTTPClient(hc), WithUserAgent("my-app/1.0"))
	if c.baseURL != SandboxURL || c.httpClient != hc || c.userAgent != "my-app/1.0" {
		t.Errorf("options not applied: %+v", c)
	}

	c = New("id", "secret", WithBaseURL("https://example.com///"))
	if c.baseURL != "https://example.com" {
		t.Errorf("baseURL = %q", c.baseURL)
	}
}

func TestRequestHeaders(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if ua := r.Header.Get("User-Agent"); ua != "my-app/1.0" {
			t.Errorf("User-Agent = %q", ua)
		}
		if a := r.Header.Get("Accept"); a != "application/json" {
			t.Errorf("Accept = %q", a)
		}
		io.WriteString(w, `{"success":true,"data":{"quota_remaining":1}}`)
	})
	client.userAgent = "my-app/1.0"
	if _, err := client.GetUsage(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestShortLivedTokenIsNotShortened(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"data": map[string]any{"access_token": "t", "expires_in": 60}})
	}))
	defer srv.Close()
	c := New("id", "secret", WithBaseURL(srv.URL))
	tok, err := c.accessToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if tok != "t" {
		t.Errorf("token = %q", tok)
	}
	if d := time.Until(c.tokenExpiry); d < 50*time.Second || d > 60*time.Second {
		t.Errorf("token expires in %v, want ~60s", d)
	}
}

func TestExpiredTokenIsRefetched(t *testing.T) {
	client, tokens := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success":true}`)
	})
	ctx := context.Background()
	if _, err := client.GetUsage(ctx); err != nil {
		t.Fatal(err)
	}
	client.tokenExpiry = time.Now().Add(-time.Second)
	if _, err := client.GetUsage(ctx); err != nil {
		t.Fatal(err)
	}
	if n := tokens.Load(); n != 2 {
		t.Errorf("token fetched %d times, want 2", n)
	}
}

func TestInvalidateTokenKeepsNewerToken(t *testing.T) {
	c := New("id", "secret")
	c.token = "new"
	c.invalidateToken("old")
	if c.token != "new" {
		t.Errorf("token = %q, want new", c.token)
	}
	c.invalidateToken("new")
	if c.token != "" {
		t.Errorf("token = %q, want empty", c.token)
	}
}

func TestEncodeRequestError(t *testing.T) {
	client, tokens := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("API must not be called")
	})
	err := client.do(context.Background(), http.MethodPost, "/v1/x", nil, map[string]any{"bad": make(chan int)}, nil)
	if err == nil || !strings.Contains(err.Error(), "encode request") {
		t.Fatalf("err = %v", err)
	}
	if n := tokens.Load(); n != 0 {
		t.Errorf("token fetched %d times, want 0", n)
	}
}

func TestBuildRequestError(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("API must not be called")
	})
	err := client.do(context.Background(), "BAD METHOD", "/v1/x", nil, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "build request") {
		t.Fatalf("err = %v", err)
	}
}

func TestTokenBuildRequestError(t *testing.T) {
	c := New("id", "secret", WithBaseURL("http://bad host"))
	_, err := c.GetUsage(context.Background())
	if err == nil || !strings.Contains(err.Error(), "build request") {
		t.Fatalf("err = %v", err)
	}
}

func TestTransportError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close() // nothing listens any more
	c := New("id", "secret", WithBaseURL(srv.URL))
	_, err := c.GetUsage(context.Background())
	var apiErr *APIError
	if err == nil || errors.As(err, &apiErr) || !strings.Contains(err.Error(), "POST /v1/token") {
		t.Fatalf("err = %v", err)
	}
}

// roundTripFunc adapts a function to http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("boom") }

func TestReadResponseError(t *testing.T) {
	hc := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(errReader{}), Request: r}, nil
	})}
	c := New("id", "secret", WithHTTPClient(hc))
	_, err := c.GetUsage(context.Background())
	if err == nil || !strings.Contains(err.Error(), "read response: boom") {
		t.Fatalf("err = %v", err)
	}
}

func TestDecodeErrors(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/usage":
			io.WriteString(w, "not json")
		case "/v1/invoices/ref":
			io.WriteString(w, `{"success":true,"data":{"reference_code":123}}`)
		}
	})
	ctx := context.Background()

	_, err := client.GetUsage(ctx)
	if err == nil || !strings.Contains(err.Error(), "decode response:") {
		t.Errorf("non-JSON 200: err = %v", err)
	}
	_, err = client.GetInvoice(ctx, "ref")
	if err == nil || !strings.Contains(err.Error(), "decode response data") {
		t.Errorf("bad data: err = %v", err)
	}
}

func TestErrorObjectWith200IsError(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"error":{"code":"VALIDATION_ERROR","message":"sai"}}`)
	})
	_, err := client.GetUsage(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 200 || apiErr.Code != CodeValidationError || apiErr.Message != "sai" {
		t.Fatalf("err = %#v", err)
	}
}

func TestEmptyDataIsIgnored(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"success":true,"data":null}`)
	})
	u, err := client.GetUsage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if u.QuotaRemaining != 0 {
		t.Errorf("usage = %+v", u)
	}
}

func TestAPIErrorString(t *testing.T) {
	tests := []struct {
		err  *APIError
		want string
	}{
		{&APIError{StatusCode: 400, Code: CodeValidationError, Message: "sai"}, "einvoice: 400 VALIDATION_ERROR: sai"},
		{&APIError{StatusCode: 502, Message: "bad gateway"}, "einvoice: 502: bad gateway"},
	}
	for _, tt := range tests {
		if got := tt.err.Error(); got != tt.want {
			t.Errorf("Error() = %q, want %q", got, tt.want)
		}
	}
}

func TestErrorHelpersOnOtherErrors(t *testing.T) {
	err := errors.New("plain")
	if IsNotFound(err) || HasCode(err, CodeBadRequest) || IsNotFound(nil) {
		t.Error("helpers must be false for non-API errors")
	}
	if IsNotFound(&APIError{StatusCode: 400}) || HasCode(&APIError{Code: "A"}, "B") {
		t.Error("helpers must not match other statuses or codes")
	}
}
