package bankhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// APIError is returned when the API responds with HTTP status >= 400.
type APIError struct {
	StatusCode int    // HTTP status code
	Message    string // human-readable message (may be Vietnamese); do not branch on it
	// Errors holds per-field validation messages, keyed by request field name.
	Errors map[string]string
}

func (e *APIError) Error() string {
	msg := fmt.Sprintf("bankhub: %d: %s", e.StatusCode, e.Message)
	if len(e.Errors) > 0 {
		fields := make([]string, 0, len(e.Errors))
		for k, v := range e.Errors {
			fields = append(fields, k+": "+v)
		}
		sort.Strings(fields)
		msg += " (" + strings.Join(fields, "; ") + ")"
	}
	return msg
}

// IsNotFound reports whether err is an APIError with HTTP status 404.
func IsNotFound(err error) bool {
	return hasStatus(err, http.StatusNotFound)
}

// IsValidationError reports whether err is an APIError with HTTP status 400.
func IsValidationError(err error) bool {
	return hasStatus(err, http.StatusBadRequest)
}

func hasStatus(err error, status int) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == status
}

// errorBody covers the error shapes the API returns:
//
//	{"code": 400, "message": "Bad request", "errors": {"field": "..."}}
//	{"status": 400, "error": "...", "messages": {"error": ["..."]}}
type errorBody struct {
	Message  string                     `json:"message"`
	Error    json.RawMessage            `json:"error"`
	Errors   map[string]json.RawMessage `json:"errors"`
	Messages struct {
		Error []string `json:"error"`
	} `json:"messages"`
}

func parseAPIError(status int, raw []byte) *APIError {
	apiErr := &APIError{StatusCode: status}
	var b errorBody
	if err := json.Unmarshal(raw, &b); err != nil {
		apiErr.Message = strings.TrimSpace(string(raw))
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(status)
		}
		return apiErr
	}

	apiErr.Message = b.Message
	if apiErr.Message == "" {
		apiErr.Message = rawText(b.Error)
	}
	if apiErr.Message == "" {
		apiErr.Message = strings.Join(b.Messages.Error, "; ")
	}
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(status)
	}
	if len(b.Errors) > 0 {
		apiErr.Errors = make(map[string]string, len(b.Errors))
		for k, v := range b.Errors {
			apiErr.Errors[k] = rawText(v)
		}
	}
	return apiErr
}

// rawText renders a JSON string, or a list of strings, as plain text.
// Anything else is returned as raw JSON.
func rawText(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var list []string
	if json.Unmarshal(raw, &list) == nil {
		return strings.Join(list, "; ")
	}
	return string(raw)
}
