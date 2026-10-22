package einvoice_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oxycoder/sepay-go/einvoice"
)

func Example() {
	ctx := context.Background()
	client := einvoice.New(os.Getenv("SEPAY_CLIENT_ID"), os.Getenv("SEPAY_CLIENT_SECRET"), einvoice.WithSandbox())

	// Pick a provider account and a template/series valid for this year.
	accounts, err := client.ListProviderAccounts(ctx, nil)
	if err != nil || len(accounts.Items) == 0 {
		log.Fatal("no provider account: ", err)
	}
	account, err := client.GetProviderAccount(ctx, accounts.Items[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	tpl := account.Templates[0]

	// Create and issue directly (IsDraft: false).
	res, err := client.CreateInvoice(ctx, &einvoice.CreateInvoiceRequest{
		TemplateCode:      tpl.TemplateCode,
		InvoiceSeries:     tpl.InvoiceSeries,
		Currency:          einvoice.CurrencyVND,
		ProviderAccountID: account.ID,
		ReferenceCode:     "ORDER-1001", // your own unique id makes resends safe
		PaymentMethod:     einvoice.PaymentBankTransfer,
		Buyer: einvoice.Buyer{
			Type:    einvoice.BuyerCompany,
			Name:    "Công ty ABC",
			TaxCode: "0101234567",
			Email:   "buyer@example.com",
		},
		Items: []einvoice.Item{{
			LineNumber: 1,
			LineType:   einvoice.LineTypeGoods,
			ItemCode:   "SP001",
			ItemName:   "Sản phẩm A",
			Unit:       "cái",
			Quantity:   1,
			UnitPrice:  4500000,
			TaxRate:    einvoice.Ptr(einvoice.TaxRate10),
		}},
	})
	if err != nil {
		var apiErr *einvoice.APIError
		if errors.As(err, &apiErr) && apiErr.Code == einvoice.CodeDocumentExisted {
			log.Fatal("invoice ORDER-1001 already exists")
		}
		log.Fatal(err)
	}

	// Processing is asynchronous and may stay Pending for up to 2 hours.
	waitCtx, cancel := context.WithTimeout(ctx, 2*time.Hour)
	defer cancel()
	status, err := client.WaitForCreate(waitCtx, res.TrackingCode)
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case status.Status == einvoice.StatusSuccess:
		fmt.Println("issued:", status.Invoice.InvoiceNumber, status.Invoice.PDFURL)
	case status.Unconfirmed():
		log.Fatal("outcome unknown, do not resend; contact SePay with ", res.TrackingCode)
	default:
		log.Fatal("failed: ", status.Message)
	}
}
