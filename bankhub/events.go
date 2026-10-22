package bankhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EventType names a webhook event.
type EventType string

// EventAll subscribes to every event in UpsertWebhookRequest.AllowEvents.
const EventAll EventType = "*"

// Link token events.
const (
	EventLinkTokenCreated EventType = "LINK_TOKEN_CREATED"
	EventLinkTokenRevoked EventType = "LINK_TOKEN_REVOKED"
	EventLinkTokenExpired EventType = "LINK_TOKEN_EXPIRED"
)

// Link session events.
const (
	EventLinkSessionInitialized   EventType = "LINK_SESSION_INITIALIZED"
	EventLinkSessionReinitialized EventType = "LINK_SESSION_REINITIALIZED"
	EventLinkSessionStateChanged  EventType = "LINK_SESSION_STATE_CHANGED" // see EventMetadata.State
	EventLinkSessionCompleted     EventType = "LINK_SESSION_COMPLETED"
	EventLinkSessionFailed        EventType = "LINK_SESSION_FAILED"
	EventLinkSessionAbandoned     EventType = "LINK_SESSION_ABANDONED"
	EventLinkSessionExpired       EventType = "LINK_SESSION_EXPIRED"
)

// Bank account events.
const (
	EventBankAccountLinked      EventType = "BANK_ACCOUNT_LINKED"
	EventBankAccountUnlinked    EventType = "BANK_ACCOUNT_UNLINKED" // unlinked at the bank
	EventBankAccountInactivated EventType = "BANK_ACCOUNT_INACTIVATED"
	// EventBankAccountInactived is the spelling shown in SePay's event
	// payload examples; handle it together with EventBankAccountInactivated.
	// The account is deactivated in Bank Hub without unlinking at the bank.
	EventBankAccountInactived   EventType = "BANK_ACCOUNT_INACTIVED"
	EventBankAccountReactivated EventType = "BANK_ACCOUNT_REACTIVATED"
)

// LinkState is the step of a link session, sent in EventLinkSessionStateChanged.
type LinkState string

const (
	StateRequestLink        LinkState = "REQUEST_BANK_ACCOUNT_LINK"   // bank chosen, entering account details
	StateConfirmLink        LinkState = "CONFIRM_BANK_ACCOUNT_LINK"   // details entered, entering OTP
	StateFinishedLink       LinkState = "FINISHED_BANK_ACCOUNT_LINK"  // account linked
	StateRequestUnlink      LinkState = "REQUEST_BANK_ACCOUNT_UNLINK" // confirming unlink (OTP)
	StateFinishedUnlink     LinkState = "FINISHED_BANK_ACCOUNT_UNLINK"
	StateFinishedReactivate LinkState = "FINISHED_BANK_ACCOUNT_REACTIVATE"
)

// Event is a link lifecycle event pushed to the webhook configured with UpsertWebhook.
type Event struct {
	XID       string    `json:"xid"` // unique per event; use it to drop duplicates
	Timestamp int64     `json:"timestamp"`
	Version   string    `json:"version"`
	Event     EventType `json:"event"`
	// Metadata holds the fields known for any event type; fields an event
	// does not carry are left empty. RawMetadata has the original JSON.
	Metadata    EventMetadata   `json:"-"`
	RawMetadata json.RawMessage `json:"metadata"`
}

// Time returns the event's creation time.
func (e *Event) Time() time.Time { return time.Unix(e.Timestamp, 0) }

// EventMetadata is the union of the metadata fields of all event types.
type EventMetadata struct {
	LinkTokenXID      string      `json:"link_token_xid"`
	LinkSessionXID    string      `json:"link_session_xid"`
	Purpose           Purpose     `json:"purpose"`
	State             LinkState   `json:"state"`
	BankAccountXID    string      `json:"bank_account_xid"`
	BrandName         string      `json:"brand_name"`
	AccountType       AccountType `json:"account_type"`
	AccountNumber     string      `json:"account_number"`
	AccountHolderName string      `json:"account_holder_name"`
	Unlinked          bool        `json:"unlinked"`    // set on LINK_SESSION_COMPLETED of an unlink flow
	Reactivated       bool        `json:"reactivated"` // set on LINK_SESSION_COMPLETED of a reactivate flow
}

// ParseEvent decodes a webhook event body.
func ParseEvent(body []byte) (*Event, error) {
	var e Event
	if err := json.Unmarshal(body, &e); err != nil {
		return nil, fmt.Errorf("bankhub: decode event: %w", err)
	}
	if len(e.RawMetadata) > 0 && string(e.RawMetadata) != "null" {
		if err := json.Unmarshal(e.RawMetadata, &e.Metadata); err != nil {
			return nil, fmt.Errorf("bankhub: decode event metadata: %w", err)
		}
	}
	return &e, nil
}

// EventHandler returns an http.Handler for the event webhook. It checks that
// the X-Secret-Key header equals secretKey (skipped if secretKey is empty),
// decodes the event and calls fn. If fn returns an error the handler responds
// 500; SePay only retries on network failures, so do not rely on redelivery.
func EventHandler(secretKey string, fn func(context.Context, *Event) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if secretKey != "" && !secureEqual(r.Header.Get("X-Secret-Key"), secretKey) {
			respond(w, http.StatusUnauthorized)
			return
		}
		body, ok := readBody(w, r)
		if !ok {
			return
		}
		e, err := ParseEvent(body)
		if err != nil {
			respond(w, http.StatusBadRequest)
			return
		}
		if err := fn(r.Context(), e); err != nil {
			respond(w, http.StatusInternalServerError)
			return
		}
		respond(w, http.StatusOK)
	})
}
