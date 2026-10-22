package einvoice

import (
	"context"
	"net/http"
)

// BusinessLocation is a seller business location (địa điểm kinh doanh).
type BusinessLocation struct {
	XID      string `json:"xid"` // use as SellerStoreXID when creating invoices
	Code     string `json:"code"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	IsUsed   bool   `json:"is_used"`   // used locations cannot be deleted, only deactivated
	IsActive bool   `json:"is_active"` // inactive locations cannot be used on new invoices
}

// CreateBusinessLocationRequest is the body of POST /v1/business-locations/create.
// Code, Name and Address are required and unique within the company.
type CreateBusinessLocationRequest struct {
	Code    string `json:"code"`    // max 50 chars, immutable
	Name    string `json:"name"`    // max 400 chars
	Address string `json:"address"` // max 255 chars
}

// UpdateBusinessLocationRequest is the body of POST /v1/business-locations/update.
// Only non-empty / non-nil fields are changed.
type UpdateBusinessLocationRequest struct {
	XID      string `json:"xid"`
	Name     string `json:"name,omitempty"`
	Address  string `json:"address,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// ListBusinessLocations returns all business locations of the company.
func (c *Client) ListBusinessLocations(ctx context.Context) ([]BusinessLocation, error) {
	var out []BusinessLocation
	if err := c.do(ctx, http.MethodGet, "/v1/business-locations", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// CreateBusinessLocation adds a business location. New locations are active.
func (c *Client) CreateBusinessLocation(ctx context.Context, req *CreateBusinessLocationRequest) (*BusinessLocation, error) {
	var out BusinessLocation
	if err := c.do(ctx, http.MethodPost, "/v1/business-locations/create", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateBusinessLocation renames, re-addresses or (de)activates a business location.
func (c *Client) UpdateBusinessLocation(ctx context.Context, req *UpdateBusinessLocationRequest) (*BusinessLocation, error) {
	var out BusinessLocation
	if err := c.do(ctx, http.MethodPost, "/v1/business-locations/update", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteBusinessLocation deletes a location that was never used on an invoice.
// Used locations fail with CodeSellerStoreInUse; deactivate them instead.
func (c *Client) DeleteBusinessLocation(ctx context.Context, xid string) error {
	body := map[string]string{"xid": xid}
	return c.do(ctx, http.MethodPost, "/v1/business-locations/delete", nil, body, nil)
}
