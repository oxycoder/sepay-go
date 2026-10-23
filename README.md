# sepay-go

Go client for the [SePay](https://sepay.vn) APIs.

| Package | API |
| ------- | --- |
| [`einvoice`](./einvoice) | SePay eInvoice API: create, issue, adjust and replace Vietnamese e-invoices |
| [`bankhub`](./bankhub) | SePay Bank Hub API: link bank accounts, query transactions, receive balance notifications |
| [`webhook`](./webhook) | SePay Webhooks: receive and verify transaction notifications for your linked bank accounts |
| [`paymentgateway`](./paymentgateway) | SePay Payment Gateway: signed checkout forms (VietQR, NAPAS QR, cards), IPN, order lookup, cancel and void |

No dependencies outside the standard library. Requires Go 1.22+.

```sh
go get github.com/oxycoder/sepay-go
```

## eInvoice

```go
import "github.com/oxycoder/sepay-go/einvoice"

client := einvoice.New(clientID, clientSecret)                         // production
client := einvoice.New(clientID, clientSecret, einvoice.WithSandbox()) // sandbox
```

The client gets the Bearer token (`POST /v1/token`), caches it for its 24h lifetime, and fetches a new one when the API returns `401`. It is safe for concurrent use; create one per credential set and reuse it.

Options: `WithSandbox`, `WithBaseURL`, `WithHTTPClient`, `WithUserAgent`, `WithPollInterval`.

### Create and issue an invoice

```go
res, err := client.CreateInvoice(ctx, &einvoice.CreateInvoiceRequest{
	TemplateCode:      tpl.TemplateCode,  // from client.GetProviderAccount(...).Templates
	InvoiceSeries:     tpl.InvoiceSeries, // don't hardcode: must match the current year (C26 = 2026)
	Currency:          einvoice.CurrencyVND,
	ProviderAccountID: accountID,
	ReferenceCode:     "ORDER-1001", // optional, unique; lets you check before resending
	IsDraft:           false,        // true = draft, issue later with client.IssueInvoice
	Buyer:             einvoice.Buyer{Name: "Công ty ABC", TaxCode: "0101234567", Email: "buyer@example.com"},
	Items: []einvoice.Item{{
		LineNumber: 1,
		LineType:   einvoice.LineTypeGoods,
		ItemName:   "Sản phẩm A",
		Unit:       "cái",
		Quantity:   1,
		UnitPrice:  4500000,
		TaxRate:    einvoice.Ptr(einvoice.TaxRate10),
	}},
})
if err != nil { ... }

status, err := client.WaitForCreate(ctx, res.TrackingCode) // polls until Success or Failed
```

`IssuedDate` defaults to now and is always sent in Vietnam time (`YYYY-MM-DD HH:mm:ss`). SePay only returns the date, so the exact time sent is kept on `res.IssuedDate` for create, adjust and replace; store it if you need the time. Optional fields where zero means something (`TaxRate`, `TaxAmount`, `TotalAmount`, `SubAmount`, `Amount`) are pointers. Set them with `einvoice.Ptr`.

### Asynchronous processing

Create, adjust, replace and issue requests are queued. When SePay hits a temporary error it retries for up to about 2 hours, so a request can stay `Pending` that long.

- `WaitForCreate` and `WaitForIssue` poll every 3s by default. When the response says when the next retry is (`retry.next_retry_at`), they wait until then. Use a `ctx` deadline to stop them.
- A `Failed` status is a result, not an error. Check `status.Status`.
- If `status.Unconfirmed()` is true, SePay could not confirm the outcome with the invoice provider. **Do not resend.** Contact SePay with the tracking code.
- Before you resend a failed invoice, call `client.GetInvoice(ctx, referenceCode)`. If `einvoice.IsNotFound(err)` is true, no invoice was recorded and it is safe to resend.

### Errors

API errors are `*einvoice.APIError` with `StatusCode`, `Code` and `Message`. Branch on `Code` and never on `Message`:

```go
var apiErr *einvoice.APIError
if errors.As(err, &apiErr) {
	switch apiErr.Code {
	case einvoice.CodeDocumentExisted: // reference_code already used
	case einvoice.CodeQuotaUsedUp:     // renew quota, don't retry
	}
}
// shortcuts
einvoice.HasCode(err, einvoice.CodeSellerStoreInactive)
einvoice.IsNotFound(err)
```

### Endpoints

| Method | Endpoint |
| ------ | -------- |
| `ListProviderAccounts`, `GetProviderAccount` | `GET /v1/provider-accounts[/{id}]` |
| `CreateInvoice` | `POST /v1/invoices/create` |
| `AdjustInvoice` | `POST /v1/invoices/adjust` |
| `ReplaceInvoice` | `POST /v1/invoices/replace` |
| `GetCreateStatus`, `WaitForCreate` | `GET /v1/invoices/create/check/{tracking_code}` (create, adjust and replace) |
| `IssueInvoice` | `POST /v1/invoices/issue` |
| `GetIssueStatus`, `WaitForIssue` | `GET /v1/invoices/issue/check/{tracking_code}` |
| `GetInvoice` | `GET /v1/invoices/{reference_code}` |
| `ListInvoices` | `GET /v1/invoices` |
| `DownloadInvoice` | `GET /v1/invoices/{tracking_code}/download` (content already base64-decoded) |
| `DeleteDraftInvoice` | `POST /v1/invoices/delete/{reference_code}` |
| `GetUsage` | `GET /v1/usage` |
| `ListEmailLogs` | `GET /v1/email-logs` |
| `ListBusinessLocations`, `CreateBusinessLocation`, `UpdateBusinessLocation`, `DeleteBusinessLocation` | `/v1/business-locations[/create\|/update\|/delete]` |

The API reference (Vietnamese) is in [docs/einvoice](./docs/einvoice).

## Bank Hub

```go
import "github.com/oxycoder/sepay-go/bankhub"

client := bankhub.New(clientID, clientSecret)                        // production
client := bankhub.New(clientID, clientSecret, bankhub.WithSandbox()) // sandbox
```

Token handling is the same as in `einvoice`. The client gets a token from `POST /v1/token`, caches it for its `ttl`, and fetches a new one when the API returns `401`.

Options: `WithSandbox`, `WithBaseURL`, `WithHTTPClient`, `WithUserAgent`.

### Link a bank account

```go
company, err := client.CreateCompany(ctx, &bankhub.CreateCompanyRequest{
	FullName: "Công ty TNHH ABC",
	Status:   bankhub.CompanyActive,
})

link, err := client.CreateLinkToken(ctx, &bankhub.CreateLinkTokenRequest{
	CompanyXID: company.XID,
	Purpose:    bankhub.PurposeLink, // PurposeUnlink / PurposeReactivate need BankAccountXID
})
// Open link.HostedLinkURL in an iframe or WebView.
```

### Receive notifications

SePay sends two kinds of notifications. Both handlers check the credential with a constant-time compare, reply in the format SePay expects, and return `500` when your function returns an error.

```go
// Balance changes (new transactions). SePay sends "Authorization: Apikey {key}".
mux.Handle("/sepay/ipn", bankhub.IPNHandler(apiKey, func(ctx context.Context, ipn *bankhub.IPN) error {
	// SePay can send the same IPN more than once. Skip IDs you have already processed.
	return processPayment(ctx, ipn.TransactionID, ipn.PaymentCode, float64(ipn.Amount))
}))

// Link lifecycle events. SePay sends the secret from UpsertWebhook in X-Secret-Key.
mux.Handle("/sepay/events", bankhub.EventHandler(secretKey, func(ctx context.Context, e *bankhub.Event) error {
	if e.Event == bankhub.EventBankAccountLinked { /* e.Metadata.BankAccountXID */ }
	return nil
}))
```

If you check the request yourself, use `bankhub.ParseIPN` and `bankhub.ParseEvent`. In the sandbox, `client.SimulateTransaction` sends a test IPN.

A company receives IPNs only when its status is `Active`, `active` is `1`, and its `transaction_amount` is above 0 or `Unlimited`. Set these with `UpdateCompany`.

### Errors

API errors are `*bankhub.APIError` with `StatusCode`, `Message` and, for validation errors, a per-field `Errors` map. Use `bankhub.IsNotFound(err)` and `bankhub.IsValidationError(err)` to check for common cases.

### Endpoints

| Method | Endpoint |
| ------ | -------- |
| `CreateLinkToken`, `GetLinkToken`, `RevokeLinkToken` | `POST /v1/link-token/{create,get,revoke}` |
| `ListCompanies`, `GetCompany` | `GET /v1/company[/{xid}]` |
| `CreateCompany` | `POST /v1/company/create` |
| `UpdateCompany` | `POST /v1/company/edit/{xid}` |
| `GetCompanyCounter` | `GET /v1/company/counter/{xid}` |
| `ListBankAccounts`, `GetBankAccount` | `GET /v1/bank-account[/{xid}]` |
| `ListTransactions`, `GetTransaction` | `GET /v1/transaction[/{transaction_id}]` |
| `SimulateTransaction` | `POST /v1/transaction/create` (sandbox only) |
| `GetMerchantConfig`, `UpdateMerchantConfig` | `GET`/`POST /v1/merchant-config` |
| `UploadMerchantLogo` | `POST /v1/merchant-config/logo` (multipart, PNG/JPEG ≤ 1MB) |
| `GetMerchantCounter` | `GET /v1/merchant-config/counter` |
| `GetWebhook`, `UpsertWebhook` | `GET`/`POST /v1/webhook` |

The API reference (Vietnamese) is in [docs/bankhub](./docs/bankhub).

## Webhooks

```go
import "github.com/oxycoder/sepay-go/webhook"
```

SePay sends an HTTP POST to your URL for every new transaction on a linked bank account. `webhook.Handler` checks the credential, decodes the body (JSON, form-urlencoded or multipart), calls your function and replies `{"success": true}`.

```go
mux.Handle("/sepay/webhook", webhook.Handler(
	func(ctx context.Context, tx *webhook.Transaction) error {
		// SePay can send the same transaction more than once (retries, manual replays). Dedupe on tx.ID.
		if tx.TransferType == webhook.TransferIn && tx.Code != "" {
			return orders.MarkPaid(ctx, tx.ID, tx.Code, tx.TransferAmount)
		}
		return nil
	},
	webhook.WithHMAC(os.Getenv("SEPAY_WEBHOOK_SECRET")),
))
```

Pick the option that matches the authentication set on the webhook:

| Authentication | Option | What is checked |
| -------------- | ------ | --------------- |
| HMAC-SHA256 (recommended) | `WithHMAC(secret)` | `X-SePay-Signature` = `sha256=HMAC(secret, "{timestamp}.{raw body}")`, and `X-SePay-Timestamp` within 5 minutes (`WithTolerance` to change) |
| API Key | `WithAPIKey(key)` | `Authorization: Apikey {key}` |
| OAuth 2.0 | `WithBearer(func(ctx, token) bool)` | `Authorization: Bearer {token}` accepted by your function |
| None | no option | nothing; local testing only |

Responses: `200` when your function returns nil; `401` for a bad credential, `400` for an unreadable body and `500` when your function returns an error. SePay retries any non-2xx response up to 7 times over about 33 minutes, and waits at most 30 seconds for a reply, so push slow work to a queue.

`tx.Raw` holds the body exactly as received; store it for audits. `tx.Time()` parses `TransactionDate` in Vietnam time. With another router, call `webhook.VerifySignature` (or `VerifyAPIKey`) on the raw body, then `webhook.Parse`. `webhook.Sign` builds a valid signature for your tests.

The guide (Vietnamese) is in [docs/webhook](./docs/webhook).

## Payment Gateway

```go
import "github.com/oxycoder/sepay-go/paymentgateway"

client := paymentgateway.New(merchantID, secretKey)                               // production
client := paymentgateway.New(merchantID, secretKey, paymentgateway.WithSandbox()) // sandbox
```

`WithSandbox` switches both hosts: the REST API (`pgapi-sandbox.sepay.vn`) and the checkout page (`pay-sandbox.sepay.vn`). Sandbox and production use different `merchant_id`/`secret_key` pairs.

Options: `WithSandbox`, `WithBaseURL`, `WithCheckoutURL`, `WithHTTPClient`, `WithUserAgent`.

### Send the customer to checkout

The checkout is an HTML form the customer's browser POSTs to `/v1/checkout/init`. `NewCheckout` validates the request, fills in `merchant`, `currency` (VND) and `operation` (PURCHASE), and signs the fields in the order SePay expects. It makes no network call.

```go
form, err := client.NewCheckout(&paymentgateway.CheckoutRequest{
	OrderInvoiceNumber: "INV-1001", // unique; IPNs, CancelOrder and VoidTransaction refer to it
	OrderAmount:        100000,     // VND
	OrderDescription:   "Thanh toán đơn hàng #1001",
	PaymentMethod:      paymentgateway.PaymentBankTransfer, // optional; empty lets the customer choose
	SuccessURL:         "https://example.com/payment/success",
	ErrorURL:           "https://example.com/payment/error",
	CancelURL:          "https://example.com/payment/cancel",
})
io.WriteString(w, form.HTML("Thanh toán")) // or render form.Action + form.Fields in your own template, in order
```

If you build the form yourself, `paymentgateway.Sign(fields, secretKey)` computes the signature.

### Confirm payments with the IPN

Don't mark an order paid from `SuccessURL`: the customer may close the browser before the redirect. SePay POSTs an IPN to the URL set at my.sepay.vn → Cổng thanh toán → Cấu hình → IPN.

```go
mux.Handle("/sepay/ipn", paymentgateway.IPNHandler(secretKey, func(ctx context.Context, ipn *paymentgateway.IPN) error {
	// SePay can send the same IPN more than once: make this idempotent.
	if ipn.NotificationType == paymentgateway.NotificationOrderPaid {
		return orders.MarkPaid(ctx, ipn.Order.OrderInvoiceNumber, float64(ipn.Order.OrderAmount))
	}
	return nil
}))
```

The handler checks `X-Secret-Key` with a constant-time compare (pass `""` if the IPN auth type is not `SECRET_KEY`) and answers `{"success": true}`. It returns `401` for a bad key, `400` for an unreadable body and `500` when your function returns an error. With another router, call `paymentgateway.VerifyIPN` and `paymentgateway.ParseIPN`.

### Errors

API errors are `*paymentgateway.APIError` with `StatusCode`, `Message` and, for validation errors, a per-field `Errors` map. Use `paymentgateway.IsNotFound`, `IsValidationError` (400/422) and `IsRateLimited` (429) to check for common cases. Times such as `Order.CreatedAt` are Vietnam time; parse them with `paymentgateway.ParseTime`.

### Endpoints

| Method | Endpoint |
| ------ | -------- |
| `NewCheckout`, `Sign` | form `POST {checkout}/v1/checkout/init` (no API call) |
| `ListOrders` | `GET /v1/order` |
| `GetOrder` | `GET /v1/order/detail/{order_id}` |
| `CancelOrder` | `POST /v1/order/cancel` (unpaid bank transfer / NAPAS orders) |
| `VoidTransaction` | `POST /v1/order/voidTransaction` (captured card payments, before settlement) |

The API reference (Vietnamese) is in [docs/payment-gateway](./docs/payment-gateway).
