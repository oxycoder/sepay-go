package bankhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Ptr returns a pointer to v. Use it for optional fields where the zero value
// is meaningful, e.g. Active: bankhub.Ptr(false).
func Ptr[T any](v T) *T { return &v }

// CompanyStatus is the state of a company.
type CompanyStatus string

const (
	CompanyPending    CompanyStatus = "Pending"
	CompanyActive     CompanyStatus = "Active"
	CompanySuspended  CompanyStatus = "Suspended"
	CompanyTerminated CompanyStatus = "Terminated"
	CompanyCancelled  CompanyStatus = "Cancelled"
	CompanyFraud      CompanyStatus = "Fraud"
)

// TransferType is the direction of a transaction.
type TransferType string

const (
	TransferCredit TransferType = "credit" // Tiền vào
	TransferDebit  TransferType = "debit"  // Tiền ra
)

// AccountType of a bank account. The API uses "individual" and "personal"
// interchangeably for personal accounts.
type AccountType string

const (
	AccountIndividual AccountType = "individual"
	AccountPersonal   AccountType = "personal"
	AccountEnterprise AccountType = "enterprise"
)

// Flag is a boolean the API encodes as 1/0, "1"/"0" or true/false.
type Flag bool

// UnmarshalJSON implements json.Unmarshaler.
func (f *Flag) UnmarshalJSON(b []byte) error {
	switch string(bytes.Trim(b, `"`)) {
	case "1", "true":
		*f = true
	case "0", "false", "", "null":
		*f = false
	default:
		return fmt.Errorf("bankhub: invalid flag %s", b)
	}
	return nil
}

// MarshalJSON implements json.Marshaler, encoding as 1 or 0.
func (f Flag) MarshalJSON() ([]byte, error) {
	if f {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

// Number is a numeric response value. The API encodes some amounts as JSON
// numbers and others as decimal strings; Number accepts both.
type Number float64

// UnmarshalJSON implements json.Unmarshaler.
func (n *Number) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return fmt.Errorf("bankhub: invalid number %q", b)
	}
	*n = Number(f)
	return nil
}

// Text is a string response value the API may also encode as a number.
type Text string

// UnmarshalJSON implements json.Unmarshaler.
func (t *Text) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*t = Text(s)
		return nil
	}
	*t = Text(b)
	return nil
}

// Meta describes a page of a list response.
type Meta struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	Total       int  `json:"total"`
	PageCount   int  `json:"page_count"`
	HasMore     bool `json:"has_more"`
}

// UnmarshalJSON accepts both paging shapes the API uses:
// {current_page, page_count, has_more} and {page, total_pages}.
func (m *Meta) UnmarshalJSON(b []byte) error {
	var w struct {
		CurrentPage int   `json:"current_page"`
		Page        int   `json:"page"`
		PerPage     int   `json:"per_page"`
		Total       int   `json:"total"`
		PageCount   int   `json:"page_count"`
		TotalPages  int   `json:"total_pages"`
		HasMore     *bool `json:"has_more"`
	}
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}
	*m = Meta{
		CurrentPage: firstNonZero(w.CurrentPage, w.Page),
		PerPage:     w.PerPage,
		Total:       w.Total,
		PageCount:   firstNonZero(w.PageCount, w.TotalPages),
	}
	if w.HasMore != nil {
		m.HasMore = *w.HasMore
	} else {
		m.HasMore = m.CurrentPage < m.PageCount
	}
	return nil
}

// firstNonZero returns the first non-zero value.
func firstNonZero(a, b int) int {
	if a != 0 {
		return a
	}
	return b
}

// List is one page of a paginated list.
type List[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// PageParams selects a page of a list. Zero values use the API defaults.
type PageParams struct {
	Page    int
	PerPage int
}

// vnZone is Vietnam time (UTC+7); the API interprets dates in this zone.
var vnZone = time.FixedZone("ICT", 7*60*60)

// formatDate formats t as "YYYY-MM-DD" in Vietnam time; zero t gives "".
func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(vnZone).Format(time.DateOnly)
}
