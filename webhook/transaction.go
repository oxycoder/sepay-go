package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/url"
	"strconv"
	"time"
)

// TransferType is the direction of a transaction.
type TransferType string

const (
	TransferIn  TransferType = "in"  // money in
	TransferOut TransferType = "out" // money out
)

// Content types SePay can send, set per webhook.
const (
	ContentTypeJSON      = "application/json" // default
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeMultipart = "multipart/form-data"
)

// DateLayout is the format of Transaction.TransactionDate, in Vietnam time.
const DateLayout = "2006-01-02 15:04:05"

var vnZone = time.FixedZone("ICT", 7*60*60)

// Transaction is the payload of a webhook delivery.
type Transaction struct {
	// ID is the transaction ID on SePay. It stays the same across retries and
	// manual replays, so use it as the unique key to drop duplicates.
	ID              int64  `json:"id"`
	Gateway         string `json:"gateway"`         // bank name, e.g. "Vietcombank"
	TransactionDate string `json:"transactionDate"` // DateLayout, Vietnam time; see Time
	AccountNumber   string `json:"accountNumber"`
	SubAccount      string `json:"subAccount"` // matched virtual account (VA), empty if none
	// Code is the payment code extracted from Content using the prefixes set
	// in Company → Configuration. Empty when no payment code matched (null in
	// the JSON payload).
	Code           string       `json:"code"`
	Content        string       `json:"content"` // transfer content as sent by the bank
	TransferType   TransferType `json:"transferType"`
	Description    string       `json:"description"`    // may be empty for some banks
	TransferAmount int64        `json:"transferAmount"` // VND, always positive
	Accumulated    int64        `json:"accumulated"`    // balance after the transaction; 0 if the bank does not report it
	ReferenceCode  string       `json:"referenceCode"`  // bank reference, may be empty

	// Raw is the request body exactly as received. Store it before processing
	// for debugging, audits and reconciliation.
	Raw []byte `json:"-"`
}

// Time parses TransactionDate in Vietnam time.
func (t *Transaction) Time() (time.Time, error) {
	return time.ParseInLocation(DateLayout, t.TransactionDate, vnZone)
}

// Parse decodes a webhook body. contentType is the request's Content-Type
// header; an empty value is treated as JSON.
func Parse(contentType string, body []byte) (*Transaction, error) {
	mediaType := ContentTypeJSON
	var params map[string]string
	if contentType != "" {
		var err error
		if mediaType, params, err = mime.ParseMediaType(contentType); err != nil {
			return nil, fmt.Errorf("webhook: content type: %w", err)
		}
	}

	var (
		tx  *Transaction
		err error
	)
	switch mediaType {
	case ContentTypeJSON:
		tx = new(Transaction)
		err = json.Unmarshal(body, tx)
	case ContentTypeForm:
		var form url.Values
		if form, err = url.ParseQuery(string(body)); err == nil {
			tx, err = fromForm(form)
		}
	case ContentTypeMultipart:
		var form *multipart.Form
		form, err = multipart.NewReader(bytes.NewReader(body), params["boundary"]).ReadForm(int64(len(body)))
		if err == nil {
			tx, err = fromForm(form.Value)
		}
	default:
		return nil, fmt.Errorf("webhook: unsupported content type %q", mediaType)
	}
	if err != nil {
		return nil, fmt.Errorf("webhook: decode payload: %w", err)
	}
	tx.Raw = body
	return tx, nil
}

func fromForm(v url.Values) (*Transaction, error) {
	tx := &Transaction{
		Gateway:         v.Get("gateway"),
		TransactionDate: v.Get("transactionDate"),
		AccountNumber:   v.Get("accountNumber"),
		SubAccount:      v.Get("subAccount"),
		Code:            v.Get("code"),
		Content:         v.Get("content"),
		TransferType:    TransferType(v.Get("transferType")),
		Description:     v.Get("description"),
		ReferenceCode:   v.Get("referenceCode"),
	}
	for _, f := range []struct {
		name string
		dst  *int64
	}{
		{"id", &tx.ID},
		{"transferAmount", &tx.TransferAmount},
		{"accumulated", &tx.Accumulated},
	} {
		s := v.Get(f.name)
		if s == "" {
			continue
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", f.name, err)
		}
		*f.dst = n
	}
	return tx, nil
}
