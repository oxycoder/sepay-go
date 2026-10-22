package webhook

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name, contentType, body string
	}{
		{"bad content type", "application/json; =", payload},
		{"bad multipart", "multipart/form-data; boundary=xyz", "not multipart"},
		{"bad form escape", ContentTypeForm, "id=%zz"},
	}
	for _, tt := range tests {
		if _, err := Parse(tt.contentType, []byte(tt.body)); err == nil {
			t.Errorf("%s: want error", tt.name)
		}
	}
}

func TestParseFormMissingNumbers(t *testing.T) {
	tx, err := Parse(ContentTypeForm, []byte("gateway=MBBank&code=&transferType=out"))
	if err != nil {
		t.Fatal(err)
	}
	if tx.ID != 0 || tx.TransferAmount != 0 || tx.Accumulated != 0 || tx.Gateway != "MBBank" || tx.TransferType != TransferOut {
		t.Errorf("tx = %+v", tx)
	}
}

func TestVerifySignatureBadTimestamp(t *testing.T) {
	h := http.Header{}
	h.Set(HeaderSignature, Sign("s3cret", 1, []byte(payload)))
	h.Set(HeaderTimestamp, "abc")
	if err := VerifySignature(h, []byte(payload), "s3cret", DefaultTolerance); !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("got %v, want %v", err, ErrInvalidSignature)
	}
}

func TestVerifySignatureFutureTimestamp(t *testing.T) {
	ts := time.Now().Add(10 * time.Minute).Unix()
	h := http.Header{}
	h.Set(HeaderSignature, Sign("s3cret", ts, []byte(payload)))
	h.Set(HeaderTimestamp, strconv.FormatInt(ts, 10))
	if err := VerifySignature(h, []byte(payload), "s3cret", DefaultTolerance); !errors.Is(err, ErrTimestampExpired) {
		t.Errorf("got %v, want %v", err, ErrTimestampExpired)
	}
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		auth  string
		token string
		ok    bool
	}{
		{"Bearer abc", "abc", true},
		{"bearer abc", "abc", true},
		{"Bearer ", "", false},
		{"Bear", "", false},
		{"Apikey abc", "", false},
	}
	for _, tt := range tests {
		h := http.Header{"Authorization": {tt.auth}}
		if token, ok := bearerToken(h); token != tt.token || ok != tt.ok {
			t.Errorf("%q: got %q, %v", tt.auth, token, ok)
		}
	}
}

func TestHandlerWithTolerance(t *testing.T) {
	fn := func(context.Context, *Transaction) error { return nil }
	old := time.Now().Add(-time.Hour).Unix()
	h := http.Header{}
	h.Set(HeaderSignature, Sign("s3cret", old, []byte(payload)))
	h.Set(HeaderTimestamp, strconv.FormatInt(old, 10))

	if rec := serve(Handler(fn, WithHMAC("s3cret")), http.MethodPost, h, payload); rec.Code != http.StatusUnauthorized {
		t.Errorf("default tolerance: status %d", rec.Code)
	}
	if rec := serve(Handler(fn, WithHMAC("s3cret"), WithTolerance(2*time.Hour)), http.MethodPost, h, payload); rec.Code != http.StatusOK {
		t.Errorf("2h tolerance: status %d", rec.Code)
	}
	if rec := serve(Handler(fn, WithHMAC("s3cret"), WithTolerance(0)), http.MethodPost, h, payload); rec.Code != http.StatusOK {
		t.Errorf("tolerance disabled: status %d", rec.Code)
	}
}

func TestHandlerBodyTooLarge(t *testing.T) {
	called := false
	fn := func(context.Context, *Transaction) error { called = true; return nil }
	body := strings.Repeat("a", maxBodySize+1)
	rec := serve(Handler(fn), http.MethodPost, nil, body)
	if rec.Code != http.StatusBadRequest || called {
		t.Errorf("status %d, called %v", rec.Code, called)
	}
	if rec.Header().Get("Allow") != "" {
		t.Error("Allow header set on POST")
	}
}

func TestHandlerAllowHeader(t *testing.T) {
	rec := serve(Handler(func(context.Context, *Transaction) error { return nil }), http.MethodPut, nil, "")
	if rec.Code != http.StatusMethodNotAllowed || rec.Header().Get("Allow") != http.MethodPost {
		t.Errorf("status %d, Allow %q", rec.Code, rec.Header().Get("Allow"))
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type %q", ct)
	}
}

func TestTimeInvalid(t *testing.T) {
	tx := &Transaction{TransactionDate: "bad"}
	if _, err := tx.Time(); err == nil {
		t.Error("want error")
	}
}
