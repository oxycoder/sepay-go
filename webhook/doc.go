// Package webhook receives SePay Webhooks: the HTTP POST SePay sends to your
// server for every new transaction on a linked bank account.
//
// Mount Handler on the URL configured at my.sepay.vn → Webhooks, choosing the
// option that matches the webhook's authentication setting:
//
//	mux.Handle("/sepay/webhook", webhook.Handler(
//		func(ctx context.Context, tx *webhook.Transaction) error {
//			// SePay can send the same transaction more than once: dedupe on tx.ID.
//			return orders.MarkPaid(ctx, tx.ID, tx.Code, tx.TransferAmount)
//		},
//		webhook.WithHMAC(os.Getenv("SEPAY_WEBHOOK_SECRET")), // or WithAPIKey, WithBearer
//	))
//
// The handler checks the credential, decodes JSON, form-urlencoded or
// multipart bodies, calls your function and answers {"success": true}. When
// your function returns an error it answers 500 and SePay retries the delivery
// up to 7 times over about 33 minutes. SePay waits at most 30 seconds for a
// response, so hand slow work to a queue.
//
// If you use another router or framework, read the raw body yourself and call
// VerifyAPIKey or VerifySignature, then Parse.
package webhook
