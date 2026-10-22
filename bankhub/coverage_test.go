package bankhub

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var ctx = context.Background()

func TestOptions(t *testing.T) {
	hc := &http.Client{}
	c := New("id", "secret", WithSandbox(), WithHTTPClient(hc), WithUserAgent("my-app/1.0"))
	if c.baseURL != SandboxURL || c.httpClient != hc || c.userAgent != "my-app/1.0" {
		t.Errorf("client = %+v", c)
	}
	c = New("id", "secret", WithBaseURL("https://x.example/"))
	if c.baseURL != "https://x.example" || c.userAgent != defaultUserAgent {
		t.Errorf("client = %+v", c)
	}
	if New("id", "secret").baseURL != ProductionURL {
		t.Error("default base URL")
	}
}

func TestUserAgentSent(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if ua := r.Header.Get("User-Agent"); ua != "my-app/1.0" {
			t.Errorf("User-Agent = %q", ua)
		}
		io.WriteString(w, `{"data":{}}`)
	})
	WithUserAgent("my-app/1.0")(client)
	if _, err := client.GetMerchantConfig(ctx); err != nil {
		t.Fatal(err)
	}
}

// expectRequest returns a handler checking method, path and (if non-empty)
// the query and JSON body, then writing resp.
func expectRequest(t *testing.T, method, path, query, reqBody, resp string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method || r.URL.Path != path {
			t.Errorf("got %s %s, want %s %s", r.Method, r.URL.Path, method, path)
		}
		if got := r.URL.Query().Encode(); got != query {
			t.Errorf("query = %s, want %s", got, query)
		}
		b, _ := io.ReadAll(r.Body)
		if string(b) != reqBody {
			t.Errorf("body = %s, want %s", b, reqBody)
		}
		if reqBody != "" && r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		io.WriteString(w, resp)
	}
}

func TestListBankAccounts(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodGet, "/v1/bank-account",
		"bank_id=2&company_xid=c-1&page=3&per_page=5&q=0123", "",
		`{"data":[{"xid":"b-1","account_number":"0123","account_type":"individual","accumulated":"1000.50","active":1,"bank_api_connected":"0"}],
		"meta":{"current_page":3,"per_page":5,"total":11,"page_count":3}}`))
	list, err := client.ListBankAccounts(ctx, &ListBankAccountsParams{
		PageParams: PageParams{Page: 3, PerPage: 5}, CompanyXID: "c-1", BankID: "2", Query: "0123",
	})
	if err != nil {
		t.Fatal(err)
	}
	a := list.Data[0]
	if a.XID != "b-1" || a.AccountType != AccountIndividual || a.Accumulated != 1000.5 || !a.Active || a.BankAPIConnected {
		t.Errorf("account = %+v", a)
	}
	if list.Meta.HasMore {
		t.Errorf("meta = %+v", list.Meta)
	}
}

func TestListBankAccountsNilParams(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodGet, "/v1/bank-account", "", "", `{"data":[],"meta":{}}`))
	list, err := client.ListBankAccounts(ctx, nil)
	if err != nil || len(list.Data) != 0 {
		t.Fatalf("list = %+v, err = %v", list, err)
	}
}

func TestGetBankAccount(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodGet, "/v1/bank-account/b 1", "", "",
		`{"data":{"xid":"b 1","phone_number":"0900","va_list":[{"xid":"va-1","sub_account":"VA1","active":"1"}]}}`))
	a, err := client.GetBankAccount(ctx, "b 1")
	if err != nil {
		t.Fatal(err)
	}
	if a.PhoneNumber != "0900" || len(a.VirtualAccounts) != 1 || !a.VirtualAccounts[0].Active {
		t.Errorf("account = %+v", a)
	}
}

