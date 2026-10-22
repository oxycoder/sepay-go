package webhook

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// maxBodySize caps webhook request bodies.
const maxBodySize = 1 << 20

// HandlerFunc processes one webhook delivery. Returning an error makes the
// handler answer 500 so SePay retries the delivery.
type HandlerFunc func(ctx context.Context, tx *Transaction) error

// Option configures Handler.
type Option func(*config)

type config struct {
	verify    func(r *http.Request, body []byte) error
	tolerance time.Duration
}

// WithAPIKey requires "Authorization: Apikey {apiKey}" on every request.
// Use it when the webhook's authentication is set to API Key.
func WithAPIKey(apiKey string) Option {
	return func(c *config) {
		c.verify = func(r *http.Request, _ []byte) error { return VerifyAPIKey(r.Header, apiKey) }
	}
}

// WithHMAC requires a valid X-SePay-Signature for secret on every request and
// rejects timestamps outside the tolerance (DefaultTolerance unless set with
// WithTolerance). Use it when the webhook's authentication is set to
// HMAC-SHA256, which SePay recommends.
func WithHMAC(secret string) Option {
	return func(c *config) {
		c.verify = func(r *http.Request, body []byte) error {
			return VerifySignature(r.Header, body, secret, c.tolerance)
		}
	}
}

// WithBearer requires "Authorization: Bearer {token}" where valid(token)
// reports true. Use it when the webhook's authentication is set to OAuth 2.0:
// SePay fetches the token from your token endpoint, so valid should check it
// against the tokens that endpoint issued.
func WithBearer(valid func(ctx context.Context, token string) bool) Option {
	return func(c *config) {
		c.verify = func(r *http.Request, _ []byte) error {
			if token, ok := bearerToken(r.Header); ok && valid(r.Context(), token) {
				return nil
			}
			return ErrUnauthorized
		}
	}
}

// WithTolerance sets how far X-SePay-Timestamp may be from the current time
// when using WithHMAC. d <= 0 disables the timestamp check.
func WithTolerance(d time.Duration) Option {
	return func(c *config) { c.tolerance = d }
}

// Handler returns an http.Handler for SePay webhook deliveries. It accepts
// only POST, checks the credential set by WithAPIKey, WithHMAC or WithBearer,
// decodes the body (see Parse) and calls fn. Without one of those options
// every request is accepted, which is only safe for local testing.
//
// Responses follow what SePay expects: 200 with {"success": true} when fn
// returns nil; 401, 400 or 500 with {"success": false} otherwise. SePay retries
// any non-2xx response, and may also deliver a transaction more than once
// through manual replays, so make fn idempotent on Transaction.ID.
func Handler(fn HandlerFunc, opts ...Option) http.Handler {
	c := &config{tolerance: DefaultTolerance}
	for _, opt := range opts {
		opt(c)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			respond(w, http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodySize))
		if err != nil {
			respond(w, http.StatusBadRequest)
			return
		}
		if c.verify != nil {
			if err := c.verify(r, body); err != nil {
				respond(w, http.StatusUnauthorized)
				return
			}
		}
		tx, err := Parse(r.Header.Get("Content-Type"), body)
		if err != nil {
			respond(w, http.StatusBadRequest)
			return
		}
		if err := fn(r.Context(), tx); err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		respond(w, http.StatusOK)
	})
}

// respond writes {"success": status == 200} with the given status.
func respond(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"success":%t}`, status == http.StatusOK)
}
