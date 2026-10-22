package paymentgateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Base URLs of the SePay Payment Gateway. The API hosts serve the REST API;
// the checkout hosts receive the payment form.
const (
	ProductionURL         = "https://pgapi.sepay.vn"
	SandboxURL            = "https://pgapi-sandbox.sepay.vn"
	ProductionCheckoutURL = "https://pay.sepay.vn"
	SandboxCheckoutURL    = "https://pay-sandbox.sepay.vn"
)

const defaultUserAgent = "sepay-go"

// Client talks to the SePay Payment Gateway. It is safe for concurrent use.
type Client struct {
	baseURL     string
	checkoutURL string
	merchantID  string
	secretKey   string
	httpClient  *http.Client
	userAgent   string
}

// Option configures a Client.
type Option func(*Client)

// WithSandbox points both the API and the checkout form at the sandbox
// environment.
func WithSandbox() Option {
	return func(c *Client) {
		c.baseURL = SandboxURL
		c.checkoutURL = SandboxCheckoutURL
	}
}

// WithBaseURL overrides the REST API base URL (default ProductionURL).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithCheckoutURL overrides the checkout base URL the payment form is posted
// to (default ProductionCheckoutURL).
func WithCheckoutURL(checkoutURL string) Option {
	return func(c *Client) { c.checkoutURL = strings.TrimRight(checkoutURL, "/") }
}

// WithHTTPClient sets the HTTP client used for requests (default: a client with a 30s timeout).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// New creates a Client for the given merchant_id and secret_key.
func New(merchantID, secretKey string, opts ...Option) *Client {
	c := &Client{
		baseURL:     ProductionURL,
		checkoutURL: ProductionCheckoutURL,
		merchantID:  merchantID,
		secretKey:   secretKey,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		userAgent:   defaultUserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// MerchantID returns the merchant_id the client was created with.
func (c *Client) MerchantID() string { return c.merchantID }

// call sends an authenticated JSON request and decodes the response into out
// (if non-nil).
func (c *Client) call(ctx context.Context, method, path string, query url.Values, in, out any) error {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	var r io.Reader
	if in != nil {
		data, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("paymentgateway: encode request: %w", err)
		}
		r = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, r)
	if err != nil {
		return fmt.Errorf("paymentgateway: build request: %w", err)
	}
	req.SetBasicAuth(c.merchantID, c.secretKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("paymentgateway: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("paymentgateway: read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return parseAPIError(resp.StatusCode, raw)
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("paymentgateway: decode response: %w", err)
	}
	return nil
}
