package einvoice

import (
	"errors"
	"fmt"
	"net/http"
)

// Error codes returned in APIError.Code.
const (
	CodeBadRequest                   = "BAD_REQUEST"
	CodeUnauthorized                 = "UNAUTHORIZED"
	CodeTokenExpired                 = "TOKEN_EXPIRED"
	CodeValidationError              = "VALIDATION_ERROR"
	CodeMethodNotAllowed             = "METHOD_NOT_ALLOWED"
	CodeInternalError                = "INTERNAL_ERROR"
	CodeDocumentExisted              = "EINVOICE_DOCUMENT_EXISTED"
	CodeQuotaUsedUp                  = "QUOTA_HAS_BEEN_USERD_UP" // sic: spelled exactly as the API returns it
	CodeRegistrationNotTaxApproved   = "REGISTRATION_NOT_TAX_APPROVED"
	CodeInvoiceDateBeforeTaxApproval = "INVOICE_DATE_BEFORE_TAX_APPROVAL"
	CodeRegistrationNotComplete      = "REGISTRATION_NOT_COMPLETE"
	CodeBillingIsUnpaid              = "BILLING_IS_UNPAID"
	CodeInvalidOriginInvoice         = "INVALID_ORIGIN_INVOICE"
	CodeIssuedDateInvalid            = "ISSUED_DATE_INVALID" // issued_date is not today
	CodeInvalidIssuedDate            = "INVALID_ISSUED_DATE" // issued_date is not "Y-m-d H:i:s"
	CodeInvalidAdjustmentType        = "INVALID_ADJUSTMENT_TYPE"
	CodeInvalidIsDraftType           = "INVALID_IS_DRAFT_TYPE"
	CodeAccountNotFound              = "EINVOICE_ACCOUNT_NOT_FOUND"
	CodeDraftNotFound                = "EINVOICE_DRAFT_NOT_FOUND"
	CodeReferenceCodeRequired        = "REFERENCE_CODE_REQUIRED"
	CodeCreateInvoiceDataFail        = "CREATE_INVOICE_DATA_FAIL"
	CodeSellerStoreNotFound          = "SELLER_STORE_NOT_FOUND"
	CodeSellerStoreInactive          = "SELLER_STORE_INACTIVE"
	CodeSellerStoreInUse             = "SELLER_STORE_IN_USE"
	CodeSellerStoreDuplicateCode     = "SELLER_STORE_DUPLICATE_CODE"
	CodeSellerStoreDuplicateName     = "SELLER_STORE_DUPLICATE_NAME"
	CodeSellerStoreDuplicateAddress  = "SELLER_STORE_DUPLICATE_ADDRESS"
)

// APIError is returned when the API responds with an error.
type APIError struct {
	StatusCode int    // HTTP status code
	Code       string // error.code, e.g. CodeValidationError; may be empty
	Message    string // human-readable message (Vietnamese); do not branch on it
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("einvoice: %d %s: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("einvoice: %d: %s", e.StatusCode, e.Message)
}

// IsNotFound reports whether err is an APIError with HTTP status 404.
//
// Before resending a Failed invoice, call GetInvoice with your reference code:
// a not-found error means no invoice was recorded and it is safe to resend.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

// HasCode reports whether err is an APIError with the given error code.
func HasCode(err error, code string) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.Code == code
}
