package paymentgateway

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestServer returns a client wired to a fake API that checks Basic auth
// before calling handler.
func newTestServer(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, secret, ok := r.BasicAuth()
		if !ok || id != "MERCHANT_123" || secret != "SECRET_KEY" {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, `{"message":"Unauthenticated"}`)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	return New("MERCHANT_123", "SECRET_KEY", WithBaseURL(srv.URL))
}

func TestListOrders(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/order" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		want := "customer_id=CUST_001&end_created_at=2024-01-01&from_created_at=2023-12-01&order_status=CAPTURED&page=2&per_page=50&q=INV&sort=created_at%3Adesc"
		if got := r.URL.RawQuery; got != want {
			t.Errorf("query = %s\nwant    %s", got, want)
		}
		io.WriteString(w, `{"data":[{"id":"9","customer_id":"2427","order_id":"SEPAY-68BA83CE637C1","order_invoice_number":"DH1757053857","order_status":"AUTHENTICATION_NOT_NEEDED","order_amount":"300000.00","order_currency":"VND","order_description":"Đơn hàng","authentication_status":null,"created_at":"2025-09-05 13:31:42","updated_at":"2025-09-05 13:31:42"}],"meta":{"per_page":50,"total":51,"has_more":false,"current_page":2,"page_count":2}}`)
	})
	list, err := client.ListOrders(context.Background(), &ListOrdersParams{
		PageParams:  PageParams{Page: 2, PerPage: 50},
		Query:       "INV",
		OrderStatus: OrderCaptured,
		CustomerID:  "CUST_001",
		FromDate:    time.Date(2023, 12, 1, 0, 0, 0, 0, vnZone),
		ToDate:      time.Date(2023, 12, 31, 20, 0, 0, 0, time.UTC), // 2024-01-01 in Vietnam
		Sort:        SortCreatedDesc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 1 || list.Meta.Total != 51 || list.Meta.CurrentPage != 2 {
		t.Fatalf("list = %+v", list)
	}
	o := list.Data[0]
	if o.OrderID != "SEPAY-68BA83CE637C1" || o.OrderStatus != OrderAuthenticationNotNeeded || o.OrderAmount != 300000 {
		t.Errorf("order = %+v", o)
	}
	created, err := ParseTime(o.CreatedAt)
	if err != nil || !created.Equal(time.Date(2025, 9, 5, 6, 31, 42, 0, time.UTC)) {
		t.Errorf("ParseTime = %v, %v", created, err)
	}
}

func TestListOrdersNilParams(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("query = %q", r.URL.RawQuery)
		}
		io.WriteString(w, `{"data":[],"meta":{}}`)
	})
	if _, err := client.ListOrders(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
}

func TestGetOrder(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/order/detail/SEPAY-68B01673A77FF" {
			t.Errorf("path = %s", r.URL.Path)
		}
		io.WriteString(w, `{"data":{"id":"1","customer_id":null,"order_id":"SEPAY-68B01673A77FF","order_invoice_number":"DH1756370479","order_status":"CAPTURED","order_amount":"300000.00","order_currency":"VND","transactions":[{"id":"1","payment_method":"CARD","transaction_type":"PAYMENT","transaction_amount":"300000","transaction_currency":"VND","transaction_status":"APPROVED","card_number":"512345xxxxxx0008","card_funding_method":"DEBIT","card_brand":"MASTERCARD","transaction_date":"2025-08-28 15:43:41"}]}}`)
	})
	o, err := client.GetOrder(context.Background(), "SEPAY-68B01673A77FF")
	if err != nil {
		t.Fatal(err)
	}
	if o.OrderStatus != OrderCaptured || o.CustomerID != "" || len(o.Transactions) != 1 {
		t.Fatalf("order = %+v", o)
	}
	tx := o.Transactions[0]
	if tx.PaymentMethod != PaymentCard || tx.TransactionStatus != TransactionApproved || tx.TransactionAmount != 300000 || tx.CardBrand != CardMastercard {
		t.Errorf("transaction = %+v", tx)
	}
}

func TestCancelAndVoid(t *testing.T) {
	for path, call := range map[string]func(*Client) error{
		"/v1/order/cancel":          func(c *Client) error { return c.CancelOrder(context.Background(), "DH1") },
		"/v1/order/voidTransaction": func(c *Client) error { return c.VoidTransaction(context.Background(), "DH1") },
	} {
		client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost || r.URL.Path != path {
				t.Errorf("%s %s, want POST %s", r.Method, r.URL.Path, path)
			}
			if ct := r.Header.Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q", ct)
			}
			var body map[string]string
			json.NewDecoder(r.Body).Decode(&body)
			if body["order_invoice_number"] != "DH1" {
				t.Errorf("body = %v", body)
			}
			io.WriteString(w, `{"message":"Đã hủy thành công"}`)
		})
		if err := call(client); err != nil {
			t.Errorf("%s: %v", path, err)
		}
	}
}

func TestAPIError(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		io.WriteString(w, `{"message":"The given data was invalid.","errors":{"order_invoice_number":["Đơn hàng đã được thanh toán"]}}`)
	})
	err := client.CancelOrder(context.Background(), "DH1")
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %v", err)
	}
	if apiErr.StatusCode != 422 || apiErr.Errors["order_invoice_number"] != "Đơn hàng đã được thanh toán" {
		t.Errorf("apiErr = %+v", apiErr)
	}
	if !IsValidationError(err) || IsNotFound(err) || IsRateLimited(err) {
		t.Error("status helpers wrong")
	}
}

func TestAPIErrorWrongCredentials(t *testing.T) {
	client := newTestServer(t, nil)
	client.secretKey = "wrong"
	_, err := client.GetOrder(context.Background(), "x")
	if !hasStatus(err, http.StatusUnauthorized) {
		t.Errorf("err = %v", err)
	}
}

func TestAPIErrorNonJSON(t *testing.T) {
	client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	_, err := client.GetOrder(context.Background(), "missing")
	if !IsNotFound(err) || err.Error() != "paymentgateway: 404: Not Found" {
		t.Errorf("err = %v", err)
	}
}
