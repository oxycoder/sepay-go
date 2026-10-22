package paymentgateway

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// docIPN is the card IPN example from the SePay docs.
const docIPN = `{
  "timestamp": 1757058220,
  "notification_type": "ORDER_PAID",
  "order": {
    "id": "e2c195be-c721-47eb-b323-99ab24e52d85",
    "order_id": "NPSETVI00101000042R",
    "order_status": "CAPTURED",
    "order_currency": "VND",
    "order_amount": "50000.00",
    "order_invoice_number": "SUB_202509_001",
    "custom_data": [],
    "user_agent": "Mozilla/5.0",
    "ip_address": "14.xxx.xxx.xxx",
    "order_description": "Thanh toán định kỳ gói Premium tháng 9/2025"
  },
  "transaction": {
    "id": "384c66dd-41e6-4316-a544-b4141682595c",
    "payment_method": "CARD",
    "transaction_id": "68ba94ac80123",
    "transaction_type": "PAYMENT",
    "transaction_date": "2025-09-01 00:00:15",
    "transaction_status": "APPROVED",
    "transaction_amount": "50000",
    "transaction_currency": "VND",
    "authentication_status": "AUTHENTICATION_SUCCESSFUL",
    "card_number": "4111XXXXXXXX1111",
    "card_holder_name": "NGUYEN VAN A",
    "card_expiry": "12/26",
    "card_funding_method": "CREDIT",
    "card_brand": "VISA"
  },
  "customer": {
    "id": "bae12d2f-0580-4669-8841-cc35cf671613",
    "customer_id": "CUST_001"
  }
}`

func TestParseIPN(t *testing.T) {
	n, err := ParseIPN([]byte(docIPN))
	if err != nil {
		t.Fatal(err)
	}
	if n.NotificationType != NotificationOrderPaid || n.Timestamp != 1757058220 {
		t.Errorf("ipn = %+v", n)
	}
	if n.Order.OrderInvoiceNumber != "SUB_202509_001" || n.Order.OrderStatus != OrderCaptured || n.Order.OrderAmount != 50000 {
		t.Errorf("order = %+v", n.Order)
	}
	if tx := n.Transaction; tx.TransactionID != "68ba94ac80123" || tx.PaymentMethod != PaymentCard || tx.CardBrand != CardVisa || tx.TransactionAmount != 50000 {
		t.Errorf("transaction = %+v", tx)
	}
	if n.Customer == nil || n.Customer.CustomerID != "CUST_001" {
		t.Errorf("customer = %+v", n.Customer)
	}
}

func TestParseIPNBankTransferNulls(t *testing.T) {
	body := `{"timestamp":1759134682,"notification_type":"ORDER_PAID","order":{"order_invoice_number":"INV-1","order_amount":"100000.00"},
	"transaction":{"payment_method":"BANK_TRANSFER","card_number":null,"card_brand":null},"customer":null,"agreement":null}`
	n, err := ParseIPN([]byte(body))
	if err != nil {
		t.Fatal(err)
	}
	if n.Customer != nil || n.Transaction.CardNumber != "" || n.Transaction.PaymentMethod != PaymentBankTransfer {
		t.Errorf("ipn = %+v", n)
	}
}

func TestIPNHandler(t *testing.T) {
	var got *IPN
	h := IPNHandler("sk", func(ctx context.Context, n *IPN) error {
		got = n
		if n.Order.OrderInvoiceNumber == "fail" {
			return errors.New("db down")
		}
		return nil
	})
	tests := []struct {
		name, method, key, body string
		status                  int
	}{
		{"ok", http.MethodPost, "sk", docIPN, http.StatusOK},
		{"wrong key", http.MethodPost, "nope", docIPN, http.StatusUnauthorized},
		{"no key", http.MethodPost, "", docIPN, http.StatusUnauthorized},
		{"bad json", http.MethodPost, "sk", "{", http.StatusBadRequest},
		{"fn error", http.MethodPost, "sk", `{"order":{"order_invoice_number":"fail"}}`, http.StatusInternalServerError},
		{"GET", http.MethodGet, "sk", "", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		got = nil
		req := httptest.NewRequest(tt.method, "/ipn", strings.NewReader(tt.body))
		if tt.key != "" {
			req.Header.Set("X-Secret-Key", tt.key)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != tt.status {
			t.Errorf("%s: status = %d, want %d", tt.name, rec.Code, tt.status)
		}
		wantBody := `{"success":` + map[bool]string{true: "true", false: "false"}[tt.status == 200] + `}`
		if rec.Body.String() != wantBody {
			t.Errorf("%s: body = %s", tt.name, rec.Body)
		}
		if tt.status == http.StatusOK && (got == nil || got.Order.OrderID != "NPSETVI00101000042R") {
			t.Errorf("%s: fn got %+v", tt.name, got)
		}
		if tt.status == http.StatusUnauthorized && got != nil {
			t.Errorf("%s: fn called despite bad key", tt.name)
		}
	}
}

func TestIPNHandlerWithoutKey(t *testing.T) {
	h := IPNHandler("", func(context.Context, *IPN) error { return nil })
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/ipn", strings.NewReader(docIPN)))
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d", rec.Code)
	}
}
