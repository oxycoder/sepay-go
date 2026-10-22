package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Headers SePay sets for HMAC-SHA256 authentication.
const (
	HeaderSignature = "X-SePay-Signature"
	HeaderTimestamp = "X-SePay-Timestamp"
)

// DefaultTolerance is how far X-SePay-Timestamp may be from the current time
// before a request is rejected as a replay.
const DefaultTolerance = 5 * time.Minute

// Verification errors. Handler answers 401 for all of them.
var (
	ErrUnauthorized     = errors.New("webhook: invalid or missing credential")
	ErrMissingSignature = errors.New("webhook: missing signature or timestamp header")
	ErrInvalidSignature = errors.New("webhook: signature does not match")
	ErrTimestampExpired = errors.New("webhook: timestamp outside tolerance")
)

// VerifyAPIKey checks that the Authorization header is "Apikey {apiKey}".
func VerifyAPIKey(h http.Header, apiKey string) error {
	if !secureEqual(h.Get("Authorization"), "Apikey "+apiKey) {
		return ErrUnauthorized
	}
	return nil
}

// VerifySignature checks the HMAC-SHA256 signature SePay sends in
// X-SePay-Signature against the raw request body. body must be the bytes as
// received, not a re-encoded copy of the parsed payload.
//
// It also rejects requests whose X-SePay-Timestamp is more than tolerance away
// from now; tolerance <= 0 disables that check.
func VerifySignature(h http.Header, body []byte, secret string, tolerance time.Duration) error {
	sig, ts := h.Get(HeaderSignature), h.Get(HeaderTimestamp)
	if sig == "" || ts == "" {
		return ErrMissingSignature
	}
	unix, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return ErrInvalidSignature
	}
	got, err := hex.DecodeString(strings.TrimPrefix(sig, "sha256="))
	if err != nil || !hmac.Equal(got, mac(secret, ts, body)) {
		return ErrInvalidSignature
	}
	if tolerance > 0 {
		if d := time.Since(time.Unix(unix, 0)); d > tolerance || d < -tolerance {
			return ErrTimestampExpired
		}
	}
	return nil
}

// Sign returns the X-SePay-Signature value SePay would send for body at the
// given Unix timestamp. Use it to build signed requests in tests.
func Sign(secret string, timestamp int64, body []byte) string {
	return "sha256=" + hex.EncodeToString(mac(secret, strconv.FormatInt(timestamp, 10), body))
}

// mac computes HMAC-SHA256(secret, "{timestamp}.{body}").
func mac(secret, timestamp string, body []byte) []byte {
	m := hmac.New(sha256.New, []byte(secret))
	m.Write([]byte(timestamp))
	m.Write([]byte("."))
	m.Write(body)
	return m.Sum(nil)
}

// bearerToken returns the token from an "Authorization: Bearer {token}" header.
func bearerToken(h http.Header) (string, bool) {
	const prefix = "Bearer "
	auth := h.Get("Authorization")
	if len(auth) <= len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}
	return auth[len(prefix):], true
}

func secureEqual(got, want string) bool {
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}
