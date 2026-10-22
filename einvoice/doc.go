// Package einvoice is a client for the SePay eInvoice API, which creates,
// issues and manages Vietnamese electronic invoices (Nghị định 123/2020/NĐ-CP).
//
// Create a client with your client_id and client_secret. The client fetches
// and caches the Bearer token itself and refreshes it when it expires:
//
//	client := einvoice.New(clientID, clientSecret, einvoice.WithSandbox())
//
// Creating and issuing invoices is asynchronous. CreateInvoice, AdjustInvoice,
// ReplaceInvoice and IssueInvoice return a tracking code; use WaitForCreate or
// WaitForIssue (or poll GetCreateStatus / GetIssueStatus yourself) until the
// status is Success or Failed:
//
//	res, err := client.CreateInvoice(ctx, &einvoice.CreateInvoiceRequest{...})
//	if err != nil { ... }
//	status, err := client.WaitForCreate(ctx, res.TrackingCode)
//
// API errors are returned as *APIError. Branch on APIError.Code, never on the
// message text:
//
//	var apiErr *einvoice.APIError
//	if errors.As(err, &apiErr) && apiErr.Code == einvoice.CodeDocumentExisted { ... }
package einvoice
