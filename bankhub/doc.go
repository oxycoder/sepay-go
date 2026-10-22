// Package bankhub is a client for the SePay Bank Hub API, which links bank
// accounts to your merchant, lists their transactions and notifies you of
// balance changes in real time.
//
// Create a client with your client_id and client_secret. The client fetches
// and caches the Bearer token itself and refreshes it when it expires:
//
//	client := bankhub.New(clientID, clientSecret, bankhub.WithSandbox())
//
// A typical integration creates a company, then a link token whose
// HostedLinkURL is opened in an iframe or WebView for the user to link a bank
// account:
//
//	company, err := client.CreateCompany(ctx, &bankhub.CreateCompanyRequest{
//		FullName: "Công ty TNHH ABC",
//		Status:   bankhub.CompanyActive,
//	})
//	link, err := client.CreateLinkToken(ctx, &bankhub.CreateLinkTokenRequest{
//		CompanyXID: company.XID,
//		Purpose:    bankhub.PurposeLink,
//	})
//	// open link.HostedLinkURL
//
// SePay pushes two kinds of notifications to your server. Use IPNHandler for
// balance changes (new transactions) and EventHandler for link lifecycle
// events (bank account linked, unlinked, ...):
//
//	mux.Handle("/sepay/ipn", bankhub.IPNHandler(apiKey, func(ctx context.Context, ipn *bankhub.IPN) error { ... }))
//	mux.Handle("/sepay/events", bankhub.EventHandler(secretKey, func(ctx context.Context, e *bankhub.Event) error { ... }))
//
// API errors are returned as *APIError:
//
//	var apiErr *bankhub.APIError
//	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusBadRequest { ... apiErr.Errors ... }
package bankhub
