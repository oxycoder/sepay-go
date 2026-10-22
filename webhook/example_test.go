package webhook_test

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/oxycoder/sepay-go/webhook"
)

func ExampleHandler() {
	http.Handle("/sepay/webhook", webhook.Handler(
		func(ctx context.Context, tx *webhook.Transaction) error {
			// SePay may deliver a transaction more than once: dedupe on tx.ID,
			// e.g. with a UNIQUE column and INSERT IGNORE.
			if tx.TransferType == webhook.TransferIn && tx.Code != "" {
				log.Printf("order %s paid %d VND", tx.Code, tx.TransferAmount)
			}
			return nil
		},
		webhook.WithHMAC(os.Getenv("SEPAY_WEBHOOK_SECRET")),
	))
}

func ExampleVerifySignature() {
	// With a framework that hands you the raw body, verify and parse yourself.
	http.HandleFunc("/sepay/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}
		if err := webhook.VerifySignature(r.Header, body, os.Getenv("SEPAY_WEBHOOK_SECRET"), webhook.DefaultTolerance); err != nil {
			http.Error(w, `{"success":false}`, http.StatusUnauthorized)
			return
		}
		tx, err := webhook.Parse(r.Header.Get("Content-Type"), body)
		if err != nil {
			http.Error(w, `{"success":false}`, http.StatusBadRequest)
			return
		}
		log.Printf("transaction %d", tx.ID)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"success":true}`)
	})
}
