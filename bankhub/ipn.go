package bankhub

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// maxBodySize caps webhook request bodies.
const maxBodySize = 1 << 20

// IPN is a balance change notification (a new transaction on a linked account).
type IPN struct {
	// TransactionID is unique per transaction. SePay may deliver the same
	// notification more than once, so drop ones you have already processed.
	TransactionID   string       `json:"transaction_id"`
	Gateway         string       `json:"gateway"` // bank, e.g. "MBBank"
	TransactionDate string       `json:"transaction_date"`
	AccountNumber   string       `json:"account_number"`
	BankAccountXID  string       `json:"bank_account_xid"`
	VA              string       `json:"va"`           // virtual account, if any
	PaymentCode     string       `json:"payment_code"` // if recognized
	Content         string       `json:"content"`
	TransferType    TransferType `json:"transfer_type"`
	Amount          Number       `json:"amount"`
	ReferenceCode   string       `json:"reference_code"`
	Accumulated     Number       `json:"accumulated"` // not supported yet, always 0
}

// ParseIPN decodes an IPN body.
func ParseIPN(body []byte) (*IPN, error) {
	var n IPN
	if err := json.Unmarshal(body, &n); err != nil {
		return nil, fmt.Errorf("bankhub: decode IPN: %w", err)
	}
	return &n, nil
}

// IPNHandler returns an http.Handler for balance change notifications. It
// checks that the Authorization header is "Apikey {apiKey}" (skipped if apiKey
// is empty, e.g. when you allowlist SePay's IPs instead), decodes the IPN and
// calls fn. It answers {"success": true} when fn returns nil and 500 otherwise.
//
// SePay only retries on network failures and may still deliver duplicates;
// make fn idempotent on IPN.TransactionID.
func IPNHandler(apiKey string, fn func(context.Context, *IPN) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiKey != "" && !secureEqual(r.Header.Get("Authorization"), "Apikey "+apiKey) {
			respond(w, http.StatusUnauthorized)
			return
		}
		body, ok := readBody(w, r)
		if !ok {
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

func readBody(w http.ResponseWriter, r *http.Request) ([]byte, bool) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		respond(w, http.StatusMethodNotAllowed)
		return nil, false
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodySize))
	if err != nil {
		respond(w, http.StatusBadRequest)
		return nil, false
	}
	return body, true
}

// respond writes {"success": status == 200} with the given status.
func respond(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"success":%t}`, status == http.StatusOK)
}

func secureEqual(got, want string) bool {
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}
