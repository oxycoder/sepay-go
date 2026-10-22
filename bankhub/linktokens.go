package bankhub

import (
	"context"
	"encoding/json"
	"net/http"
)

// Purpose of a link token.
type Purpose string

const (
	PurposeLink       Purpose = "LINK_BANK_ACCOUNT"       // link a new bank account
	PurposeUnlink     Purpose = "UNLINK_BANK_ACCOUNT"     // unlink (deactivate) an account
	PurposeReactivate Purpose = "REACTIVATE_BANK_ACCOUNT" // reactivate a deactivated account
)

// Language of the hosted link UI.
type Language string

const (
	LanguageVI Language = "vi"
	LanguageEN Language = "en"
)

// LinkTokenStatus is the state of a link token.
type LinkTokenStatus string

const (
	LinkTokenActive  LinkTokenStatus = "Active"
	LinkTokenRevoked LinkTokenStatus = "Revoked"
	LinkTokenExpired LinkTokenStatus = "Expired"
)

// SessionState is the state of a link session in LinkTokenDetail.Sessions.
type SessionState string

const (
	SessionPending    SessionState = "PENDING"
	SessionInProgress SessionState = "IN_PROGRESS"
	SessionCompleted  SessionState = "COMPLETED"
	SessionFailed     SessionState = "FAILED"
)

// CreateLinkTokenRequest is the body of CreateLinkToken.
type CreateLinkTokenRequest struct {
	CompanyXID string  // required
	Purpose    Purpose // required
	// BankAccountXID is required for PurposeUnlink and PurposeReactivate.
	BankAccountXID string
	// CompletionRedirectURI is where the user is sent after finishing the flow.
	CompletionRedirectURI string
	IsMobileApp           bool // true when opened from a mobile app, false for web
	Language              Language
}

func (r *CreateLinkTokenRequest) MarshalJSON() ([]byte, error) {
	w := struct {
		CompanyXID            string   `json:"company_xid"`
		Purpose               Purpose  `json:"purpose"`
		BankAccountXID        string   `json:"bank_account_xid,omitempty"`
		CompletionRedirectURI string   `json:"completion_redirect_uri,omitempty"`
		IsMobileApp           *Flag    `json:"is_mobile_app,omitempty"`
		Language              Language `json:"language,omitempty"`
	}{
		CompanyXID:            r.CompanyXID,
		Purpose:               r.Purpose,
		BankAccountXID:        r.BankAccountXID,
		CompletionRedirectURI: r.CompletionRedirectURI,
		Language:              r.Language,
	}
	if r.IsMobileApp {
		w.IsMobileApp = Ptr(Flag(true))
	}
	return json.Marshal(w)
}

// LinkToken is a newly created link token.
type LinkToken struct {
	XID string `json:"xid"` // pass to RevokeLinkToken
	// HostedLinkURL is the page to open in an iframe, WebView or the SDK.
	HostedLinkURL string `json:"hosted_link_url"`
	LinkToken     string `json:"link_token"` // pass to GetLinkToken
	ExpiresAt     string `json:"expires_at"` // "YYYY-MM-DD HH:mm:ss"
}

// LinkTokenDetail is a link token with its sessions and lifecycle events.
type LinkTokenDetail struct {
	XID                   string          `json:"xid"`
	Purpose               Purpose         `json:"purpose"`
	Status                LinkTokenStatus `json:"status"`
	CompletionRedirectURI string          `json:"completion_redirect_uri"`
	IsMobileApp           Flag            `json:"is_mobile_app"`
	Language              Language        `json:"language"`
	SessionCount          int             `json:"session_count"`
	CreatedAt             string          `json:"created_at"`
	UpdatedAt             string          `json:"updated_at"`
	ExpiresAt             string          `json:"expires_at"`
	Company               struct {
		XID      string        `json:"xid"`
		FullName string        `json:"full_name"`
		Status   CompanyStatus `json:"status"`
	} `json:"company"`
	Sessions []LinkSession    `json:"sessions"`
	Events   []LinkTokenEvent `json:"events"`
}

// LinkSession is one visit of a user to a link token.
type LinkSession struct {
	XID        string          `json:"xid"`
	State      SessionState    `json:"state"`
	Status     string          `json:"status"`
	Context    json.RawMessage `json:"context"`
	Result     json.RawMessage `json:"result"`
	FinishedAt string          `json:"finished_at"`
	ExpiresAt  string          `json:"expires_at"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
}

// LinkTokenEvent is a lifecycle event of a link token.
type LinkTokenEvent struct {
	XID       string          `json:"xid"`
	Name      EventType       `json:"name"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt string          `json:"created_at"`
}

// CreateLinkToken creates a one-time link for the user to link, unlink or
// reactivate a bank account. Open the returned HostedLinkURL for the user.
func (c *Client) CreateLinkToken(ctx context.Context, req *CreateLinkTokenRequest) (*LinkToken, error) {
	var out LinkToken
	if err := c.get(ctx, http.MethodPost, "/v1/link-token/create", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetLinkToken returns a link token, by its LinkToken value, with its sessions and events.
func (c *Client) GetLinkToken(ctx context.Context, linkToken string) (*LinkTokenDetail, error) {
	var out LinkTokenDetail
	in := map[string]string{"link_token": linkToken}
	if err := c.get(ctx, http.MethodPost, "/v1/link-token/get", nil, in, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// RevokeLinkToken invalidates a link token, by its XID, immediately. Users on
// the hosted link are disconnected. This cannot be undone.
func (c *Client) RevokeLinkToken(ctx context.Context, linkTokenXID string) error {
	in := map[string]string{"link_token_xid": linkTokenXID}
	_, err := c.call(ctx, http.MethodPost, "/v1/link-token/revoke", nil, in)
	return err
}
