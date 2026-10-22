package bankhub_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/oxycoder/sepay-go/bankhub"
)

func Example() {
	ctx := context.Background()
	client := bankhub.New(os.Getenv("SEPAY_CLIENT_ID"), os.Getenv("SEPAY_CLIENT_SECRET"), bankhub.WithSandbox())

	company, err := client.CreateCompany(ctx, &bankhub.CreateCompanyRequest{
		FullName: "Công ty TNHH ABC",
		Status:   bankhub.CompanyActive,
	})
	if err != nil {
		log.Fatal(err)
	}

	link, err := client.CreateLinkToken(ctx, &bankhub.CreateLinkTokenRequest{
		CompanyXID:            company.XID,
		Purpose:               bankhub.PurposeLink,
		CompletionRedirectURI: "https://example.com/bank-linked",
		Language:              bankhub.LanguageVI,
	})
	if err != nil {
		log.Fatal(err)
	}
	// Show this page to the user in an iframe or WebView.
	fmt.Println(link.HostedLinkURL)
}

func ExampleIPNHandler() {
	http.Handle("/sepay/ipn", bankhub.IPNHandler(os.Getenv("SEPAY_IPN_API_KEY"),
		func(ctx context.Context, ipn *bankhub.IPN) error {
			// SePay may deliver an IPN more than once: dedupe on TransactionID.
			if ipn.TransferType == bankhub.TransferCredit {
				log.Printf("received %.0f VND for %q", float64(ipn.Amount), ipn.PaymentCode)
			}
			return nil
		}))
}

func ExampleEventHandler() {
	http.Handle("/sepay/events", bankhub.EventHandler(os.Getenv("SEPAY_WEBHOOK_SECRET"),
		func(ctx context.Context, e *bankhub.Event) error {
			switch e.Event {
			case bankhub.EventBankAccountLinked:
				log.Printf("linked %s %s", e.Metadata.BrandName, e.Metadata.AccountNumber)
			case bankhub.EventBankAccountUnlinked, bankhub.EventBankAccountInactivated, bankhub.EventBankAccountInactived:
				log.Printf("bank account %s removed", e.Metadata.BankAccountXID)
			}
			return nil
		}))
}

func ExampleClient_UpsertWebhook() {
	client := bankhub.New(os.Getenv("SEPAY_CLIENT_ID"), os.Getenv("SEPAY_CLIENT_SECRET"))
	_, err := client.UpsertWebhook(context.Background(), &bankhub.UpsertWebhookRequest{
		WebhookURL:  "https://example.com/sepay/events",
		AuthType:    bankhub.AuthSecretKey,
		SecretKey:   os.Getenv("SEPAY_WEBHOOK_SECRET"),
		Active:      bankhub.Ptr(true),
		AllowEvents: []bankhub.EventType{bankhub.EventAll},
	})
	if err != nil {
		log.Fatal(err)
	}
}
