package einvoice

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestListProviderAccounts(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/provider-accounts" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("page") != "3" || q.Get("per_page") != "100" {
			t.Errorf("query = %s", r.URL.RawQuery)
		}
		io.WriteString(w, `{"success":true,"data":{"paging":{"per_page":100,"total":1,"current_page":3},
			"items":[{"id":"acc-1","provider":"matbao","active":true,"tax_authority_approved_date":"2026-01-01"}]}}`)
	})
	list, err := client.ListProviderAccounts(context.Background(), &PageParams{Page: 3, PerPage: 100})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].ID != "acc-1" || !list.Items[0].Active || list.Items[0].TaxAuthorityApprovedDate != "2026-01-01" {
		t.Errorf("list = %+v", list)
	}
}

func TestListProviderAccountsNilParams(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("query = %s, want empty", r.URL.RawQuery)
		}
		io.WriteString(w, `{"success":true,"data":{"items":[]}}`)
	})
	list, err := client.ListProviderAccounts(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 0 {
		t.Errorf("list = %+v", list)
	}
}

func TestGetProviderAccount(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/v1/provider-accounts/a%2Fb" {
			t.Errorf("path = %s", r.URL.EscapedPath())
		}
		io.WriteString(w, `{"success":true,"data":{"id":"a/b","provider":"bkav","active":true,
			"templates":[{"template_code":"1","invoice_series":"C26TSE","invoice_label":"VAT"}]}}`)
	})
	acc, err := client.GetProviderAccount(context.Background(), "a/b")
	if err != nil {
		t.Fatal(err)
	}
	if acc.ID != "a/b" || len(acc.Templates) != 1 || acc.Templates[0].InvoiceSeries != "C26TSE" {
		t.Errorf("account = %+v", acc)
	}
}

func TestAccountEndpointsReturnErrors(t *testing.T) {
	client, _ := newTestServer(t, notFoundHandler)
	ctx := context.Background()

	if res, err := client.ListProviderAccounts(ctx, nil); !IsNotFound(err) || res != nil {
		t.Errorf("ListProviderAccounts = %v, %v", res, err)
	}
	if res, err := client.GetProviderAccount(ctx, "x"); !IsNotFound(err) || res != nil {
		t.Errorf("GetProviderAccount = %v, %v", res, err)
	}
}

// notFoundHandler answers every API call with a 404 APIError.
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, `{"success":false,"error":{"code":"NOT_FOUND","message":"missing"}}`)
}
