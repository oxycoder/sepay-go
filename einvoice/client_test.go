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
