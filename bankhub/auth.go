package bankhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// tokenRefreshMargin renews the token this long before it actually expires.
const tokenRefreshMargin = time.Minute

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TTL         int64  `json:"ttl"` // seconds
}

// accessToken returns a cached Bearer token, fetching a new one via POST /v1/token when needed.
func (c *Client) accessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/token", nil, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(c.clientID, c.clientSecret)

	raw, err := c.send(req)
	if err != nil {
		return "", err
	}
	var tok tokenResponse
	if err := json.Unmarshal(raw, &tok); err != nil {
		return "", fmt.Errorf("bankhub: decode token: %w", err)
	}
	if tok.AccessToken == "" {
		return "", fmt.Errorf("bankhub: token response has no access_token")
	}

	ttl := time.Duration(tok.TTL) * time.Second
	if ttl > 2*tokenRefreshMargin {
		ttl -= tokenRefreshMargin
	}
	c.token = tok.AccessToken
	c.tokenExpiry = time.Now().Add(ttl)
	return c.token, nil
}

// invalidateToken drops the cached token if it is still the one that was rejected.
func (c *Client) invalidateToken(rejected string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token == rejected {
		c.token = ""
	}
}
