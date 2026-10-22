package bankhub

import (
	"context"
	"net/http"
	"net/url"
)

// BankAccount is a bank account linked through Bank Hub.
type BankAccount struct {
	XID               string      `json:"xid"`
	CompanyXID        string      `json:"company_xid"`
	BrandName         string      `json:"brand_name"` // e.g. "OCB", "Vietcombank"
	BankBIN           string      `json:"bank_bin"`   // NAPAS BIN, e.g. "970416"; empty if unknown
	BankCode          string      `json:"bank_code"`  // e.g. "ACB"; empty if unknown
	AccountHolderName string      `json:"account_holder_name"`
	AccountNumber     string      `json:"account_number"`
	AccountType       AccountType `json:"account_type"`
	Accumulated       Number      `json:"accumulated"`
	BankAPIConnected  Flag        `json:"bank_api_connected"`
	Active            Flag        `json:"active"`
	LastTransaction   string      `json:"last_transaction"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`

	// The fields below are only filled by GetBankAccount.
	IdentificationNumber string           `json:"identification_number"` // CMND/CCCD
	PhoneNumber          string           `json:"phone_number"`
	VirtualAccounts      []VirtualAccount `json:"va_list"`
}

// VirtualAccount (VA) is a sub-account of a bank account, used to tell payers apart.
type VirtualAccount struct {
	XID           string `json:"xid"`
	SubAccount    string `json:"sub_account"`
	SubHolderName string `json:"sub_holder_name"`
	Active        Flag   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// ListBankAccountsParams filters ListBankAccounts. Zero values are ignored.
type ListBankAccountsParams struct {
	PageParams
	CompanyXID string
	BankID     string
	Query      string
}

// ListBankAccounts returns a page of linked bank accounts, newest first. params may be nil.
func (c *Client) ListBankAccounts(ctx context.Context, params *ListBankAccountsParams) (*List[BankAccount], error) {
	q := url.Values{}
	if params != nil {
		pagination(q, params.Page, params.PerPage)
		setIf(q, "company_xid", params.CompanyXID)
		setIf(q, "bank_id", params.BankID)
		setIf(q, "q", params.Query)
	}
	var out List[BankAccount]
	if err := c.list(ctx, "/v1/bank-account", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetBankAccount returns a bank account with its virtual accounts.
func (c *Client) GetBankAccount(ctx context.Context, xid string) (*BankAccount, error) {
	var out BankAccount
	if err := c.get(ctx, http.MethodGet, "/v1/bank-account/"+url.PathEscape(xid), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
