package webhook

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

const payload = `{
  "id": 92704,
  "gateway": "Vietcombank",
  "transactionDate": "2024-07-02 11:08:33",
  "accountNumber": "1017588888",
  "subAccount": "",
  "code": "SEVN63DC8E5C",
  "content": "SEVN63DC8E5C chuyen tien",
  "transferType": "in",
  "description": "NGUYEN VAN A chuyen tien",
  "transferAmount": 5000000,
  "accumulated": 105000000,
  "referenceCode": "FT24012345678"
}`

func checkTx(t *testing.T, tx *Transaction) {
	t.Helper()
	if tx == nil {
		t.Fatal("no transaction")
	}
	if tx.ID != 92704 || tx.Gateway != "Vietcombank" || tx.Code != "SEVN63DC8E5C" ||
		tx.TransferType != TransferIn || tx.TransferAmount != 5000000 || tx.Accumulated != 105000000 ||
		tx.ReferenceCode != "FT24012345678" {
		t.Errorf("tx = %+v", tx)
	}
}

func TestParseJSON(t *testing.T) {
	tx, err := Parse("application/json; charset=utf-8", []byte(payload))
	if err != nil {
		t.Fatal(err)
	}
	checkTx(t, tx)
	if string(tx.Raw) != payload {
		t.Error("Raw not set")
	}
	at, err := tx.Time()
	if err != nil || !at.Equal(time.Date(2024, 7, 2, 4, 8, 33, 0, time.UTC)) {
		t.Errorf("Time() = %v, %v", at, err)
	}

	tx, err = Parse("", []byte(`{"id":1,"code":null}`))
	if err != nil || tx.Code != "" {
		t.Errorf("null code: %+v, %v", tx, err)
	}
}

var formValues = map[string]string{
	"id": "92704", "gateway": "Vietcombank", "code": "SEVN63DC8E5C", "transferType": "in",
	"transferAmount": "5000000", "accumulated": "105000000", "referenceCode": "FT24012345678",
}

func TestParseForm(t *testing.T) {
	var parts []string
	for k, v := range formValues {
		parts = append(parts, k+"="+v)
	}
	tx, err := Parse(ContentTypeForm, []byte(strings.Join(parts, "&")))
	if err != nil {
		t.Fatal(err)
	}
	checkTx(t, tx)

	if _, err := Parse(ContentTypeForm, []byte("id=abc")); err == nil {
		t.Error("bad id: want error")
	}
}

func TestParseMultipart(t *testing.T) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	for k, v := range formValues {
		mw.WriteField(k, v)
	}
	mw.Close()
	tx, err := Parse(mw.FormDataContentType(), buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	checkTx(t, tx)
}

func TestParseUnsupported(t *testing.T) {
	if _, err := Parse("text/plain", []byte(payload)); err == nil {
		t.Error("want error")
	}
}

func TestVerifySignature(t *testing.T) {
	body := []byte(payload)
	now := time.Now().Unix()
	header := func(sig string, ts int64) http.Header {
		h := http.Header{}
		h.Set(HeaderSignature, sig)
		h.Set(HeaderTimestamp, strconv.FormatInt(ts, 10))
		return h
	}

	tests := []struct {
		name string
		h    http.Header
		body []byte
		want error
	}{
		{"valid", header(Sign("s3cret", now, body), now), body, nil},
		{"wrong secret", header(Sign("other", now, body), now), body, ErrInvalidSignature},
		{"tampered body", header(Sign("s3cret", now, body), now), []byte(`{"id":1}`), ErrInvalidSignature},
		{"timestamp not signed", header(Sign("s3cret", now, body), now+1), body, ErrInvalidSignature},
		{"expired", header(Sign("s3cret", now-600, body), now-600), body, ErrTimestampExpired},
		{"missing", http.Header{}, body, ErrMissingSignature},
		{"bad hex", header("sha256=zz", now), body, ErrInvalidSignature},
	}
	for _, tt := range tests {
		if err := VerifySignature(tt.h, tt.body, "s3cret", DefaultTolerance); !errors.Is(err, tt.want) {
			t.Errorf("%s: got %v, want %v", tt.name, err, tt.want)
		}
	}

	old := header(Sign("s3cret", now-600, body), now-600)
	if err := VerifySignature(old, body, "s3cret", 0); err != nil {
		t.Errorf("tolerance 0: %v", err)
	}
}

func serve(h http.Handler, method string, header http.Header, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, "/hook", strings.NewReader(body))
	for k, v := range header {
		req.Header[k] = v
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHandler(t *testing.T) {
	var got *Transaction
	fn := func(_ context.Context, tx *Transaction) error {
		got = tx
		if tx.ID == 0 {
			return errors.New("boom")
		}
		return nil
	}
	now := time.Now().Unix()
	signed := http.Header{}
	signed.Set(HeaderSignature, Sign("s3cret", now, []byte(payload)))
	signed.Set(HeaderTimestamp, strconv.FormatInt(now, 10))
	signed.Set("Content-Type", ContentTypeJSON)

	tests := []struct {
		name     string
		h        http.Handler
		method   string
		header   http.Header
		body     string
		wantCode int
	}{
		{"hmac ok", Handler(fn, WithHMAC("s3cret")), http.MethodPost, signed, payload, http.StatusOK},
		{"hmac wrong secret", Handler(fn, WithHMAC("other")), http.MethodPost, signed, payload, http.StatusUnauthorized},
		{"apikey ok", Handler(fn, WithAPIKey("k")), http.MethodPost, http.Header{"Authorization": {"Apikey k"}}, payload, http.StatusOK},
		{"apikey wrong", Handler(fn, WithAPIKey("k")), http.MethodPost, http.Header{"Authorization": {"Apikey x"}}, payload, http.StatusUnauthorized},
		{"bearer ok", Handler(fn, WithBearer(func(_ context.Context, tok string) bool { return tok == "t" })), http.MethodPost, http.Header{"Authorization": {"Bearer t"}}, payload, http.StatusOK},
		{"bearer wrong", Handler(fn, WithBearer(func(_ context.Context, tok string) bool { return tok == "t" })), http.MethodPost, http.Header{"Authorization": {"Bearer x"}}, payload, http.StatusUnauthorized},
		{"no auth", Handler(fn), http.MethodPost, nil, payload, http.StatusOK},
		{"bad body", Handler(fn), http.MethodPost, nil, "{", http.StatusBadRequest},
		{"fn error", Handler(fn), http.MethodPost, nil, `{"id":0}`, http.StatusInternalServerError},
		{"GET", Handler(fn), http.MethodGet, nil, "", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		got = nil
		rec := serve(tt.h, tt.method, tt.header, tt.body)
		if rec.Code != tt.wantCode {
			t.Errorf("%s: status %d, want %d", tt.name, rec.Code, tt.wantCode)
		}
		wantBody := `{"success":false}`
		if tt.wantCode == http.StatusOK {
			wantBody = `{"success":true}`
			checkTx(t, got)
		}
		if rec.Body.String() != wantBody {
			t.Errorf("%s: body %s", tt.name, rec.Body)
		}
		if tt.wantCode == http.StatusUnauthorized && got != nil {
			t.Errorf("%s: fn called on unauthorized request", tt.name)
		}
	}
}
