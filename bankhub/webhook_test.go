package bankhub

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const linkedEvent = `{
  "timestamp": 1768312023,
  "xid": "e343bdec-97dc-4c60-b953-f81b183e0096",
  "version": "1.0",
  "event": "BANK_ACCOUNT_LINKED",
  "metadata": {
    "bank_account_xid": "587eabfb-f086-11f0-b16e-52c7e9b4f41b",
    "brand_name": "MBBank",
    "account_type": "personal",
    "account_number": "0000000000011111",
    "link_token_xid": "ef4cb765-c297-4a51-a9eb-6f7a751978ae",
    "link_session_xid": null
  }
}`

const ipnBody = `{
  "gateway": "MBBank",
  "transaction_date": "2024-01-17 10:30:00",
  "account_number": "0123",
  "bank_account_xid": "b-1",
  "va": null,
  "payment_code": "DH123",
  "content": "Thanh toan DH123",
  "transfer_type": "credit",
  "amount": 500000,
  "reference_code": "FT123",
  "accumulated": 0,
  "transaction_id": "t-1"
}`

func TestParseEvent(t *testing.T) {
	e, err := ParseEvent([]byte(linkedEvent))
	if err != nil {
		t.Fatal(err)
	}
	if e.Event != EventBankAccountLinked || e.Time().Unix() != 1768312023 {
		t.Errorf("event = %+v", e)
	}
	m := e.Metadata
	if m.BankAccountXID != "587eabfb-f086-11f0-b16e-52c7e9b4f41b" || m.AccountType != AccountPersonal || m.LinkSessionXID != "" {
		t.Errorf("metadata = %+v", m)
	}
}

func post(h http.Handler, header, value, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/hook", strings.NewReader(body))
	if header != "" {
		req.Header.Set(header, value)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestEventHandler(t *testing.T) {
	var got *Event
	h := EventHandler("s3cret", func(_ context.Context, e *Event) error {
		got = e
		return nil
	})

	if rec := post(h, "X-Secret-Key", "wrong", linkedEvent); rec.Code != http.StatusUnauthorized || got != nil {
		t.Fatalf("wrong secret: %d", rec.Code)
	}
	rec := post(h, "X-Secret-Key", "s3cret", linkedEvent)
	if rec.Code != http.StatusOK || rec.Body.String() != `{"success":true}` {
		t.Fatalf("got %d %s", rec.Code, rec.Body)
	}
	if got == nil || got.Metadata.BrandName != "MBBank" {
		t.Errorf("event = %+v", got)
	}
	if rec := post(h, "X-Secret-Key", "s3cret", "not json"); rec.Code != http.StatusBadRequest {
		t.Errorf("bad body: %d", rec.Code)
	}
}

func TestIPNHandler(t *testing.T) {
	var got *IPN
	fail := false
	h := IPNHandler("key", func(_ context.Context, n *IPN) error {
		got = n
		if fail {
			return errors.New("db down")
		}
		return nil
	})

	if rec := post(h, "Authorization", "Apikey nope", ipnBody); rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong key: %d", rec.Code)
	}
	rec := post(h, "Authorization", "Apikey key", ipnBody)
	if rec.Code != http.StatusOK || rec.Body.String() != `{"success":true}` || rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("got %d %s", rec.Code, rec.Body)
	}
	if got.TransactionID != "t-1" || got.Amount != 500000 || got.TransferType != TransferCredit || got.VA != "" {
		t.Errorf("ipn = %+v", got)
	}

	fail = true
	if rec := post(h, "Authorization", "Apikey key", ipnBody); rec.Code != http.StatusInternalServerError || rec.Body.String() != `{"success":false}` {
		t.Errorf("handler error: %d %s", rec.Code, rec.Body)
	}

	req := httptest.NewRequest(http.MethodGet, "/hook", nil)
	req.Header.Set("Authorization", "Apikey key")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET: %d", rr.Code)
	}
}
