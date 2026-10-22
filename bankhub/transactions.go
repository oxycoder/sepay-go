package bankhub

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// Transaction is a transaction on a linked bank account.
type Transaction struct {
	TransactionID      string       `json:"transaction_id"`
	TransactionDate    string       `json:"transaction_date"`
	BankAccountXID     string       `json:"bank_account_xid"`
	CompanyXID         string       `json:"company_xid"`
	BrandName          string       `json:"brand_name"`
	AccountNumber      string       `json:"account_number"`
	VA                 string       `json:"va"`     // virtual account number, if any
	VAXID              string       `json:"va_xid"` // virtual account id, if any
	ReferenceNumber    string       `json:"reference_number"`
	TransferType       TransferType `json:"transfer_type"`
	Amount             Number       `json:"amount"` // VND
	TransactionContent string       `json:"transaction_content"`
	PaymentCode        string       `json:"payment_code"` // extracted from the content
	WebhookSuccess     Flag         `json:"webhook_success"`
}

// ListTransactionsParams filters ListTransactions. Zero values are ignored.
// Dates are compared in Vietnam time.
type ListTransactionsParams struct {
	PageParams
	CompanyXID     string
	BankAccountXID string
	BrandName      string    // e.g. "MBBank", "Vietcombank"
	Date           time.Time // a single day
	StartDate      time.Time
	EndDate        time.Time
	TransferType   TransferType
	// WebhookSuccess filters by webhook delivery result when non-nil.
	WebhookSuccess *bool
	Query          string // searches content, payment code and account number
}

// ListTransactions returns a page of transactions, newest first. params may be nil.
func (c *Client) ListTransactions(ctx context.Context, params *ListTransactionsParams) (*List[Transaction], error) {
	q := url.Values{}
	if p := params; p != nil {
		pagination(q, p.Page, p.PerPage)
		setIf(q, "company_xid", p.CompanyXID)
		setIf(q, "bank_account_xid", p.BankAccountXID)
		setIf(q, "brand_name", p.BrandName)
		setIf(q, "transaction_date", formatDate(p.Date))
		setIf(q, "start_transaction_date", formatDate(p.StartDate))
		setIf(q, "end_transaction_date", formatDate(p.EndDate))
		setIf(q, "transfer_type", string(p.TransferType))
		if p.WebhookSuccess != nil {
			q.Set("webhook_success", map[bool]string{true: "1", false: "0"}[*p.WebhookSuccess])
		}
		setIf(q, "q", p.Query)
	}
	var out List[Transaction]
	if err := c.list(ctx, "/v1/transaction", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTransaction returns a transaction by its transaction_id.
func (c *Client) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	var out Transaction
	if err := c.get(ctx, http.MethodGet, "/v1/transaction/"+url.PathEscape(transactionID), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SimulateTransactionRequest is the body of SimulateTransaction.
type SimulateTransactionRequest struct {
	BankAccountXID     string       `json:"bank_account_xid"`
	TransferType       TransferType `json:"transfer_type"`
	Amount             float64      `json:"amount"` // VND
	TransactionContent string       `json:"transaction_content"`
}

// SimulatedTransaction is the result of SimulateTransaction.
type SimulatedTransaction struct {
	TransactionID   string       `json:"transaction_id"`
	ReferenceNumber string       `json:"reference_number"`
	TransactionDate string       `json:"transaction_date"`
	Amount          Number       `json:"amount"`
	TransferType    TransferType `json:"transfer_type"`
	Accumulated     Number       `json:"accumulated"`
}

// SimulateTransaction creates a fake transaction to test your IPN handler.
// It only works in the sandbox; production returns 403.
func (c *Client) SimulateTransaction(ctx context.Context, req *SimulateTransactionRequest) (*SimulatedTransaction, error) {
	var out SimulatedTransaction
	if err := c.get(ctx, http.MethodPost, "/v1/transaction/create", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
