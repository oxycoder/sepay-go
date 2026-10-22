package bankhub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Base URLs of the SePay Bank Hub API.
const (
	ProductionURL = "https://bankhub-api.sepay.vn"
	SandboxURL    = "https://bankhub-api-sandbox.sepay.vn"
)

const defaultUserAgent = "sepay-go"

// Client talks to the SePay Bank Hub API. It is safe for concurrent use.
type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	userAgent    string

	mu          sync.Mutex
	token       string
	tokenExpiry time.Time
}

// Option configures a Client.
type Option func(*Client)

// WithSandbox points the client at the sandbox environment.
func WithSandbox() Option {
	return WithBaseURL(SandboxURL)
}

// WithBaseURL overrides the API base URL (default ProductionURL).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(baseURL, "/") }
}

// WithHTTPClient sets the HTTP client used for requests (default: a client with a 30s timeout).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// New creates a Client authenticating with the given client_id and client_secret.
func New(clientID, clientSecret string, opts ...Option) *Client {
	c := &Client{
		baseURL:      ProductionURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		userAgent:    defaultUserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// body is an encoded request body.
type body struct {
	contentType string
	data        []byte
}

func jsonBody(v any) (*body, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("bankhub: encode request: %w", err)
	}
	return &body{contentType: "application/json", data: data}, nil
}

// call sends an authenticated JSON request and returns the raw response body.
// If the token is rejected with 401, it is refreshed and the request retried once.
func (c *Client) call(ctx context.Context, method, path string, query url.Values, in any) ([]byte, error) {
	var b *body
	if in != nil {
		var err error
		if b, err = jsonBody(in); err != nil {
			return nil, err
		}
	}
	return c.callBody(ctx, method, path, query, b)
}

func (c *Client) callBody(ctx context.Context, method, path string, query url.Values, b *body) ([]byte, error) {
	for attempt := 0; ; attempt++ {
		token, err := c.accessToken(ctx)
		if err != nil {
			return nil, err
		}
		req, err := c.newRequest(ctx, method, path, query, b)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)

		raw, err := c.send(req)
		var apiErr *APIError
		if attempt == 0 && errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
			c.invalidateToken(token)
			continue
		}
		return raw, err
	}
}

// get sends a request and decodes the response's "data" field into out.
func (c *Client) get(ctx context.Context, method, path string, query url.Values, in, out any) error {
	raw, err := c.call(ctx, method, path, query, in)
	if err != nil {
		return err
	}
	_, err = decodeData(raw, out)
	return err
}

// list sends a GET request and decodes the whole response ("data" and "meta") into out.
func (c *Client) list(ctx context.Context, path string, query url.Values, out any) error {
	raw, err := c.call(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("bankhub: decode response: %w", err)
	}
	return nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, b *body) (*http.Request, error) {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	var r io.Reader
	if b != nil {
		r = bytes.NewReader(b.data)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, r)
	if err != nil {
		return nil, fmt.Errorf("bankhub: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if b != nil {
		req.Header.Set("Content-Type", b.contentType)
	}
	return req, nil
}

// send performs req and returns the response body, or an *APIError for HTTP status >= 400.
func (c *Client) send(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bankhub: %s %s: %w", req.Method, req.URL.Path, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("bankhub: read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, raw)
	}
	return raw, nil
}

// decodeData decodes the "data" field of raw into out. Responses without a
// "data" field (e.g. link tokens) are decoded as a whole. It reports false if
// "data" is present but null.
func decodeData(raw []byte, out any) (bool, error) {
	if out == nil {
		return true, nil
	}
	var env map[string]json.RawMessage
	if err := json.Unmarshal(raw, &env); err != nil {
		return false, fmt.Errorf("bankhub: decode response: %w", err)
	}
	data, ok := env["data"]
	if !ok {
		data = raw
	}
	if string(data) == "null" {
		return false, nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		return false, fmt.Errorf("bankhub: decode response data: %w", err)
	}
	return true, nil
}

// pagination adds page/per_page to q when set.
func pagination(q url.Values, page, perPage int) {
	if page > 0 {
		q.Set("page", fmt.Sprint(page))
	}
	if perPage > 0 {
		q.Set("per_page", fmt.Sprint(perPage))
	}
}

// setIf sets q[key] = v when v is not empty.
func setIf(q url.Values, key, v string) {
	if v != "" {
		q.Set(key, v)
	}
}
