// Package paymentgateway is a client for the SePay Payment Gateway (Cổng thanh
// toán), which accepts VietQR bank transfers, NAPAS QR and international cards
// (Visa, Mastercard, JCB) through a hosted checkout page.
//
// Create a client with your merchant_id and secret_key. Sandbox and production
// use different credentials:
//
//	client := paymentgateway.New(merchantID, secretKey, paymentgateway.WithSandbox())
//
// A payment starts with a signed HTML form that the customer's browser POSTs
// to SePay's checkout page. NewCheckout fills in the merchant, signs the
// fields and returns the form:
//
//	form, err := client.NewCheckout(&paymentgateway.CheckoutRequest{
//		OrderInvoiceNumber: "INV-1001", // unique per order
//		OrderAmount:        100000,     // VND
//		OrderDescription:   "Thanh toán đơn hàng #1001",
//		SuccessURL:         "https://example.com/payment/success",
//		ErrorURL:           "https://example.com/payment/error",
//		CancelURL:          "https://example.com/payment/cancel",
//	})
//	io.WriteString(w, form.HTML("Thanh toán"))
//
// The customer is redirected back to SuccessURL, ErrorURL or CancelURL, but
// may close the browser first. Confirm payments from the IPN (Instant Payment
// Notification) that SePay POSTs to the URL configured at my.sepay.vn:
//
//	mux.Handle("/sepay/ipn", paymentgateway.IPNHandler(secretKey, func(ctx context.Context, ipn *paymentgateway.IPN) error {
//		if ipn.NotificationType == paymentgateway.NotificationOrderPaid {
//			return orders.MarkPaid(ctx, ipn.Order.OrderInvoiceNumber)
//		}
//		return nil
//	}))
//
// The REST API (ListOrders, GetOrder, CancelOrder, VoidTransaction) uses HTTP
// Basic authentication with the same credentials. API errors are returned as
// *APIError:
//
//	var apiErr *paymentgateway.APIError
//	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnprocessableEntity { ... }
package paymentgateway
