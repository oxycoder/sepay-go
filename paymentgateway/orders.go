package paymentgateway

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Order is a payment order. Times are "YYYY-MM-DD HH:mm:ss" in Vietnam time;
// parse them with ParseTime.
type Order struct {
	ID                   string        `json:"id"`          // SePay internal id
	CustomerID           string        `json:"customer_id"` // empty if none
	OrderID              string        `json:"order_id"`    // e.g. "SEPAY-68BA83CE637C1"
	OrderInvoiceNumber   string        `json:"order_invoice_number"`
	OrderStatus          OrderStatus   `json:"order_status"`
	OrderAmount          Number        `json:"order_amount"` // VND
	OrderCurrency        Currency      `json:"order_currency"`
	OrderDescription     string        `json:"order_description"`
	AuthenticationStatus string        `json:"authentication_status"`
	CreatedAt            string        `json:"created_at"`
	UpdatedAt            string        `json:"updated_at"`
	Transactions         []Transaction `json:"transactions"` // only set by GetOrder
}

// Transaction is a payment, refund or void on an order. Card fields are only
// set for PaymentCard.
type Transaction struct {
	ID                         string            `json:"id"`
	PaymentMethod              PaymentMethod     `json:"payment_method"`
	TransactionType            TransactionType   `json:"transaction_type"`
	TransactionAmount          Number            `json:"transaction_amount"`
	TransactionCurrency        Currency          `json:"transaction_currency"`
	TransactionStatus          TransactionStatus `json:"transaction_status"`
	AuthenticationStatus       string            `json:"authentication_status"`
	CardNumber                 string            `json:"card_number"` // masked
	CardHolderName             string            `json:"card_holder_name"`
	CardExpiry                 string            `json:"card_expiry"` // MMYY
	CardFundingMethod          CardFundingMethod `json:"card_funding_method"`
	CardBrand                  CardBrand         `json:"card_brand"`
	TransactionDate            string            `json:"transaction_date"`
	TransactionLastUpdatedDate string            `json:"transaction_last_updated_date"`
}

// SortOrder orders ListOrders results.
type SortOrder string

const (
	SortCreatedAsc  SortOrder = "created_at:asc"
	SortCreatedDesc SortOrder = "created_at:desc"
)

// ListOrdersParams filters ListOrders. Zero values are ignored. Dates are
// compared in Vietnam time.
type ListOrdersParams struct {
	PageParams
	Query       string // keyword search
	OrderStatus OrderStatus
	CustomerID  string
	CreatedAt   time.Time // a single day
	FromDate    time.Time // first day, inclusive
	ToDate      time.Time // last day, inclusive
	Sort        SortOrder
}

// ListOrders returns a page of orders. params may be nil.
func (c *Client) ListOrders(ctx context.Context, params *ListOrdersParams) (*List[Order], error) {
	q := url.Values{}
	if p := params; p != nil {
		if p.Page > 0 {
			q.Set("page", fmt.Sprint(p.Page))
		}
		if p.PerPage > 0 {
			q.Set("per_page", fmt.Sprint(p.PerPage))
		}
		setIf(q, "q", p.Query)
		setIf(q, "order_status", string(p.OrderStatus))
		setIf(q, "customer_id", p.CustomerID)
		setIf(q, "created_at", formatDate(p.CreatedAt))
		setIf(q, "from_created_at", formatDate(p.FromDate))
		setIf(q, "end_created_at", formatDate(p.ToDate))
		setIf(q, "sort", string(p.Sort))
	}
	var out List[Order]
	if err := c.call(ctx, http.MethodGet, "/v1/order", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOrder returns an order and its transactions by SePay order id (e.g.
// "SEPAY-68BA83CE637C1", Order.OrderID or IPNOrder.OrderID).
func (c *Client) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	var out struct {
		Data Order `json:"data"`
	}
	if err := c.call(ctx, http.MethodGet, "/v1/order/detail/"+url.PathEscape(orderID), nil, nil, &out); err != nil {
		return nil, err
	}
	return &out.Data, nil
}

type invoiceRequest struct {
	OrderInvoiceNumber string `json:"order_invoice_number"`
}

// CancelOrder cancels an unpaid order by its invoice number. It only applies
// to PaymentBankTransfer and PaymentNapasBankTransfer orders that are neither
// OrderCaptured nor OrderCancelled.
func (c *Client) CancelOrder(ctx context.Context, orderInvoiceNumber string) error {
	return c.call(ctx, http.MethodPost, "/v1/order/cancel", nil, invoiceRequest{orderInvoiceNumber}, nil)
}

// VoidTransaction voids the card payment of an order by its invoice number.
// It only applies to PaymentCard orders that are OrderCaptured and not yet
// settled: payments before 16:00 settle on the next working day (T+1), later
// ones after two working days (T+2).
func (c *Client) VoidTransaction(ctx context.Context, orderInvoiceNumber string) error {
	return c.call(ctx, http.MethodPost, "/v1/order/voidTransaction", nil, invoiceRequest{orderInvoiceNumber}, nil)
}

// setIf sets q[key] = v when v is not empty.
func setIf(q url.Values, key, v string) {
	if v != "" {
		q.Set(key, v)
	}
}
