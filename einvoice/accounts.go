package einvoice

import (
	"context"
	"net/http"
	"net/url"
)

// ProviderAccount is an e-invoice provider account (e.g. matbao, bkav).
type ProviderAccount struct {
	ID       string `json:"id"` // use as ProviderAccountID when creating invoices
	Provider string `json:"provider"`
	Active   bool   `json:"active"`
	// TaxAuthorityApprovedDate (YYYY-MM-DD) is when the tax authority approved
	// the registration; invoices dated before it are rejected. Empty in sandbox.
	TaxAuthorityApprovedDate string `json:"tax_authority_approved_date"`
	// Templates is only filled by GetProviderAccount.
	Templates []Template `json:"templates"`
}

// Template is an invoice template/series pair usable with a provider account.
type Template struct {
	TemplateCode  string `json:"template_code"`
	InvoiceSeries string `json:"invoice_series"`
	InvoiceLabel  string `json:"invoice_label"`
}

// PageParams selects a page of a list. Zero values use the API defaults.
type PageParams struct {
	Page    int
	PerPage int // max 100
}

// ListProviderAccounts returns a page of provider accounts. params may be nil.
func (c *Client) ListProviderAccounts(ctx context.Context, params *PageParams) (*List[ProviderAccount], error) {
	q := url.Values{}
	if params != nil {
		pagination(q, params.Page, params.PerPage)
	}
	var out List[ProviderAccount]
	if err := c.do(ctx, http.MethodGet, "/v1/provider-accounts", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetProviderAccount returns a provider account with its templates and series.
func (c *Client) GetProviderAccount(ctx context.Context, id string) (*ProviderAccount, error) {
	var out ProviderAccount
	if err := c.do(ctx, http.MethodGet, "/v1/provider-accounts/"+url.PathEscape(id), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
