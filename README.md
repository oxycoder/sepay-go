# sepay-go

Go client for the [SePay](https://sepay.vn) APIs.

| Package | API |
| ------- | --- |
| [`einvoice`](./einvoice) | SePay eInvoice API: create, issue, adjust and replace Vietnamese e-invoices |
| [`bankhub`](./bankhub) | SePay Bank Hub API: link bank accounts, query transactions, receive balance notifications |

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

`IssuedDate` defaults to now and is always sent in Vietnam time (`YYYY-MM-DD HH:mm:ss`). Optional fields where zero means something (`TaxRate`, `TaxAmount`, `TotalAmount`, `SubAmount`, `Amount`) are pointers. Set them with `einvoice.Ptr`.

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
