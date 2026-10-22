package einvoice

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

// Base URLs of the SePay eInvoice API.
const (
	ProductionURL = "https://einvoice-api.sepay.vn"
	SandboxURL    = "https://einvoice-api-sandbox.sepay.vn"
)

const (
	defaultUserAgent    = "sepay-go"
	defaultPollInterval = 3 * time.Second
)

// Client talks to the SePay eInvoice API. It is safe for concurrent use.
type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	userAgent    string
	pollInterval time.Duration

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

// WithPollInterval sets the base interval used by WaitForCreate and WaitForIssue (default 3s).
func WithPollInterval(d time.Duration) Option {
	return func(c *Client) { c.pollInterval = d }
}

// New creates a Client authenticating with the given client_id and client_secret.
func New(clientID, clientSecret string, opts ...Option) *Client {
	c := &Client{
		baseURL:      ProductionURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		userAgent:    defaultUserAgent,
		pollInterval: defaultPollInterval,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// envelope is the common response wrapper: {"success": ..., "data": ..., "error": {...}}.
type envelope struct {
	Success *bool           `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Message string `json:"message"`
}

// do sends an authenticated request and decodes the "data" field into out.
// If the token is rejected with 401, it is refreshed and the request retried once.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	var payload []byte
	if body != nil {
		var err error
		if payload, err = json.Marshal(body); err != nil {
			return fmt.Errorf("einvoice: encode request: %w", err)
		}
	}

	for attempt := 0; ; attempt++ {
		token, err := c.accessToken(ctx)
		if err != nil {
			return err
		}
		req, err := c.newRequest(ctx, method, path, query, payload)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)

		err = c.send(req, out)
		var apiErr *APIError
		if attempt == 0 && errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
			c.invalidateToken(token)
			continue
		}
		return err
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, query url.Values, payload []byte) (*http.Request, error) {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	var body io.Reader
	if payload != nil {
		body = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, fmt.Errorf("einvoice: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) send(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("einvoice: %s %s: %w", req.Method, req.URL.Path, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("einvoice: read response: %w", err)
	}

	var env envelope
	jsonErr := json.Unmarshal(raw, &env)

	if resp.StatusCode >= 400 || (jsonErr == nil && (env.Error != nil || (env.Success != nil && !*env.Success))) {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		switch {
		case jsonErr != nil:
			apiErr.Message = strings.TrimSpace(string(raw))
		case env.Error != nil:
			apiErr.Code, apiErr.Message = env.Error.Code, env.Error.Message
		default:
			apiErr.Message = env.Message
		}
		return apiErr
	}
	if jsonErr != nil {
		return fmt.Errorf("einvoice: decode response: %w", jsonErr)
	}
	if out == nil || len(env.Data) == 0 || string(env.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(env.Data, out); err != nil {
		return fmt.Errorf("einvoice: decode response data: %w", err)
	}
	return nil
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
