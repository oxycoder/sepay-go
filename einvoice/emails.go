package einvoice

import (
	"context"
	"net/http"
	"net/url"
)

// EmailLog is one invoice email delivery record.
type EmailLog struct {
	InvoiceSeries     string      `json:"invoice_series"`
	InvoiceNumber     string      `json:"invoice_number"`
	FullInvoiceNumber string      `json:"full_invoice_number"`
	ReferenceCode     string      `json:"reference_code"`
	ToEmail           string      `json:"to_email"`
	Status            EmailStatus `json:"status"`
	SentAt            string      `json:"sent_at"`    // YYYY-MM-DD HH:mm:ss; empty until sent
	CreatedAt         string      `json:"created_at"` // YYYY-MM-DD HH:mm:ss
}

// ListEmailLogsParams filters ListEmailLogs. All fields are optional;
// string filters match partially.
type ListEmailLogsParams struct {
	Page          int
	PerPage       int // default 20, max 100
	InvoiceNumber string
	Status        EmailStatus
	ToEmail       string
	ReferenceCode string
}

// ListEmailLogs returns a page of invoice email delivery records. params may be nil.
func (c *Client) ListEmailLogs(ctx context.Context, params *ListEmailLogsParams) (*List[EmailLog], error) {
	q := url.Values{}
	if params != nil {
		pagination(q, params.Page, params.PerPage)
		setIf(q, "invoice_number", params.InvoiceNumber)
		setIf(q, "status", string(params.Status))
		setIf(q, "to_email", params.ToEmail)
		setIf(q, "reference_code", params.ReferenceCode)
	}
	var out List[EmailLog]
	if err := c.do(ctx, http.MethodGet, "/v1/email-logs", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func setIf(q url.Values, key, value string) {
	if value != "" {
		q.Set(key, value)
	}
}
