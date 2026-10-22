package bankhub

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// TransactionAmountUnlimited lifts a company's transaction limit.
const TransactionAmountUnlimited = "Unlimited"

// PaymentCode switches payment code recognition on or off.
type PaymentCode string

const (
	PaymentCodeOn  PaymentCode = "on"
	PaymentCodeOff PaymentCode = "off"
)

// SuffixCharacterType is the character set of a payment code suffix.
type SuffixCharacterType string

const (
	SuffixNumberOnly      SuffixCharacterType = "NumberOnly"
	SuffixNumberAndLetter SuffixCharacterType = "NumberAndLetter"
)

// Company is a company (organization) owning linked bank accounts.
//
// A company only receives IPN notifications when Status is CompanyActive,
// Active is true and Configuration.TransactionAmount is above 0 or
// TransactionAmountUnlimited.
type Company struct {
	XID      string        `json:"xid"`
	FullName string        `json:"full_name"`
	Status   CompanyStatus `json:"status"`
	Active   Flag          `json:"active"`
	// Configuration is only filled by GetCompany.
	Configuration *CompanyConfiguration `json:"configuration"`
	CreatedAt     string                `json:"created_at"`
	UpdatedAt     string                `json:"updated_at"`
}

// CompanyConfiguration holds the payment code and limit settings of a company.
type CompanyConfiguration struct {
	PaymentCode                    PaymentCode         `json:"payment_code"`
	PaymentCodePrefix              string              `json:"payment_code_prefix"`
	PaymentCodeSuffixFrom          int                 `json:"payment_code_suffix_from"`
	PaymentCodeSuffixTo            int                 `json:"payment_code_suffix_to"`
	PaymentCodeSuffixCharacterType SuffixCharacterType `json:"payment_code_suffix_character_type"`
	DataStorageTime                Text                `json:"data_storage_time"`
	// TransactionAmount is a number or TransactionAmountUnlimited.
	TransactionAmount Text `json:"transaction_amount"`
}

// ListCompaniesParams filters ListCompanies. Zero values are ignored.
type ListCompaniesParams struct {
	PageParams
	Query  string // searches full_name
	Status CompanyStatus
	// SortCreatedAt is "asc" or "desc".
	SortCreatedAt string
}

// ListCompanies returns a page of the merchant's companies. params may be nil.
func (c *Client) ListCompanies(ctx context.Context, params *ListCompaniesParams) (*List[Company], error) {
	q := url.Values{}
	if params != nil {
		pagination(q, params.Page, params.PerPage)
		setIf(q, "q", params.Query)
		setIf(q, "status", string(params.Status))
		setIf(q, "sort[created_at]", params.SortCreatedAt)
	}
	var out List[Company]
	if err := c.list(ctx, "/v1/company", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCompany returns a company with its configuration.
func (c *Client) GetCompany(ctx context.Context, xid string) (*Company, error) {
	var out Company
	if err := c.get(ctx, http.MethodGet, "/v1/company/"+url.PathEscape(xid), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCompanyRequest is the body of CreateCompany.
type CreateCompanyRequest struct {
	FullName string `json:"full_name"` // unique within the merchant
	// Status is CompanyActive (usable right away) or CompanyPending.
	Status CompanyStatus `json:"status"`
}

// CreateCompany creates a company. Only XID and FullName are set on the result.
func (c *Client) CreateCompany(ctx context.Context, req *CreateCompanyRequest) (*Company, error) {
	var out Company
	if err := c.get(ctx, http.MethodPost, "/v1/company/create", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCompanyRequest is the body of UpdateCompany. Only set fields are sent.
type UpdateCompanyRequest struct {
	FullName  string        `json:"full_name,omitempty"` // at least 3 characters
	ShortName string        `json:"short_name,omitempty"`
	Status    CompanyStatus `json:"status,omitempty"`

	PaymentCode PaymentCode `json:"payment_code,omitempty"`
	// PaymentCodePrefix is 2-5 letters.
	PaymentCodePrefix              string              `json:"payment_code_prefix,omitempty"`
	PaymentCodeSuffixFrom          *int                `json:"payment_code_suffix_from,omitempty"`
	PaymentCodeSuffixTo            *int                `json:"payment_code_suffix_to,omitempty"`
	PaymentCodeSuffixCharacterType SuffixCharacterType `json:"payment_code_suffix_character_type,omitempty"`

	// TransactionAmount is a number as a string (e.g. "10000000") or TransactionAmountUnlimited.
	TransactionAmount string `json:"transaction_amount,omitempty"`
}

// UpdateCompany updates the given fields of a company.
func (c *Client) UpdateCompany(ctx context.Context, xid string, req *UpdateCompanyRequest) error {
	_, err := c.call(ctx, http.MethodPost, "/v1/company/edit/"+url.PathEscape(xid), nil, req)
	return err
}

// Counter is a per-day transaction count.
type Counter struct {
	Dates []DayCount `json:"dates"`
	Total Count      `json:"total"`
}

// Count is a number of transactions.
type Count struct {
	Transaction    int `json:"transaction"`
	TransactionIn  int `json:"transaction_in"`
	TransactionOut int `json:"transaction_out"`
}

// DayCount is the transaction count of one day.
type DayCount struct {
	Date      string `json:"date"`       // YYYY-MM-DD
	CompanyID Text   `json:"company_id"` // only set by GetMerchantCounter
	Count
}

// GetCompanyCounter counts a company's transactions per day. A zero date
// returns all days, otherwise only that day.
func (c *Client) GetCompanyCounter(ctx context.Context, xid string, date time.Time) (*Counter, error) {
	q := url.Values{}
	setIf(q, "date", formatDate(date))
	var out Counter
	if err := c.get(ctx, http.MethodGet, "/v1/company/counter/"+url.PathEscape(xid), q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
