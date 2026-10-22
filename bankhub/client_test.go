package bankhub

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// newTestServer returns a client wired to a fake API. The token endpoint is
// handled automatically and hands out "tok-1", "tok-2", ...
func newTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *atomic.Int32) {
	t.Helper()
	var tokens atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/token" {
			id, secret, ok := r.BasicAuth()
			if !ok || id != "id" || secret != "secret" {
				w.WriteHeader(http.StatusUnauthorized)
				io.WriteString(w, `{"code":401,"message":"Unauthenticated"}`)
				return
			}
			n := tokens.Add(1)
			w.WriteHeader(http.StatusCreated)
			writeJSON(w, map[string]any{"code": 201, "access_token": "tok-" + string(rune('0'+n)), "ttl": 60000})
			return
		}
		handler(w, r)
	}))
	t.Cleanup(srv.Close)
	return New("id", "secret", WithBaseURL(srv.URL)), &tokens
}

func writeJSON(w http.ResponseWriter, v any) {
	json.NewEncoder(w).Encode(v)
}

func TestTokenIsCachedAndSent(t *testing.T) {
	client, tokens := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer tok-1" {
			t.Errorf("Authorization = %q", got)
		}
		writeJSON(w, map[string]any{"data": map[string]any{"logo_url": "https://x/logo.png"}})
	})
	for range 3 {
		if _, err := client.GetMerchantConfig(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if n := tokens.Load(); n != 1 {
		t.Errorf("token fetched %d times, want 1", n)
	}
}

func TestTokenRefreshedOn401(t *testing.T) {
	client, tokens := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "Bearer tok-1" {
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, `{"code":401,"message":"Unauthenticated"}`)
			return
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"full_name":"ABC"`) {
			t.Errorf("body not resent on retry: %s", body)
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"code":201,"message":"ok","data":{"xid":"c-1","full_name":"ABC"}}`)
	})
	c, err := client.CreateCompany(context.Background(), &CreateCompanyRequest{FullName: "ABC", Status: CompanyActive})
	if err != nil {
		t.Fatal(err)
	}
	if c.XID != "c-1" {
		t.Errorf("XID = %q", c.XID)
	}
	if n := tokens.Load(); n != 2 {
		t.Errorf("token fetched %d times, want 2", n)
	}
}

func TestBadCredentials(t *testing.T) {
	client, _ := newTestServer(t, nil)
	client.clientSecret = "wrong"
	_, err := client.GetMerchantConfig(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 401 || apiErr.Message != "Unauthenticated" {
		t.Fatalf("err = %v", err)
	}
}

func TestValidationError(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"code":400,"message":"Bad request","errors":{"full_name":"Tên quá ngắn","status":["a","b"]}}`)
	})
	_, err := client.CreateCompany(context.Background(), &CreateCompanyRequest{})
	if !IsValidationError(err) {
		t.Fatalf("err = %v", err)
	}
	var apiErr *APIError
	errors.As(err, &apiErr)
	if apiErr.Errors["full_name"] != "Tên quá ngắn" || apiErr.Errors["status"] != "a; b" {
		t.Errorf("Errors = %v", apiErr.Errors)
	}
	if !strings.Contains(err.Error(), "full_name: Tên quá ngắn") {
		t.Errorf("Error() = %q", err.Error())
	}
}

func TestWebhookStyleError(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"status":400,"error":"","messages":{"error":["URL không hợp lệ"]}}`)
	})
	_, err := client.UpsertWebhook(context.Background(), &UpsertWebhookRequest{WebhookURL: "http://10.0.0.1"})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Message != "URL không hợp lệ" {
		t.Fatalf("err = %v", err)
	}
}

func TestNotFound(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"code":404,"message":"Resource not found"}`)
	})
	_, err := client.GetCompany(context.Background(), "missing")
	if !IsNotFound(err) {
		t.Fatalf("err = %v", err)
	}
}

func TestLinkTokenFlatResponse(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/link-token/create" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["company_xid"] != "c-1" || body["purpose"] != "UNLINK_BANK_ACCOUNT" ||
			body["bank_account_xid"] != "b-1" || body["is_mobile_app"] != float64(1) {
			t.Errorf("body = %v", body)
		}
		if _, ok := body["language"]; ok {
			t.Error("empty language sent")
		}
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"xid":"lt-1","hosted_link_url":"https://bankhub.sepay.vn/link/lt-1","link_token":"tok","expires_at":"2024-01-17 10:30:00"}`)
	})
	lt, err := client.CreateLinkToken(context.Background(), &CreateLinkTokenRequest{
		CompanyXID: "c-1", Purpose: PurposeUnlink, BankAccountXID: "b-1", IsMobileApp: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if lt.HostedLinkURL != "https://bankhub.sepay.vn/link/lt-1" || lt.LinkToken != "tok" {
		t.Errorf("lt = %+v", lt)
	}
}

func TestListCompanies(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Encode(); got != "page=2&per_page=10&sort%5Bcreated_at%5D=desc&status=Active" {
			t.Errorf("query = %s", got)
		}
		io.WriteString(w, `{"data":[{"xid":"c-1","full_name":"ABC","status":"Active","active":"1"},{"xid":"c-2","active":0}],
			"meta":{"per_page":10,"total":12,"has_more":false,"current_page":2,"page_count":2}}`)
	})
	list, err := client.ListCompanies(context.Background(), &ListCompaniesParams{
		PageParams: PageParams{Page: 2, PerPage: 10}, Status: CompanyActive, SortCreatedAt: "desc",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Data) != 2 || !list.Data[0].Active || list.Data[1].Active {
		t.Errorf("data = %+v", list.Data)
	}
	if list.Meta != (Meta{CurrentPage: 2, PerPage: 10, Total: 12, PageCount: 2}) {
		t.Errorf("meta = %+v", list.Meta)
	}
}

