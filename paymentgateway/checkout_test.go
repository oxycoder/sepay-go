package paymentgateway

import (
	"strings"
	"testing"
)

// docRequest is the example form from the SePay docs.
var docRequest = &CheckoutRequest{
	OrderAmount:        100000,
	OrderDescription:   "Thanh toán đơn hàng #12345",
	OrderInvoiceNumber: "INV_20231201_001",
	SuccessURL:         "https://yoursite.com/payment/success",
	ErrorURL:           "https://yoursite.com/payment/error",
	CancelURL:          "https://yoursite.com/payment/cancel",
}

// docSignature was computed with:
//
//	printf '%s' 'order_amount=100000,merchant=MERCHANT_123,...' | openssl dgst -sha256 -hmac SECRET_KEY -binary | base64
const docSignature = "LddMASma1ObPVFJzgklR1Jye6/mv66LwWNBfydEsTQI="

func TestNewCheckoutMatchesDocs(t *testing.T) {
	form, err := New("MERCHANT_123", "SECRET_KEY").NewCheckout(docRequest)
	if err != nil {
		t.Fatal(err)
	}
	if form.Action != "https://pay.sepay.vn/v1/checkout/init" {
		t.Errorf("Action = %q", form.Action)
	}
	var names []string
	for _, f := range form.Fields {
		names = append(names, f.Name)
	}
	want := "order_amount,merchant,currency,operation,order_description,order_invoice_number,success_url,error_url,cancel_url,signature"
	if got := strings.Join(names, ","); got != want {
		t.Errorf("fields = %s\nwant     %s", got, want)
	}
	v := form.Values()
	if v.Get("signature") != docSignature {
		t.Errorf("signature = %q, want %q", v.Get("signature"), docSignature)
	}
	if v.Get("currency") != "VND" || v.Get("operation") != "PURCHASE" || v.Get("merchant") != "MERCHANT_123" {
		t.Errorf("defaults not applied: %v", v)
	}
}

func TestSignUsesFixedOrderAndIgnoresOtherFields(t *testing.T) {
	shuffled := []Field{
		{"signature", "ignored"},
		{"cancel_url", "https://yoursite.com/payment/cancel"},
		{"merchant", "MERCHANT_123"},
		{"order_description", "Thanh toán đơn hàng #12345"},
		{"error_url", "https://yoursite.com/payment/error"},
		{"currency", "VND"},
		{"order_invoice_number", "INV_20231201_001"},
		{"success_url", "https://yoursite.com/payment/success"},
		{"operation", "PURCHASE"},
		{"order_amount", "100000"},
		{"unknown", "x"},
	}
	if got := Sign(shuffled, "SECRET_KEY"); got != docSignature {
		t.Errorf("Sign = %q, want %q", got, docSignature)
	}
}

func TestNewCheckoutOptionalFieldsAndSandbox(t *testing.T) {
	form, err := New("M", "S", WithSandbox()).NewCheckout(&CheckoutRequest{
		OrderAmount:        5000,
		OrderDescription:   "x",
		OrderInvoiceNumber: "INV-1",
		PaymentMethod:      PaymentBankTransfer,
		CustomerID:         "CUST_001",
	})
	if err != nil {
		t.Fatal(err)
	}
	if form.Action != "https://pay-sandbox.sepay.vn/v1/checkout/init" {
		t.Errorf("Action = %q", form.Action)
	}
	v := form.Values()
	if v.Get("payment_method") != "BANK_TRANSFER" || v.Get("customer_id") != "CUST_001" {
		t.Errorf("optional fields missing: %v", v)
	}
	if v.Has("success_url") {
		t.Error("empty success_url should be omitted")
	}
	if got, want := v.Get("signature"), Sign(form.Fields, "S"); got != want {
		t.Errorf("signature = %q, want %q", got, want)
	}
}

func TestNewCheckoutValidation(t *testing.T) {
	valid := func() CheckoutRequest { return *docRequest }
	tests := map[string]func(*CheckoutRequest){
		"no invoice":     func(r *CheckoutRequest) { r.OrderInvoiceNumber = "" },
		"no description": func(r *CheckoutRequest) { r.OrderDescription = "" },
		"zero amount":    func(r *CheckoutRequest) { r.OrderAmount = 0 },
		"negative":       func(r *CheckoutRequest) { r.OrderAmount = -1 },
	}
	for name, mutate := range tests {
		r := valid()
		mutate(&r)
		if _, err := New("M", "S").NewCheckout(&r); err == nil {
			t.Errorf("%s: want error", name)
		}
	}
	r := valid()
	r.Operation, r.OrderAmount = OperationVerify, 0
	if _, err := New("M", "S").NewCheckout(&r); err != nil {
		t.Errorf("VERIFY with zero amount: %v", err)
	}
	if _, err := New("", "S").NewCheckout(docRequest); err == nil {
		t.Error("empty merchant: want error")
	}
}

func TestCheckoutFormHTMLEscapes(t *testing.T) {
	form := &CheckoutForm{
		Action: "https://pay.sepay.vn/v1/checkout/init",
		Fields: []Field{{"order_description", `"><script>alert(1)</script>`}},
	}
	out := form.HTML("Pay & go")
	if strings.Contains(out, "<script>") {
		t.Errorf("value not escaped:\n%s", out)
	}
	for _, want := range []string{
		`<form action="https://pay.sepay.vn/v1/checkout/init" method="POST">`,
		`name="order_description" value="&#34;&gt;&lt;script&gt;`,
		`<button type="submit">Pay &amp; go</button>`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}
