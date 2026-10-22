package paymentgateway

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Currency of an order. Only VND is supported.
type Currency string

const CurrencyVND Currency = "VND"

// Operation is the kind of checkout.
type Operation string

const (
	OperationPurchase Operation = "PURCHASE"
	OperationVerify   Operation = "VERIFY"
)

// PaymentMethod is how the customer pays.
type PaymentMethod string

const (
	PaymentCard              PaymentMethod = "CARD"                // international card
	PaymentBankTransfer      PaymentMethod = "BANK_TRANSFER"       // VietQR bank transfer
	PaymentNapasBankTransfer PaymentMethod = "NAPAS_BANK_TRANSFER" // NAPAS QR
)

// OrderStatus is the state of an order.
type OrderStatus string

const (
	OrderCaptured                OrderStatus = "CAPTURED"                  // paid
	OrderCancelled               OrderStatus = "CANCELLED"                 // cancelled
	OrderAuthenticationNotNeeded OrderStatus = "AUTHENTICATION_NOT_NEEDED" // waiting for payment
)

// TransactionType is the kind of a transaction on an order.
type TransactionType string

const (
	TransactionPayment TransactionType = "PAYMENT"
	TransactionRefund  TransactionType = "REFUND"
	TransactionVoid    TransactionType = "VOID"
)

// TransactionStatus is the result of a transaction.
type TransactionStatus string

const (
	TransactionApproved TransactionStatus = "APPROVED"
	TransactionDeclined TransactionStatus = "DECLINED"
	TransactionPending  TransactionStatus = "PENDING"
)

// CardFundingMethod tells debit and credit cards apart.
type CardFundingMethod string

const (
	CardDebit  CardFundingMethod = "DEBIT"
	CardCredit CardFundingMethod = "CREDIT"
)

// CardBrand is the card network.
type CardBrand string

const (
	CardVisa       CardBrand = "VISA"
	CardMastercard CardBrand = "MASTERCARD"
	CardJCB        CardBrand = "JCB"
)

// Number is a numeric response value. The API encodes amounts as decimal
// strings (e.g. "50000.00"); Number also accepts JSON numbers.
type Number float64

// UnmarshalJSON implements json.Unmarshaler.
func (n *Number) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return fmt.Errorf("paymentgateway: invalid number %q", b)
	}
	*n = Number(f)
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

// List is one page of a paginated list.
type List[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

// PageParams selects a page of a list. Zero values use the API defaults
// (page 1, 20 per page, at most 100).
type PageParams struct {
	Page    int
	PerPage int
}

// vnZone is Vietnam time (UTC+7); the API reads and writes times in this zone.
var vnZone = time.FixedZone("ICT", 7*60*60)

// timeLayout is the API's "YYYY-MM-DD HH:mm:ss" format.
const timeLayout = time.DateTime

// ParseTime parses an API timestamp ("YYYY-MM-DD HH:mm:ss", Vietnam time),
// such as Order.CreatedAt or Transaction.TransactionDate.
func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation(timeLayout, s, vnZone)
}

// formatDate formats t as "YYYY-MM-DD" in Vietnam time; zero t gives "".
func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(vnZone).Format(time.DateOnly)
}