func TestListTransactions(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("start_transaction_date") != "2024-01-01" || q.Get("transfer_type") != "credit" || q.Get("webhook_success") != "0" {
			t.Errorf("query = %s", q.Encode())
		}
		io.WriteString(w, `{"data":[{"transaction_id":"t-1","amount":500000,"transfer_type":"credit","webhook_success":"1"}],
			"meta":{"page":1,"per_page":20,"total":30,"total_pages":2}}`)
	})
	// 2023-12-31 20:00 UTC is 2024-01-01 in Vietnam.
	start := time.Date(2023, 12, 31, 20, 0, 0, 0, time.UTC)
	list, err := client.ListTransactions(context.Background(), &ListTransactionsParams{
		StartDate: start, TransferType: TransferCredit, WebhookSuccess: Ptr(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	if tx := list.Data[0]; tx.Amount != 500000 || !tx.WebhookSuccess {
		t.Errorf("tx = %+v", tx)
	}
	if m := list.Meta; m.CurrentPage != 1 || m.PageCount != 2 || !m.HasMore {
		t.Errorf("meta = %+v", m)
	}
}

func TestGetCompanyConfiguration(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"data":{"xid":"c-1","status":"Active","active":"1","configuration":{
			"payment_code":"on","payment_code_suffix_from":3,"transaction_amount":"Unlimited","data_storage_time":30}}}`)
	})
	c, err := client.GetCompany(context.Background(), "c-1")
	if err != nil {
		t.Fatal(err)
	}
	if cfg := c.Configuration; cfg == nil || cfg.TransactionAmount != TransactionAmountUnlimited || cfg.DataStorageTime != "30" || cfg.PaymentCodeSuffixFrom != 3 {
		t.Errorf("configuration = %+v", c.Configuration)
	}
}

func TestGetWebhookNotConfigured(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"status":200,"error":null,"messages":{"success":[],"error":[]},"data":null}`)
	})
	cfg, err := client.GetWebhook(context.Background())
	if err != nil || cfg != nil {
		t.Fatalf("cfg = %v, err = %v", cfg, err)
	}
}

func TestUpsertWebhookBody(t *testing.T) {
	var got string
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = string(b)
		io.WriteString(w, `{"status":200,"data":{"webhook_url":"https://m/wh","active":1,"allow_events":[]}}`)
	})
	ctx := context.Background()

	cfg, err := client.UpsertWebhook(ctx, &UpsertWebhookRequest{WebhookURL: "https://m/wh", Active: Ptr(false), AllowEvents: []EventType{}})
	if err != nil {
		t.Fatal(err)
	}
	if got != `{"webhook_url":"https://m/wh","active":0,"allow_events":[]}` {
		t.Errorf("body = %s", got)
	}
	if !cfg.Active {
		t.Errorf("cfg = %+v", cfg)
	}

	client.UpsertWebhook(ctx, &UpsertWebhookRequest{SecretKey: "s"})
	if got != `{"secret_key":"s"}` {
		t.Errorf("partial body = %s", got)
	}
}

func TestUpdateMerchantConfigClearsKeywords(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		if string(b) != `{"transaction_whitelist_remark_keywords":[]}` {
			t.Errorf("body = %s", b)
		}
		io.WriteString(w, `{"data":{"transaction_whitelist_remark_keywords":[]}}`)
	})
	if _, err := client.UpdateMerchantConfig(context.Background(), &UpdateMerchantConfigRequest{}); err != nil {
		t.Fatal(err)
	}
}

func TestUploadMerchantLogo(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		f, h, err := r.FormFile("logo")
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(f)
		if h.Filename != "logo.png" || h.Header.Get("Content-Type") != "image/png" || string(data) != "PNGDATA" {
			t.Errorf("file = %s %v %q", h.Filename, h.Header, data)
		}
		io.WriteString(w, `{"data":{"logo_url":"https://s3/1-2.png"}}`)
	})
	u, err := client.UploadMerchantLogo(context.Background(), "/tmp/logo.png", strings.NewReader("PNGDATA"))
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://s3/1-2.png" {
		t.Errorf("url = %q", u)
	}
	if _, err := client.UploadMerchantLogo(context.Background(), "logo.gif", strings.NewReader("")); err == nil {
		t.Error("gif accepted")
	}
}
