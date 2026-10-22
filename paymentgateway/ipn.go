package paymentgateway

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// maxBodySize caps IPN request bodies.
const maxBodySize = 1 << 20

// NotificationType is the kind of IPN.
type NotificationType string

const (
	NotificationOrderPaid       NotificationType = "ORDER_PAID"       // payment succeeded
	NotificationTransactionVoid NotificationType = "TRANSACTION_VOID" // card payment voided
)

// IPN is the payment notification SePay POSTs to your IPN URL.
type IPN struct {
	Timestamp        int64            `json:"timestamp"` // Unix seconds
	NotificationType NotificationType `json:"notification_type"`
	Order            IPNOrder         `json:"order"`
	Transaction      IPNTransaction   `json:"transaction"`
	Customer         *IPNCustomer     `json:"customer"`  // nil if the checkout had no customer
	Agreement        json.RawMessage  `json:"agreement"` // recurring payments; not released yet
}

// IPNOrder is the order an IPN is about.
type IPNOrder struct {
	ID                 string          `json:"id"`       // SePay internal id (UUID)
	OrderID            string          `json:"order_id"` // pass to GetOrder
	OrderStatus        OrderStatus     `json:"order_status"`
	OrderCurrency      Currency        `json:"order_currency"`
	OrderAmount        Number          `json:"order_amount"`
	OrderInvoiceNumber string          `json:"order_invoice_number"` // CheckoutRequest.OrderInvoiceNumber
	OrderDescription   string          `json:"order_description"`
	CustomData         json.RawMessage `json:"custom_data"`
	UserAgent          string          `json:"user_agent"` // customer's browser
	IPAddress          string          `json:"ip_address"` // customer's IP
}

// IPNTransaction is the transaction an IPN is about. Card fields are only set
// for PaymentCard.
type IPNTransaction struct {
	ID                   string            `json:"id"`             // SePay internal id (UUID)
	TransactionID        string            `json:"transaction_id"` // unique per transaction
	PaymentMethod        PaymentMethod     `json:"payment_method"`
	TransactionType      TransactionType   `json:"transaction_type"`
	TransactionDate      string            `json:"transaction_date"` // parse with ParseTime
	TransactionStatus    TransactionStatus `json:"transaction_status"`
	TransactionAmount    Number            `json:"transaction_amount"`
	TransactionCurrency  Currency          `json:"transaction_currency"`
	AuthenticationStatus string            `json:"authentication_status"`
	CardNumber           string            `json:"card_number"` // masked
	CardHolderName       string            `json:"card_holder_name"`
	CardExpiry           string            `json:"card_expiry"`
	CardFundingMethod    CardFundingMethod `json:"card_funding_method"`
	CardBrand            CardBrand         `json:"card_brand"`
}

// IPNCustomer is the customer of an IPN.
type IPNCustomer struct {
	ID         string `json:"id"`          // SePay internal id (UUID)
	CustomerID string `json:"customer_id"` // CheckoutRequest.CustomerID
}

// ParseIPN decodes an IPN body.
func ParseIPN(body []byte) (*IPN, error) {
	var n IPN
	if err := json.Unmarshal(body, &n); err != nil {
		return nil, fmt.Errorf("paymentgateway: decode IPN: %w", err)
	}
	return &n, nil
}

// VerifyIPN reports whether the X-Secret-Key header of an IPN request equals
// secretKey, using a constant-time compare.
func VerifyIPN(h http.Header, secretKey string) bool {
	return subtle.ConstantTimeCompare([]byte(h.Get("X-Secret-Key")), []byte(secretKey)) == 1
}

// IPNHandler returns an http.Handler for the IPN URL. It checks the
// X-Secret-Key header against secretKey (skipped if secretKey is empty, i.e.
// the IPN auth type is not SECRET_KEY), decodes the IPN and calls fn. It
// answers {"success": true} when fn returns nil, 401 for a bad key, 400 for
// an unreadable body and 500 when fn returns an error.
//
// SePay may deliver the same IPN more than once; make fn idempotent, e.g. on
// IPN.Transaction.TransactionID or IPN.Order.OrderInvoiceNumber.
func IPNHandler(secretKey string, fn func(context.Context, *IPN) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			respond(w, http.StatusMethodNotAllowed)
			return
		}
		if secretKey != "" && !VerifyIPN(r.Header, secretKey) {
			respond(w, http.StatusUnauthorized)
			return
		}
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodySize))
		if err != nil {
			respond(w, http.StatusBadRequest)
			return
		}
		n, err := ParseIPN(body)
		if err != nil {
			respond(w, http.StatusBadRequest)
			return
		}
		if err := fn(r.Context(), n); err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		respond(w, http.StatusOK)
	})
}

// respond writes {"success": status == 200} with the given status.
func respond(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"success":%t}`, status == http.StatusOK)
}
