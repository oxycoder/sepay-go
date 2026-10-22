package einvoice

import (
	"context"
	"io"
	"net/http"
	"testing"
)

func TestListEmailLogs(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/email-logs" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		q := r.URL.Query()
		want := map[string]string{
			"page": "1", "per_page": "20", "invoice_number": "589", "status": "sent",
			"to_email": "a@b.vn", "reference_code": "ref-1",
		}
		for k, v := range want {
			if q.Get(k) != v {
				t.Errorf("%s = %q, want %q", k, q.Get(k), v)
			}
		}
		io.WriteString(w, `{"success":true,"data":{"paging":{"total":1},"items":[{"invoice_number":"589",
			"to_email":"a@b.vn","status":"sent","sent_at":"2026-01-01 10:00:00","created_at":"2026-01-01 09:59:00"}]}}`)
	})
	list, err := client.ListEmailLogs(context.Background(), &ListEmailLogsParams{
		Page: 1, PerPage: 20, InvoiceNumber: "589", Status: EmailSent, ToEmail: "a@b.vn", ReferenceCode: "ref-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if list.Paging.Total != 1 || len(list.Items) != 1 || list.Items[0].Status != EmailSent || list.Items[0].SentAt == "" {
		t.Errorf("list = %+v", list)
	}
}

func TestListEmailLogsOmitsEmptyFilters(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "status=failed" {
			t.Errorf("query = %s", r.URL.RawQuery)
		}
		io.WriteString(w, `{"success":true,"data":{"items":[]}}`)
	})
	if _, err := client.ListEmailLogs(context.Background(), &ListEmailLogsParams{Status: EmailFailed}); err != nil {
		t.Fatal(err)
	}
}

func TestListEmailLogsNilParamsAndError(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("query = %s, want empty", r.URL.RawQuery)
		}
		notFoundHandler(w, r)
	})
	if res, err := client.ListEmailLogs(context.Background(), nil); !IsNotFound(err) || res != nil {
		t.Errorf("ListEmailLogs = %v, %v", res, err)
	}
}
