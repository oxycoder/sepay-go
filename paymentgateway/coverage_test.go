package paymentgateway

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOptions(t *testing.T) {
	hc := &http.Client{}
	c := New("M", "S",
		WithBaseURL("https://api.example/"),
		WithCheckoutURL("https://pay.example/"),
		WithHTTPClient(hc),
		WithUserAgent("my-app/1.0"),
	)
	if c.baseURL != "https://api.example" || c.checkoutURL != "https://pay.example" ||
		c.httpClient != hc || c.userAgent != "my-app/1.0" || c.MerchantID() != "M" {
		t.Errorf("client = %+v", c)
	}

	c = New("M", "S")
	if c.baseURL != ProductionURL || c.checkoutURL != ProductionCheckoutURL || c.userAgent != defaultUserAgent {
		t.Errorf("defaults = %+v", c)
	}
}

func TestUserAgentSent(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if ua := r.Header.Get("User-Agent"); ua != "my-app/1.0" {
			t.Errorf("User-Agent = %q", ua)
		}
		if r.Header.Get("Accept") != "application/json" || r.Header.Get("Content-Type") != "" {
			t.Errorf("headers = %v", r.Header)
		}
		io.WriteString(w, `{"data":{"order_id":"X"}}`)
	})
	WithUserAgent("my-app/1.0")(client)
	if _, err := client.GetOrder(context.Background(), "X"); err != nil {
		t.Fatal(err)
	}
}

func TestNewCheckoutNilRequest(t *testing.T) {
	if _, err := New("M", "S").NewCheckout(nil); err == nil {
		t.Error("want error")
	}
}

func TestCallEncodeError(t *testing.T) {
	err := New("M", "S").call(context.Background(), http.MethodPost, "/x", nil, make(chan int), nil)
	if err == nil || !strings.Contains(err.Error(), "encode request") {
		t.Errorf("err = %v", err)
	}
}

func TestCallBuildRequestError(t *testing.T) {
	c := New("M", "S", WithBaseURL("http://bad\x7fhost"))
	if err := c.CancelOrder(context.Background(), "DH1"); err == nil || !strings.Contains(err.Error(), "build request") {
		t.Errorf("err = %v", err)
	}
}

func TestCallTransportError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close()
	c := New("M", "S", WithBaseURL(srv.URL))
	if _, err := c.GetOrder(context.Background(), "X"); err == nil || !strings.Contains(err.Error(), "GET /v1/order/detail/X") {
		t.Errorf("err = %v", err)
	}
}

func TestCallReadError(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		io.WriteString(w, `{"data"`)
	})
	if _, err := client.GetOrder(context.Background(), "X"); err == nil || !strings.Contains(err.Error(), "read response") {
		t.Errorf("err = %v", err)
	}
}

func TestCallDecodeError(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `not json`)
	})
	if _, err := client.GetOrder(context.Background(), "X"); err == nil || !strings.Contains(err.Error(), "decode response") {
		t.Errorf("err = %v", err)
	}
}

func TestListOrdersError(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	list, err := client.ListOrders(context.Background(), nil)
	if list != nil || !IsRateLimited(err) {
		t.Errorf("list = %v, err = %v", list, err)
	}
}

func TestAPIErrorMessage(t *testing.T) {
	err := &APIError{StatusCode: 422, Message: "invalid", Errors: map[string]string{"b": "y", "a": "x"}}
	if got, want := err.Error(), "paymentgateway: 422: invalid (a: x; b: y)"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	err = &APIError{StatusCode: 404, Message: "missing"}
	if got, want := err.Error(), "paymentgateway: 404: missing"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestParseAPIError(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		raw     string
		message string
		errors  map[string]string
	}{
		{"plain text", 500, " oops \n", "oops", nil},
		{"empty", 502, "", "Bad Gateway", nil},
		{"error string", 401, `{"error":"Unauthorized access"}`, "Unauthorized access", nil},
		{"error list", 400, `{"error":["a","b"]}`, "a; b", nil},
		{"error object", 400, `{"error":{"code":1}}`, `{"code":1}`, nil},
		{"error null", 403, `{"error":null}`, "Forbidden", nil},
		{"no message", 404, `{}`, "Not Found", nil},
		{"field errors", 422, `{"message":"m","errors":{"a":"x","b":["y","z"],"c":null}}`, "m",
			map[string]string{"a": "x", "b": "y; z", "c": ""}},
	}
	for _, tt := range tests {
		e := parseAPIError(tt.status, []byte(tt.raw))
		if e.StatusCode != tt.status || e.Message != tt.message {
			t.Errorf("%s: got %+v", tt.name, e)
		}
		if len(e.Errors) != len(tt.errors) {
			t.Errorf("%s: errors = %v", tt.name, e.Errors)
		}
		for k, v := range tt.errors {
			if e.Errors[k] != v {
				t.Errorf("%s: errors[%s] = %q, want %q", tt.name, k, e.Errors[k], v)
			}
		}
	}
}

func TestIsHelpersNonAPIError(t *testing.T) {
	err := errors.New("plain")
	if IsNotFound(err) || IsValidationError(err) || IsRateLimited(err) || IsNotFound(nil) {
		t.Error("helpers matched non-API error")
	}
	if !IsValidationError(&APIError{StatusCode: http.StatusBadRequest}) {
		t.Error("400 not a validation error")
	}
}

func TestNumberUnmarshal(t *testing.T) {
	tests := []struct {
		in   string
		want Number
		err  bool
	}{
		{`"50000.00"`, 50000, false},
		{`12.5`, 12.5, false},
		{`null`, 7, false},
		{`""`, 7, false},
		{`"abc"`, 7, true},
	}
	for _, tt := range tests {
		n := Number(7)
		err := json.Unmarshal([]byte(tt.in), &n)
		if (err != nil) != tt.err || n != tt.want {
			t.Errorf("%s: got %v, %v", tt.in, n, err)
		}
	}
}

func TestParseTimeAndFormatDate(t *testing.T) {
	at, err := ParseTime("2025-09-05 13:31:42")
	if err != nil || at.UTC().Hour() != 6 {
		t.Errorf("ParseTime = %v, %v", at, err)
	}
	if _, err := ParseTime("bad"); err == nil {
		t.Error("want error")
	}
	if got := formatDate(at); got != "2025-09-05" {
		t.Errorf("formatDate = %q", got)
	}
}

func TestIPNHandlerBodyTooLarge(t *testing.T) {
	called := false
	h := IPNHandler("", func(context.Context, *IPN) error { called = true; return nil })
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/ipn", strings.NewReader(strings.Repeat("a", maxBodySize+1))))
	if rec.Code != http.StatusBadRequest || called {
		t.Errorf("status = %d, called = %v", rec.Code, called)
	}
}

func TestIPNHandlerAllowHeader(t *testing.T) {
	h := IPNHandler("", func(context.Context, *IPN) error { return nil })
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ipn", nil))
	if rec.Header().Get("Allow") != http.MethodPost || rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("headers = %v", rec.Header())
	}
}

func TestCheckoutFormValues(t *testing.T) {
	f := &CheckoutForm{Fields: []Field{{"a", "1"}, {"b", "2"}}}
	if got := f.Values().Encode(); got != "a=1&b=2" {
		t.Errorf("Values = %s", got)
	}
}
