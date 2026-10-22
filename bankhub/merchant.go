package bankhub

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// MerchantConfig is the merchant-wide configuration.
type MerchantConfig struct {
	// TransactionWhitelistRemarkKeywords limits IPN notifications to
	// transactions whose content contains one of these keywords
	// (case-insensitive). Empty means every transaction is notified.
	TransactionWhitelistRemarkKeywords []string `json:"transaction_whitelist_remark_keywords"`
	LogoURL                            string   `json:"logo_url"`
}

// GetMerchantConfig returns the merchant configuration.
func (c *Client) GetMerchantConfig(ctx context.Context) (*MerchantConfig, error) {
	var out MerchantConfig
	if err := c.get(ctx, http.MethodGet, "/v1/merchant-config", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMerchantConfigRequest is the body of UpdateMerchantConfig.
type UpdateMerchantConfigRequest struct {
	// TransactionWhitelistRemarkKeywords replaces the keyword whitelist.
	// Nil or empty clears it, so every transaction is notified.
	TransactionWhitelistRemarkKeywords []string `json:"transaction_whitelist_remark_keywords"`
}

// UpdateMerchantConfig updates the merchant configuration and returns the result.
func (c *Client) UpdateMerchantConfig(ctx context.Context, req *UpdateMerchantConfigRequest) (*MerchantConfig, error) {
	in := *req
	if in.TransactionWhitelistRemarkKeywords == nil {
		in.TransactionWhitelistRemarkKeywords = []string{}
	}
	var out MerchantConfig
	if err := c.get(ctx, http.MethodPost, "/v1/merchant-config", nil, &in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UploadMerchantLogo uploads the logo shown on every hosted link and returns
// its URL. filename must end in .png, .jpg or .jpeg; the file must be at most
// 1MB. The previous logo is deleted.
func (c *Client) UploadMerchantLogo(ctx context.Context, filename string, logo io.Reader) (string, error) {
	var contentType string
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	default:
		return "", fmt.Errorf("bankhub: logo must be a .png, .jpg or .jpeg file, got %q", filename)
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="logo"; filename=%q`, filepath.Base(filename)))
	h.Set("Content-Type", contentType)
	part, err := mw.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("bankhub: encode logo: %w", err)
	}
	if _, err := io.Copy(part, logo); err != nil {
		return "", fmt.Errorf("bankhub: read logo: %w", err)
	}
	if err := mw.Close(); err != nil {
		return "", fmt.Errorf("bankhub: encode logo: %w", err)
	}

	raw, err := c.callBody(ctx, http.MethodPost, "/v1/merchant-config/logo", nil,
		&body{contentType: mw.FormDataContentType(), data: buf.Bytes()})
	if err != nil {
		return "", err
	}
	var out struct {
		LogoURL string `json:"logo_url"`
	}
	if _, err := decodeData(raw, &out); err != nil {
		return "", err
	}
	return out.LogoURL, nil
}

// MerchantCounterParams selects the days of GetMerchantCounter. Zero dates
// default to today (Vietnam time).
type MerchantCounterParams struct {
	StartDate time.Time
	EndDate   time.Time
}

// GetMerchantCounter counts the merchant's transactions per day and company.
// params may be nil. DayCount.Transaction only counts incoming transactions.
func (c *Client) GetMerchantCounter(ctx context.Context, params *MerchantCounterParams) (*Counter, error) {
	q := url.Values{}
	if params != nil {
		setIf(q, "start_date", formatDate(params.StartDate))
		setIf(q, "end_date", formatDate(params.EndDate))
	}
	var out Counter
	if err := c.get(ctx, http.MethodGet, "/v1/merchant-config/counter", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
