package paymentgateway_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/oxycoder/sepay-go/paymentgateway"
)

func Example() {
	client := paymentgateway.New(os.Getenv("SEPAY_MERCHANT_ID"), os.Getenv("SEPAY_SECRET_KEY"), paymentgateway.WithSandbox())

	http.HandleFunc("/checkout", func(w http.ResponseWriter, r *http.Request) {
		form, err := client.NewCheckout(&paymentgateway.CheckoutRequest{
			OrderInvoiceNumber: "INV-1001",
			OrderAmount:        100000,
			OrderDescription:   "Thanh toán đơn hàng #1001",
			SuccessURL:         "https://example.com/payment/success",
			ErrorURL:           "https://example.com/payment/error",
			CancelURL:          "https://example.com/payment/cancel",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, form.HTML("Thanh toán"))
	})
}

func ExampleIPNHandler() {
	http.Handle("/sepay/ipn", paymentgateway.IPNHandler(os.Getenv("SEPAY_SECRET_KEY"),
		func(ctx context.Context, ipn *paymentgateway.IPN) error {
			// SePay may deliver an IPN more than once: make this idempotent.
			switch ipn.NotificationType {
			case paymentgateway.NotificationOrderPaid:
				log.Printf("order %s paid: %.0f VND", ipn.Order.OrderInvoiceNumber, float64(ipn.Order.OrderAmount))
			case paymentgateway.NotificationTransactionVoid:
				log.Printf("order %s voided", ipn.Order.OrderInvoiceNumber)
			}
			return nil
		}))
}

func ExampleClient_ListOrders() {
	client := paymentgateway.New(os.Getenv("SEPAY_MERCHANT_ID"), os.Getenv("SEPAY_SECRET_KEY"))

	list, err := client.ListOrders(context.Background(), &paymentgateway.ListOrdersParams{
		OrderStatus: paymentgateway.OrderCaptured,
		Sort:        paymentgateway.SortCreatedDesc,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, o := range list.Data {
		fmt.Println(o.OrderInvoiceNumber, o.OrderStatus, float64(o.OrderAmount))
	}
}

func ExampleSign() {
	fields := []paymentgateway.Field{
		{Name: "order_amount", Value: "100000"},
		{Name: "merchant", Value: "MERCHANT_123"},
		{Name: "currency", Value: "VND"},
		{Name: "operation", Value: "PURCHASE"},
		{Name: "order_description", Value: "Thanh toán đơn hàng #12345"},
		{Name: "order_invoice_number", Value: "INV_20231201_001"},
		{Name: "success_url", Value: "https://yoursite.com/payment/success"},
		{Name: "error_url", Value: "https://yoursite.com/payment/error"},
		{Name: "cancel_url", Value: "https://yoursite.com/payment/cancel"},
	}
	fmt.Println(paymentgateway.Sign(fields, "SECRET_KEY"))
	// Output: LddMASma1ObPVFJzgklR1Jye6/mv66LwWNBfydEsTQI=
}