func TestUpdateCompany(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodPost, "/v1/company/edit/c-1", "",
		`{"full_name":"ABC Corp","payment_code":"on","payment_code_suffix_from":0,"transaction_amount":"Unlimited"}`,
		`{"code":200,"message":"ok"}`))
	err := client.UpdateCompany(ctx, "c-1", &UpdateCompanyRequest{
		FullName: "ABC Corp", PaymentCode: PaymentCodeOn, PaymentCodeSuffixFrom: Ptr(0), TransactionAmount: TransactionAmountUnlimited,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCompanyCounter(t *testing.T) {
	resp := `{"data":{"dates":[{"date":"2024-01-01","transaction":3,"transaction_in":2,"transaction_out":1}],"total":{"transaction":3,"transaction_in":2,"transaction_out":1}}}`
	client, _ := newTestServer(t, expectRequest(t, http.MethodGet, "/v1/company/counter/c-1", "date=2024-01-01", "", resp))
	cnt, err := client.GetCompanyCounter(ctx, "c-1", time.Date(2023, 12, 31, 20, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if cnt.Total.Transaction != 3 || len(cnt.Dates) != 1 || cnt.Dates[0].TransactionIn != 2 {
		t.Errorf("counter = %+v", cnt)
	}

	client, _ = newTestServer(t, expectRequest(t, http.MethodGet, "/v1/company/counter/c-1", "", "", resp))
	if _, err := client.GetCompanyCounter(ctx, "c-1", time.Time{}); err != nil {
		t.Fatal(err)
	}
}

func TestGetMerchantCounter(t *testing.T) {
	resp := `{"data":{"dates":[{"date":"2024-01-01","company_id":42,"transaction":1}],"total":{"transaction":1}}}`
	client, _ := newTestServer(t, expectRequest(t, http.MethodGet, "/v1/merchant-config/counter",
		"end_date=2024-01-31&start_date=2024-01-01", "", resp))
	cnt, err := client.GetMerchantCounter(ctx, &MerchantCounterParams{
		StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, vnZone),
		EndDate:   time.Date(2024, 1, 31, 0, 0, 0, 0, vnZone),
	})
	if err != nil {
		t.Fatal(err)
	}
	if cnt.Dates[0].CompanyID != "42" || cnt.Total.Transaction != 1 {
		t.Errorf("counter = %+v", cnt)
	}

	client, _ = newTestServer(t, expectRequest(t, http.MethodGet, "/v1/merchant-config/counter", "", "", resp))
	if _, err := client.GetMerchantCounter(ctx, nil); err != nil {
		t.Fatal(err)
	}
}

func TestGetLinkToken(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodPost, "/v1/link-token/get", "", `{"link_token":"tok"}`,
		`{"data":{"xid":"lt-1","purpose":"LINK_BANK_ACCOUNT","status":"Active","is_mobile_app":0,"session_count":1,
		"company":{"xid":"c-1","full_name":"ABC","status":"Active"},
		"sessions":[{"xid":"s-1","state":"COMPLETED","context":{"a":1}}],
		"events":[{"xid":"e-1","name":"LINK_TOKEN_CREATED","metadata":{}}]}}`))
	lt, err := client.GetLinkToken(ctx, "tok")
	if err != nil {
		t.Fatal(err)
	}
	if lt.Status != LinkTokenActive || lt.Company.XID != "c-1" || lt.Sessions[0].State != SessionCompleted ||
		lt.Events[0].Name != EventLinkTokenCreated || lt.IsMobileApp {
		t.Errorf("link token = %+v", lt)
	}
}

func TestRevokeLinkToken(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodPost, "/v1/link-token/revoke", "", `{"link_token_xid":"lt-1"}`, `{"message":"ok"}`))
	if err := client.RevokeLinkToken(ctx, "lt-1"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateLinkTokenMinimal(t *testing.T) {
	req := &CreateLinkTokenRequest{CompanyXID: "c-1", Purpose: PurposeLink, Language: LanguageEN, CompletionRedirectURI: "https://x/done"}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"company_xid":"c-1","purpose":"LINK_BANK_ACCOUNT","completion_redirect_uri":"https://x/done","language":"en"}`
	if string(b) != want {
		t.Errorf("body = %s", b)
	}
}

func TestGetTransaction(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodGet, "/v1/transaction/t-1", "", "",
		`{"data":{"transaction_id":"t-1","amount":"250000","transfer_type":"debit","webhook_success":true}}`))
	tx, err := client.GetTransaction(ctx, "t-1")
	if err != nil {
		t.Fatal(err)
	}
	if tx.Amount != 250000 || tx.TransferType != TransferDebit || !tx.WebhookSuccess {
		t.Errorf("tx = %+v", tx)
	}
}

func TestListTransactionsAllParams(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodGet, "/v1/transaction",
		"bank_account_xid=b-1&brand_name=MBBank&company_xid=c-1&end_transaction_date=2024-01-02&page=1&per_page=2&q=DH&transaction_date=2024-01-03&webhook_success=1",
		"", `{"data":[],"meta":{"current_page":1,"page_count":1,"has_more":true}}`))
	list, err := client.ListTransactions(ctx, &ListTransactionsParams{
		PageParams:     PageParams{Page: 1, PerPage: 2},
		CompanyXID:     "c-1",
		BankAccountXID: "b-1",
		BrandName:      "MBBank",
		Date:           time.Date(2024, 1, 3, 0, 0, 0, 0, vnZone),
		EndDate:        time.Date(2024, 1, 2, 0, 0, 0, 0, vnZone),
		WebhookSuccess: Ptr(true),
		Query:          "DH",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !list.Meta.HasMore {
		t.Errorf("has_more from response ignored: %+v", list.Meta)
	}

	client, _ = newTestServer(t, expectRequest(t, http.MethodGet, "/v1/transaction", "", "", `{"data":[]}`))
	if _, err := client.ListTransactions(ctx, nil); err != nil {
		t.Fatal(err)
	}
}

func TestSimulateTransaction(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodPost, "/v1/transaction/create", "",
		`{"bank_account_xid":"b-1","transfer_type":"credit","amount":10000,"transaction_content":"DH1"}`,
		`{"data":{"transaction_id":"t-9","amount":10000,"transfer_type":"credit","accumulated":null}}`))
	tx, err := client.SimulateTransaction(ctx, &SimulateTransactionRequest{
		BankAccountXID: "b-1", TransferType: TransferCredit, Amount: 10000, TransactionContent: "DH1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if tx.TransactionID != "t-9" || tx.Amount != 10000 || tx.Accumulated != 0 {
		t.Errorf("tx = %+v", tx)
	}
}

func TestGetWebhook(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodGet, "/v1/webhook", "", "",
		`{"status":200,"data":{"webhook_url":"https://m/wh","auth_type":"SECRET_KEY","secret_key":"***5678","active":"1","allow_events":["*"]}}`))
	cfg, err := client.GetWebhook(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthType != AuthSecretKey || !cfg.Active || len(cfg.AllowEvents) != 1 || cfg.AllowEvents[0] != EventAll {
		t.Errorf("cfg = %+v", cfg)
	}
}

func TestUpsertWebhookFullBody(t *testing.T) {
	b, err := json.Marshal(&UpsertWebhookRequest{
		WebhookURL: "https://m/wh", AuthType: AuthSecretKey, SecretKey: "s", Active: Ptr(true),
		AllowEvents: []EventType{EventBankAccountLinked},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `{"webhook_url":"https://m/wh","auth_type":"SECRET_KEY","secret_key":"s","active":1,"allow_events":["BANK_ACCOUNT_LINKED"]}`
	if string(b) != want {
		t.Errorf("body = %s", b)
	}
}

func TestUpdateMerchantConfigKeywords(t *testing.T) {
	client, _ := newTestServer(t, expectRequest(t, http.MethodPost, "/v1/merchant-config", "",
		`{"transaction_whitelist_remark_keywords":["DH","SEVN"]}`,
		`{"data":{"transaction_whitelist_remark_keywords":["DH","SEVN"],"logo_url":"https://x/l.png"}}`))
	cfg, err := client.UpdateMerchantConfig(ctx, &UpdateMerchantConfigRequest{TransactionWhitelistRemarkKeywords: []string{"DH", "SEVN"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.TransactionWhitelistRemarkKeywords) != 2 || cfg.LogoURL != "https://x/l.png" {
		t.Errorf("cfg = %+v", cfg)
	}
}

func TestUploadMerchantLogoJPEG(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, h, err := r.FormFile("logo")
		if err != nil {
			t.Fatal(err)
		}
		if h.Filename != "Logo.JPEG" || h.Header.Get("Content-Type") != "image/jpeg" {
			t.Errorf("file = %s %v", h.Filename, h.Header)
		}
		io.WriteString(w, `{"data":{"logo_url":"https://s3/l.jpg"}}`)
	})
	if u, err := client.UploadMerchantLogo(ctx, "dir/Logo.JPEG", strings.NewReader("JPG")); err != nil || u != "https://s3/l.jpg" {
		t.Errorf("url = %q, err = %v", u, err)
	}
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, errors.New("disk error") }

func TestUploadMerchantLogoErrors(t *testing.T) {
	if _, err := New("id", "secret").UploadMerchantLogo(ctx, "logo.jpg", errReader{}); err == nil || !strings.Contains(err.Error(), "read logo") {
		t.Errorf("read error: %v", err)
	}

	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		io.WriteString(w, `{"message":"File too large"}`)
	})
	if _, err := client.UploadMerchantLogo(ctx, "logo.png", strings.NewReader("x")); !hasStatus(err, http.StatusRequestEntityTooLarge) {
		t.Errorf("API error: %v", err)
	}

	client, _ = newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `not json`)
	})
	if _, err := client.UploadMerchantLogo(ctx, "logo.png", strings.NewReader("x")); err == nil {
		t.Error("decode error: want error")
	}
}

// TestEndpointErrors checks that every endpoint returns the API error.
func TestEndpointErrors(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `{"code":404,"message":"Resource not found"}`)
	})
	calls := map[string]func() error{
		"ListBankAccounts":     func() error { _, err := client.ListBankAccounts(ctx, nil); return err },
		"GetBankAccount":       func() error { _, err := client.GetBankAccount(ctx, "x"); return err },
		"ListCompanies":        func() error { _, err := client.ListCompanies(ctx, nil); return err },
		"GetCompany":           func() error { _, err := client.GetCompany(ctx, "x"); return err },
		"CreateCompany":        func() error { _, err := client.CreateCompany(ctx, &CreateCompanyRequest{}); return err },
		"UpdateCompany":        func() error { return client.UpdateCompany(ctx, "x", &UpdateCompanyRequest{}) },
		"GetCompanyCounter":    func() error { _, err := client.GetCompanyCounter(ctx, "x", time.Time{}); return err },
		"CreateLinkToken":      func() error { _, err := client.CreateLinkToken(ctx, &CreateLinkTokenRequest{}); return err },
		"GetLinkToken":         func() error { _, err := client.GetLinkToken(ctx, "x"); return err },
		"RevokeLinkToken":      func() error { return client.RevokeLinkToken(ctx, "x") },
		"GetMerchantConfig":    func() error { _, err := client.GetMerchantConfig(ctx); return err },
		"UpdateMerchantConfig": func() error { _, err := client.UpdateMerchantConfig(ctx, &UpdateMerchantConfigRequest{}); return err },
		"GetMerchantCounter":   func() error { _, err := client.GetMerchantCounter(ctx, nil); return err },
		"ListTransactions":     func() error { _, err := client.ListTransactions(ctx, nil); return err },
		"GetTransaction":       func() error { _, err := client.GetTransaction(ctx, "x"); return err },
		"SimulateTransaction":  func() error { _, err := client.SimulateTransaction(ctx, &SimulateTransactionRequest{}); return err },
		"GetWebhook":           func() error { _, err := client.GetWebhook(ctx); return err },
		"UpsertWebhook":        func() error { _, err := client.UpsertWebhook(ctx, &UpsertWebhookRequest{}); return err },
	}
	for name, call := range calls {
		if err := call(); !IsNotFound(err) {
			t.Errorf("%s: err = %v", name, err)
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/company":
			io.WriteString(w, `{"data":"not a list"}`)
		case "/v1/webhook":
			io.WriteString(w, `[`)
		default:
			io.WriteString(w, `{"data":{"xid":5}}`)
		}
	})
	if _, err := client.ListCompanies(ctx, nil); err == nil || !strings.Contains(err.Error(), "decode response") {
		t.Errorf("list: %v", err)
	}
	if _, err := client.GetWebhook(ctx); err == nil || !strings.Contains(err.Error(), "decode response") {
		t.Errorf("envelope: %v", err)
	}
	if _, err := client.GetCompany(ctx, "c-1"); err == nil || !strings.Contains(err.Error(), "decode response data") {
		t.Errorf("data: %v", err)
	}
}

func TestDecodeDataNilOut(t *testing.T) {
	if found, err := decodeData([]byte("garbage"), nil); !found || err != nil {
		t.Errorf("found = %v, err = %v", found, err)
	}
}

func TestCallEncodeError(t *testing.T) {
	_, err := New("id", "secret").call(ctx, http.MethodPost, "/x", nil, make(chan int))
	if err == nil || !strings.Contains(err.Error(), "encode request") {
		t.Errorf("err = %v", err)
	}
}

func TestBuildRequestErrors(t *testing.T) {
	// The token request fails first on a bad base URL.
	c := New("id", "secret", WithBaseURL("http://bad\x7fhost"))
	if _, err := c.GetMerchantConfig(ctx); err == nil || !strings.Contains(err.Error(), "build request") {
		t.Errorf("token: %v", err)
	}

	// With a valid token, an invalid method fails building the API request.
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {})
	if _, err := client.call(ctx, "BAD METHOD", "/x", nil, nil); err == nil || !strings.Contains(err.Error(), "build request") {
		t.Errorf("request: %v", err)
	}
}

func TestTransportError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close()
	c := New("id", "secret", WithBaseURL(srv.URL))
	if _, err := c.GetMerchantConfig(ctx); err == nil || !strings.Contains(err.Error(), "POST /v1/token") {
		t.Errorf("err = %v", err)
	}
}

func TestReadResponseError(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		io.WriteString(w, `{"data"`)
	})
	if _, err := client.GetMerchantConfig(ctx); err == nil || !strings.Contains(err.Error(), "read response") {
		t.Errorf("err = %v", err)
	}
}

func tokenServer(t *testing.T, tokenResp string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, tokenResp)
	}))
	t.Cleanup(srv.Close)
	return New("id", "secret", WithBaseURL(srv.URL))
}

func TestTokenErrors(t *testing.T) {
	if _, err := tokenServer(t, `nope`).GetMerchantConfig(ctx); err == nil || !strings.Contains(err.Error(), "decode token") {
		t.Errorf("bad JSON: %v", err)
	}
	if _, err := tokenServer(t, `{"ttl":60}`).GetMerchantConfig(ctx); err == nil || !strings.Contains(err.Error(), "no access_token") {
		t.Errorf("no token: %v", err)
	}
}

func TestTokenShortTTL(t *testing.T) {
	c := tokenServer(t, `{"access_token":"t","ttl":60}`)
	tok, err := c.accessToken(ctx)
	if err != nil || tok != "t" {
		t.Fatalf("tok = %q, err = %v", tok, err)
	}
	// A TTL too short for the refresh margin is used as is.
	if d := time.Until(c.tokenExpiry); d < 50*time.Second || d > 60*time.Second {
		t.Errorf("expiry in %v", d)
	}
}

func TestInvalidateTokenKeepsNewer(t *testing.T) {
	c := New("id", "secret")
	c.token = "new"
	c.invalidateToken("old")
	if c.token != "new" {
		t.Error("newer token dropped")
	}
	c.invalidateToken("new")
	if c.token != "" {
		t.Error("rejected token kept")
	}
}

func TestAPIErrorString(t *testing.T) {
	err := &APIError{StatusCode: 500, Message: "boom"}
	if err.Error() != "bankhub: 500: boom" {
		t.Errorf("Error() = %q", err.Error())
	}
}

func TestParseAPIError(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		raw     string
		message string
		errors  map[string]string
	}{
		{"plain text", 500, " oops \n", "oops", nil},
		{"empty", 502, "", "Bad Gateway", nil},
		{"error string", 401, `{"error":"Unauthorized"}`, "Unauthorized", nil},
		{"error object", 400, `{"error":{"code":1}}`, `{"code":1}`, nil},
		{"messages", 400, `{"error":null,"messages":{"error":["a","b"]}}`, "a; b", nil},
		{"nothing", 403, `{}`, "Forbidden", nil},
		{"field errors", 400, `{"message":"m","errors":{"a":null,"b":1}}`, "m", map[string]string{"a": "", "b": "1"}},
	}
	for _, tt := range tests {
		e := parseAPIError(tt.status, []byte(tt.raw))
		if e.StatusCode != tt.status || e.Message != tt.message || len(e.Errors) != len(tt.errors) {
			t.Errorf("%s: got %+v", tt.name, e)
		}
		for k, v := range tt.errors {
			if e.Errors[k] != v {
				t.Errorf("%s: errors[%s] = %q, want %q", tt.name, k, e.Errors[k], v)
			}
		}
	}
}

func TestIsHelpersNonAPIError(t *testing.T) {
	if IsNotFound(errors.New("x")) || IsValidationError(nil) {
		t.Error("helpers matched non-API error")
	}
}

func TestFlag(t *testing.T) {
	for in, want := range map[string]Flag{`1`: true, `"1"`: true, `true`: true, `0`: false, `"0"`: false, `false`: false, `""`: false, `null`: false} {
		f := !want
		if err := json.Unmarshal([]byte(in), &f); err != nil || f != want {
			t.Errorf("%s: got %v, %v", in, f, err)
		}
	}
	var f Flag
	if err := json.Unmarshal([]byte(`2`), &f); err == nil {
		t.Error("2: want error")
	}
	if b, _ := json.Marshal(Flag(true)); string(b) != "1" {
		t.Errorf("Marshal(true) = %s", b)
	}
	if b, _ := json.Marshal(Flag(false)); string(b) != "0" {
		t.Errorf("Marshal(false) = %s", b)
	}
}

func TestNumber(t *testing.T) {
	tests := []struct {
		in   string
		want Number
		err  bool
	}{
		{`"1.5"`, 1.5, false},
		{`2`, 2, false},
		{`null`, 7, false},
		{`""`, 7, false},
		{`"x"`, 7, true},
	}
	for _, tt := range tests {
		n := Number(7)
		if err := json.Unmarshal([]byte(tt.in), &n); (err != nil) != tt.err || n != tt.want {
			t.Errorf("%s: got %v, %v", tt.in, n, err)
		}
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		in   string
		want Text
	}{
		{`"abc"`, "abc"},
		{`30`, "30"},
		{`null`, "keep"},
	}
	for _, tt := range tests {
		txt := Text("keep")
		if err := json.Unmarshal([]byte(tt.in), &txt); err != nil || txt != tt.want {
			t.Errorf("%s: got %q, %v", tt.in, txt, err)
		}
	}
	var txt Text
	if err := txt.UnmarshalJSON([]byte(`"bad\q"`)); err == nil {
		t.Error("invalid string: want error")
	}
}

func TestMetaUnmarshalError(t *testing.T) {
	var l List[Company]
	if err := json.Unmarshal([]byte(`{"data":[],"meta":{"page":"one"}}`), &l); err == nil {
		t.Error("want error")
	}
}

func TestParseEventErrors(t *testing.T) {
	if _, err := ParseEvent([]byte(`{`)); err == nil {
		t.Error("bad JSON: want error")
	}
	if _, err := ParseEvent([]byte(`{"event":"X","metadata":"oops"}`)); err == nil || !strings.Contains(err.Error(), "metadata") {
		t.Errorf("bad metadata: %v", err)
	}
	e, err := ParseEvent([]byte(`{"event":"LINK_TOKEN_EXPIRED","metadata":null}`))
	if err != nil || e.Event != EventLinkTokenExpired || e.Metadata != (EventMetadata{}) {
		t.Errorf("null metadata: %+v, %v", e, err)
	}
}

func TestParseIPNError(t *testing.T) {
	if _, err := ParseIPN([]byte(`{`)); err == nil {
		t.Error("want error")
	}
}

func TestEventHandlerErrors(t *testing.T) {
	h := EventHandler("", func(context.Context, *Event) error { return errors.New("db down") })
	if rec := post(h, "", "", linkedEvent); rec.Code != http.StatusInternalServerError || rec.Body.String() != `{"success":false}` {
		t.Errorf("fn error: %d %s", rec.Code, rec.Body)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/hook", nil))
	if rec.Code != http.StatusMethodNotAllowed || rec.Header().Get("Allow") != http.MethodPost {
		t.Errorf("GET: %d %v", rec.Code, rec.Header())
	}
}

func TestIPNHandlerErrors(t *testing.T) {
	called := false
	h := IPNHandler("", func(context.Context, *IPN) error { called = true; return nil })
	if rec := post(h, "", "", "not json"); rec.Code != http.StatusBadRequest {
		t.Errorf("bad body: %d", rec.Code)
	}
	if rec := post(h, "", "", strings.Repeat("a", maxBodySize+1)); rec.Code != http.StatusBadRequest {
		t.Errorf("too large: %d", rec.Code)
	}
	if called {
		t.Error("fn called on bad request")
	}
	if rec := post(h, "", "", ipnBody); rec.Code != http.StatusOK || !called {
		t.Errorf("no key: %d, called %v", rec.Code, called)
	}
}
