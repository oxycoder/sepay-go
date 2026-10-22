package einvoice

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// Ptr returns a pointer to v. Use it for optional fields where the zero value
// is meaningful, e.g. TaxRate: einvoice.Ptr(einvoice.TaxRate0).
func Ptr[T any](v T) *T { return &v }

// Currency of an invoice.
type Currency string

const (
	CurrencyVND Currency = "VND"
	CurrencyUSD Currency = "USD"
	CurrencyCAD Currency = "CAD"
)

// PaymentMethod of an invoice.
type PaymentMethod string

const (
	PaymentCash           PaymentMethod = "TM"    // Tiền mặt
	PaymentBankTransfer   PaymentMethod = "CK"    // Chuyển khoản
	PaymentCashOrTransfer PaymentMethod = "TM/CK" // Tiền mặt và chuyển khoản
	PaymentOther          PaymentMethod = "KHAC"  // Khác
)

// Template codes commonly returned in ProviderAccount.Templates.
const (
	TemplateVAT   = "1" // Hóa đơn giá trị gia tăng
	TemplateSales = "2" // Hóa đơn bán hàng (tax_rate and tax_amount are ignored)
)

// BuyerType distinguishes individuals from companies.
type BuyerType string

const (
	BuyerPersonal BuyerType = "personal"
	BuyerCompany  BuyerType = "company"
)

// LineType is the kind of an invoice line.
type LineType int

const (
	LineTypeGoods     LineType = 1 // Hàng hóa/dịch vụ
	LineTypePromotion LineType = 2 // Hàng khuyến mại
	LineTypeDiscount  LineType = 3 // Chiết khấu thương mại
	LineTypeNote      LineType = 4 // Ghi chú (printed on the invoice)
)

// TaxRate is a VAT rate in percent.
type TaxRate int

const (
	TaxRateNotSubject  TaxRate = -2 // Không chịu thuế
	TaxRateNotDeclared TaxRate = -1 // Không kê khai, tính nộp thuế GTGT
	TaxRate0           TaxRate = 0
	TaxRate5           TaxRate = 5
	TaxRate8           TaxRate = 8
	TaxRate10          TaxRate = 10
)

// AdjustmentType is the kind of an adjustment invoice.
type AdjustmentType string

const (
	AdjustIncrease AdjustmentType = "increase" // Điều chỉnh tăng
	AdjustDecrease AdjustmentType = "decrease" // Điều chỉnh giảm
	AdjustCombined AdjustmentType = "combined" // Điều chỉnh tăng/giảm, sign per line
	AdjustInfo     AdjustmentType = "info"     // Điều chỉnh thông tin, total must be 0
)

// ProcessStatus is the state of an asynchronous create/issue request.
type ProcessStatus string

const (
	StatusPending ProcessStatus = "Pending"
	StatusSuccess ProcessStatus = "Success"
	StatusFailed  ProcessStatus = "Failed"
)

// InvoiceStatus is the SePay-side state of an invoice.
type InvoiceStatus string

const (
	InvoiceDraft     InvoiceStatus = "draft"
	InvoiceIssued    InvoiceStatus = "issued"
	InvoiceCancelled InvoiceStatus = "cancelled"
)

// DocumentType classifies an invoice document.
type DocumentType int

const (
	DocDraft          DocumentType = 1  // Hóa đơn nháp
	DocNew            DocumentType = 2  // Hóa đơn mới
	DocCancelled      DocumentType = 3  // Hóa đơn đã hủy/xóa bỏ
	DocReplacement    DocumentType = 4  // Hóa đơn thay thế
	DocReplaced       DocumentType = 5  // Hóa đơn bị thay thế
	DocAdjustIncrease DocumentType = 6  // Hóa đơn điều chỉnh tăng
	DocAdjustDecrease DocumentType = 7  // Hóa đơn điều chỉnh giảm
	DocAdjustInfo     DocumentType = 8  // Hóa đơn điều chỉnh thông tin
	DocAdjustCombined DocumentType = 9  // Hóa đơn điều chỉnh tăng/giảm
	DocAdjusted       DocumentType = 10 // Hóa đơn bị điều chỉnh
)

// Source tells where an invoice was created.
type Source string

const (
	SourceManual Source = "manual" // created in the web UI
	SourceAPI    Source = "api"
)

// FileType of a downloaded invoice file.
type FileType string

const (
	FilePDF FileType = "pdf"
	FileXML FileType = "xml"
)

// EmailStatus is the delivery state of an invoice email.
type EmailStatus string

const (
	EmailPending         EmailStatus = "pending"
	EmailSending         EmailStatus = "sending"
	EmailSent            EmailStatus = "sent"
	EmailFailed          EmailStatus = "failed" // may still be retried
	EmailRetrying        EmailStatus = "retrying"
	EmailPermanentFailed EmailStatus = "permanent_failed"
)

// Number is a numeric response value. The API encodes some amounts as JSON
// numbers and others as decimal strings (e.g. "2.00"); Number accepts both.
type Number float64

// UnmarshalJSON implements json.Unmarshaler.
func (n *Number) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return fmt.Errorf("einvoice: invalid number %q", b)
	}
	*n = Number(f)
	return nil
}

// Paging describes a page of a list response.
type Paging struct {
	PerPage     int  `json:"per_page"`
	Total       int  `json:"total"`
	HasMore     bool `json:"has_more"`
	CurrentPage int  `json:"current_page"`
	PageCount   int  `json:"page_count"`
}

// List is one page of a paginated list.
type List[T any] struct {
	Paging Paging `json:"paging"`
	Items  []T    `json:"items"`
}

// Buyer is the invoice buyer.
type Buyer struct {
	Type       BuyerType `json:"type,omitempty"`
	Name       string    `json:"name,omitempty"`
	LegalName  string    `json:"legal_name,omitempty"` // used when Type is BuyerCompany
	TaxCode    string    `json:"tax_code,omitempty"`
	Address    string    `json:"address,omitempty"`
	Email      string    `json:"email,omitempty"` // invoice is emailed here
	Phone      string    `json:"phone,omitempty"`
	BuyerCode  string    `json:"buyer_code,omitempty"`
	NationalID string    `json:"national_id,omitempty"` // CCCD / số định danh cá nhân
}

// SellerStore is the seller business location attached to an invoice.
type SellerStore struct {
	XID  string `json:"xid"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// vnZone is Vietnam time (UTC+7); the API interprets dates in this zone.
var vnZone = time.FixedZone("ICT", 7*60*60)

const dateTimeLayout = "2006-01-02 15:04:05"

// formatDateTime formats t as "YYYY-MM-DD HH:mm:ss" in Vietnam time; zero t gives "".
func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(vnZone).Format(dateTimeLayout)
}
