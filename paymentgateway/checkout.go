package paymentgateway

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"html"
	"net/url"
	"strconv"
	"strings"
)

// CheckoutPath is the path of the checkout form endpoint on the checkout host.
const CheckoutPath = "/v1/checkout/init"

// signedFields are the form fields covered by the signature, in the order
// SePay concatenates them. Do not reorder.
var signedFields = []string{
	"order_amount", "merchant", "currency", "operation",
	"order_description", "order_invoice_number", "customer_id",
	"payment_method", "success_url", "error_url", "cancel_url",
}

// CheckoutRequest describes a one-time payment.
type CheckoutRequest struct {
	// OrderInvoiceNumber identifies the order in your system. It must be
	// unique; IPNs, CancelOrder and VoidTransaction refer to it. Required.
	OrderInvoiceNumber string
	// OrderAmount in VND. Must be above 0 for OperationPurchase.
	OrderAmount int64
	// OrderDescription is shown to the customer. Required.
	OrderDescription string
	// Operation defaults to OperationPurchase.
	Operation Operation
	// Currency defaults to CurrencyVND, the only supported currency.
	Currency Currency
	// PaymentMethod preselects a method. When empty, SePay uses the only
	// method enabled for the merchant or lets the customer choose.
	PaymentMethod PaymentMethod
	// CustomerID is your id for the customer. Optional.
	CustomerID string
	// SuccessURL, ErrorURL and CancelURL are where the customer is sent after
	// paying, on failure and on cancel. They must be publicly reachable.
	// Optional.
	SuccessURL string
	ErrorURL   string
	CancelURL  string
}

// Field is one form field.
type Field struct {
	Name  string
	Value string
}

// CheckoutForm is a signed payment form. Render it with HTML, or post Fields
// to Action from your own template, keeping their order.
type CheckoutForm struct {
	Action string  // URL the form is POSTed to
	Fields []Field // form fields, signature last
}

// NewCheckout validates req and returns the signed form that sends the
// customer to SePay's checkout page. It makes no network call.
func (c *Client) NewCheckout(req *CheckoutRequest) (*CheckoutForm, error) {
	if req == nil {
		return nil, errors.New("paymentgateway: nil CheckoutRequest")
	}
	op := req.Operation
	if op == "" {
		op = OperationPurchase
	}
	currency := req.Currency
	if currency == "" {
		currency = CurrencyVND
	}
	switch {
	case c.merchantID == "":
		return nil, errors.New("paymentgateway: merchant id is empty")
	case req.OrderInvoiceNumber == "":
		return nil, errors.New("paymentgateway: OrderInvoiceNumber is required")
	case req.OrderDescription == "":
		return nil, errors.New("paymentgateway: OrderDescription is required")
	case req.OrderAmount < 0, op == OperationPurchase && req.OrderAmount == 0:
		return nil, errors.New("paymentgateway: OrderAmount must be above 0")
	}

	values := map[string]string{
		"order_amount":         strconv.FormatInt(req.OrderAmount, 10),
		"merchant":             c.merchantID,
		"currency":             string(currency),
		"operation":            string(op),
		"order_description":    req.OrderDescription,
		"order_invoice_number": req.OrderInvoiceNumber,
		"customer_id":          req.CustomerID,
		"payment_method":       string(req.PaymentMethod),
		"success_url":          req.SuccessURL,
		"error_url":            req.ErrorURL,
		"cancel_url":           req.CancelURL,
	}
	fields := make([]Field, 0, len(signedFields)+1)
	for _, name := range signedFields {
		if v := values[name]; v != "" {
			fields = append(fields, Field{Name: name, Value: v})
		}
	}
	fields = append(fields, Field{Name: "signature", Value: Sign(fields, c.secretKey)})
	return &CheckoutForm{Action: c.checkoutURL + CheckoutPath, Fields: fields}, nil
}

// Sign returns the checkout signature of fields:
// base64(HMAC-SHA256(secretKey, "name1=value1,name2=value2,...")), over the
// signed fields present in fields, in SePay's fixed order. Other fields (such
// as "signature") are ignored. NewCheckout calls it for you; use it when you
// build the form yourself.
func Sign(fields []Field, secretKey string) string {
	byName := make(map[string]string, len(fields))
	for _, f := range fields {
		byName[f.Name] = f.Value
	}
	parts := make([]string, 0, len(signedFields))
	for _, name := range signedFields {
		if v, ok := byName[name]; ok {
			parts = append(parts, name+"="+v)
		}
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(strings.Join(parts, ",")))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Values returns the form fields as url.Values, e.g. to post them from a test.
func (f *CheckoutForm) Values() url.Values {
	v := make(url.Values, len(f.Fields))
	for _, field := range f.Fields {
		v.Set(field.Name, field.Value)
	}
	return v
}

// HTML renders the form with hidden inputs and a submit button labelled
// buttonLabel. All values are HTML-escaped.
func (f *CheckoutForm) HTML(buttonLabel string) string {
	var b strings.Builder
	b.WriteString(`<form action="` + html.EscapeString(f.Action) + `" method="POST">` + "\n")
	for _, field := range f.Fields {
		b.WriteString(`  <input type="hidden" name="` + html.EscapeString(field.Name) +
			`" value="` + html.EscapeString(field.Value) + `" />` + "\n")
	}
	b.WriteString(`  <button type="submit">` + html.EscapeString(buttonLabel) + "</button>\n</form>")
	return b.String()
}
