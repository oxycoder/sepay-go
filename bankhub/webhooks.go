package bankhub

import (
	"context"
	"encoding/json"
	"net/http"
)

// AuthSecretKey makes SePay send the webhook secret in the X-Secret-Key header.
const AuthSecretKey = "SECRET_KEY"

// WebhookConfig is the merchant's event webhook configuration.
type WebhookConfig struct {
	WebhookURL string `json:"webhook_url"`
	AuthType   string `json:"auth_type"`
	// SecretKey is masked ("***5678") by GetWebhook and returned in full by UpsertWebhook.
	SecretKey          string      `json:"secret_key"`
	RequestContentType string      `json:"request_content_type"`
	Active             Flag        `json:"active"`
	AllowEvents        []EventType `json:"allow_events"`
}

// GetWebhook returns the event webhook configuration, or nil if none is configured.
func (c *Client) GetWebhook(ctx context.Context) (*WebhookConfig, error) {
	raw, err := c.call(ctx, http.MethodGet, "/v1/webhook", nil, nil)
	if err != nil {
		return nil, err
	}
	var out WebhookConfig
	found, err := decodeData(raw, &out)
	if err != nil || !found {
		return nil, err
	}
	return &out, nil
}

// UpsertWebhookRequest is the body of UpsertWebhook. Unset fields are left unchanged.
type UpsertWebhookRequest struct {
	// WebhookURL must be public HTTPS in production, at most 2000 characters.
	WebhookURL string
	AuthType   string // AuthSecretKey
	SecretKey  string // at most 500 characters; 32-64 random characters recommended
	Active     *bool
	// AllowEvents selects the events to send: nil leaves it unchanged, an
	// empty non-nil slice sends no events, []EventType{EventAll} sends all.
	AllowEvents []EventType
}

func (r *UpsertWebhookRequest) MarshalJSON() ([]byte, error) {
	w := struct {
		WebhookURL  string       `json:"webhook_url,omitempty"`
		AuthType    string       `json:"auth_type,omitempty"`
		SecretKey   string       `json:"secret_key,omitempty"`
		Active      *Flag        `json:"active,omitempty"`
		AllowEvents *[]EventType `json:"allow_events,omitempty"`
	}{
		WebhookURL: r.WebhookURL,
		AuthType:   r.AuthType,
		SecretKey:  r.SecretKey,
	}
	if r.Active != nil {
		w.Active = Ptr(Flag(*r.Active))
	}
	if r.AllowEvents != nil {
		w.AllowEvents = &r.AllowEvents
	}
	return json.Marshal(w)
}

// UpsertWebhook creates the event webhook, or updates it if it exists. A
// merchant has at most one. Receive the events with EventHandler.
func (c *Client) UpsertWebhook(ctx context.Context, req *UpsertWebhookRequest) (*WebhookConfig, error) {
	var out WebhookConfig
	if err := c.get(ctx, http.MethodPost, "/v1/webhook", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
