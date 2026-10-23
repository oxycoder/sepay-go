package einvoice

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestCreateInvoiceRequestEncoding(t *testing.T) {
	var got map[string]any
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/invoices/create" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q", ct)
		}
		json.NewDecoder(r.Body).Decode(&got)
		writeJSON(w, map[string]any{"success": true, "data": map[string]any{
			"tracking_code": "trk-1", "tracking_url": "https://x/check/trk-1", "message": "ok",
		}})
	})

	issued := time.Date(2026, 1, 25, 18, 30, 0, 0, time.UTC) // 01:30 next day in Vietnam
	res, err := client.CreateInvoice(context.Background(), &CreateInvoiceRequest{
		TemplateCode:      TemplateVAT,
		InvoiceSeries:     "C26TSE",
		IssuedDate:        issued,
		Currency:          CurrencyVND,
		ProviderAccountID: "acc",
		Buyer:             Buyer{Name: "Công ty ABC"},
		Items: []Item{
			{LineNumber: 1, LineType: LineTypeGoods, ItemName: "A", Quantity: 1.5, UnitPrice: 100000, TaxRate: Ptr(TaxRate0)},
			{LineNumber: 2, LineType: LineTypeNote, ItemName: "note"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.TrackingCode != "trk-1" {
		t.Errorf("TrackingCode = %q", res.TrackingCode)
	}

	if got["issued_date"] != "2026-01-26 01:30:00" {
		t.Errorf("issued_date = %v", got["issued_date"])
	}
	if !res.IssuedDate.Equal(issued) || res.IssuedDate.Location() != vnZone {
		t.Errorf("res.IssuedDate = %v, want %v in Vietnam time", res.IssuedDate, issued)
	}
	if got["is_draft"] != false {
		t.Errorf("is_draft = %v, want explicit false", got["is_draft"])
	}
	if _, ok := got["total_amount"]; ok {
		t.Error("total_amount should be omitted when nil")
	}
	items := got["items"].([]any)
	first := items[0].(map[string]any)
	if first["tax_rate"] != float64(0) {
		t.Errorf("tax_rate = %v, want 0 to be sent", first["tax_rate"])
	}
	if first["quantity"] != 1.5 {
		t.Errorf("quantity = %v", first["quantity"])
	}
	note := items[1].(map[string]any)
	for _, k := range []string{"tax_rate", "quantity", "unit_price", "tax_amount"} {
		if _, ok := note[k]; ok {
			t.Errorf("note line should omit %s", k)
		}
	}
}

func TestCreateInvoiceDefaultsIssuedDateToNow(t *testing.T) {
	b, err := json.Marshal(&CreateInvoiceRequest{})
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	json.Unmarshal(b, &got)
	ts, err := time.ParseInLocation(dateTimeLayout, got["issued_date"].(string), vnZone)
	if err != nil {
		t.Fatal(err)
	}
	if d := time.Since(ts); d < 0 || d > time.Minute {
		t.Errorf("issued_date %v is not now", ts)
	}
}

func TestAdjustOmitsZeroIssuedDate(t *testing.T) {
	b, err := json.Marshal(&AdjustInvoiceRequest{AdjustmentType: AdjustInfo, Items: []Item{{ItemName: "x", Amount: Ptr(-550000.0)}}})
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	json.Unmarshal(b, &got)
	if _, ok := got["issued_date"]; ok {
		t.Errorf("issued_date should be omitted: %s", b)
	}
	if _, ok := got["buyer"]; ok {
		t.Errorf("buyer should be omitted when nil: %s", b)
	}
	if got["is_draft"] != false || got["adjustment_type"] != "info" {
		t.Errorf("unexpected body: %s", b)
	}
	if got["items"].([]any)[0].(map[string]any)["amount"] != -550000.0 {
		t.Errorf("amount not sent: %s", b)
	}
}

func TestGetInvoiceDecodesStringNumbers(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/invoices/ref-1" {
			t.Errorf("path = %s", r.URL.Path)
		}
		io.WriteString(w, `{"success":true,"data":{
			"reference_code":"ref-1","invoice_number":"589","status":"issued","document_type":2,
			"xml_url":null,"seller_store":null,"total_before_tax":200000,"tax_amount":20000,"total_amount":220000,
			"items":[{"line_number":1,"line_type":1,"item_name":"A","quantity":"2.00","unit_price":"100000.00",
				"total_amount":"200000.00","tax_rate":"10.00","tax_amount":"20000.00"}]}}`)
	})
	inv, err := client.GetInvoice(context.Background(), "ref-1")
	if err != nil {
		t.Fatal(err)
	}
	if inv.Status != InvoiceIssued || inv.DocumentType != DocNew || inv.SellerStore != nil || inv.TotalAmount != 220000 {
		t.Errorf("invoice = %+v", inv)
	}
	line := inv.Items[0]
	if line.Quantity != 2 || line.UnitPrice != 100000 || line.TaxRate != 10 {
		t.Errorf("line = %+v", line)
	}
}

func TestListInvoicesQuery(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("page") != "2" || q.Get("per_page") != "50" || q.Get("source") != "api" || q.Get("document_type") != "6,7" {
			t.Errorf("query = %s", r.URL.RawQuery)
		}
		io.WriteString(w, `{"data":{"paging":{"per_page":50,"total":51,"has_more":false,"current_page":2,"page_count":2},
			"items":[{"reference_code":"r1","status":"draft","invoice_number":"0"}]}}`)
	})
	list, err := client.ListInvoices(context.Background(), &ListInvoicesParams{
		Page: 2, PerPage: 50, Source: SourceAPI, DocumentTypes: []DocumentType{DocAdjustIncrease, DocAdjustDecrease},
	})
	if err != nil {
		t.Fatal(err)
	}
	if list.Paging.Total != 51 || len(list.Items) != 1 || list.Items[0].Status != InvoiceDraft {
		t.Errorf("list = %+v", list)
	}
}

func TestDownloadInvoiceDecodesBase64(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/invoices/trk/download" || r.URL.Query().Get("type") != "pdf" {
			t.Errorf("%s?%s", r.URL.Path, r.URL.RawQuery)
		}
		io.WriteString(w, `{"success":true,"data":{"file_type":"pdf","file_name":"HD.pdf","content":"JVBERi0xLjQ="}}`)
	})
	f, err := client.DownloadInvoice(context.Background(), "trk", FilePDF)
	if err != nil {
		t.Fatal(err)
	}
	if string(f.Content) != "%PDF-1.4" || f.FileName != "HD.pdf" {
		t.Errorf("file = %+v", f)
	}
}

func TestAdjustInvoice(t *testing.T) {
	var got map[string]any
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/invoices/adjust" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&got)
		writeJSON(w, map[string]any{"success": true, "data": map[string]any{"tracking_code": "adj-1"}})
	})
	res, err := client.AdjustInvoice(context.Background(), &AdjustInvoiceRequest{
		ProviderAccountID:            "acc",
		OriginalInvoiceReferenceCode: "ref-1",
		AdjustmentType:               AdjustInfo,
		IssuedDate:                   time.Date(2026, 3, 1, 3, 0, 0, 0, time.UTC),
		Buyer:                        &Buyer{Name: "B"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.TrackingCode != "adj-1" {
		t.Errorf("TrackingCode = %q", res.TrackingCode)
	}
	if !res.IssuedDate.Equal(time.Date(2026, 3, 1, 3, 0, 0, 0, time.UTC)) {
		t.Errorf("res.IssuedDate = %v", res.IssuedDate)
	}
	if got["issued_date"] != "2026-03-01 10:00:00" || got["original_invoice_reference_code"] != "ref-1" {
		t.Errorf("body = %v", got)
	}
	if buyer, ok := got["buyer"].(map[string]any); !ok || buyer["name"] != "B" {
		t.Errorf("buyer = %v", got["buyer"])
	}
}

func TestReplaceInvoice(t *testing.T) {
	var got map[string]any
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/invoices/replace" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&got)
		writeJSON(w, map[string]any{"success": true, "data": map[string]any{"tracking_code": "rep-1"}})
	})
	res, err := client.ReplaceInvoice(context.Background(), &ReplaceInvoiceRequest{
		ProviderAccountID:            "acc",
		OriginalInvoiceReferenceCode: "ref-1",
		IssuedDate:                   time.Date(2026, 3, 1, 3, 0, 0, 0, time.UTC),
		Buyer:                        Buyer{Name: "B"},
		Items:                        []Item{{ItemName: "x", Quantity: 1, UnitPrice: 10}},
		TotalAmount:                  Ptr(int64(11)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.TrackingCode != "rep-1" {
		t.Errorf("TrackingCode = %q", res.TrackingCode)
	}
	if got["issued_date"] != "2026-03-01 10:00:00" || got["total_amount"] != float64(11) || got["is_draft"] != false {
		t.Errorf("body = %v", got)
	}
}

func TestSubmitDefaultsIssuedDateToNow(t *testing.T) {
	var sent []string
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var got map[string]any
		json.NewDecoder(r.Body).Decode(&got)
		s, _ := got["issued_date"].(string)
		sent = append(sent, s)
		writeJSON(w, map[string]any{"success": true, "data": map[string]any{"tracking_code": "t"}})
	})
	ctx := context.Background()
	before := time.Now().Truncate(time.Second)
	submit := []func() (*TrackingResult, error){
		func() (*TrackingResult, error) {
			return client.CreateInvoice(ctx, &CreateInvoiceRequest{})
		},
		func() (*TrackingResult, error) {
			return client.AdjustInvoice(ctx, &AdjustInvoiceRequest{})
		},
		func() (*TrackingResult, error) {
			return client.ReplaceInvoice(ctx, &ReplaceInvoiceRequest{})
		},
	}
	for i, f := range submit {
		res, err := f()
		if err != nil {
			t.Fatal(err)
		}
		ts := res.IssuedDate
		if ts.Before(before) || ts.After(time.Now()) || ts.Nanosecond() != 0 {
			t.Errorf("%d: IssuedDate %v is not now to the second", i, ts)
		}
		if sent[i] != ts.Format(dateTimeLayout) || ts.Location() != vnZone {
			t.Errorf("%d: sent %q, result %v", i, sent[i], ts)
		}
	}
}

func TestReplaceOmitsZeroIssuedDate(t *testing.T) {
	b, err := json.Marshal(&ReplaceInvoiceRequest{})
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	json.Unmarshal(b, &got)
	if _, ok := got["issued_date"]; ok {
		t.Errorf("issued_date should be omitted: %s", b)
	}
	if _, ok := got["currency"]; ok {
		t.Errorf("currency should be omitted when empty: %s", b)
	}
}

func TestDeleteDraftInvoice(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/invoices/delete/ref-1" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "" {
			t.Error("Content-Type set without a body")
		}
		io.WriteString(w, `{"success":true}`)
	})
	if err := client.DeleteDraftInvoice(context.Background(), "ref-1"); err != nil {
		t.Fatal(err)
	}
}

func TestListInvoicesSellerStoreAndNilParams(t *testing.T) {
	var queries []string
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		io.WriteString(w, `{"success":true,"data":{"items":[]}}`)
	})
	ctx := context.Background()
	if _, err := client.ListInvoices(ctx, &ListInvoicesParams{SellerStoreXID: "x1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.ListInvoices(ctx, nil); err != nil {
		t.Fatal(err)
	}
	if len(queries) != 2 || queries[0] != "seller_store_xid=x1" || queries[1] != "" {
		t.Errorf("queries = %q", queries)
	}
}

func TestInvoiceEndpointsReturnErrors(t *testing.T) {
	client, _ := newTestServer(t, notFoundHandler)
	ctx := context.Background()

	check := func(name string, res any, err error) {
		t.Helper()
		if !IsNotFound(err) {
			t.Errorf("%s err = %v", name, err)
		}
		if res != nil {
			t.Errorf("%s res = %v, want nil", name, res)
		}
	}
	res1, err := client.CreateInvoice(ctx, &CreateInvoiceRequest{})
	check("CreateInvoice", nilIfNil(res1), err)
	res2, err := client.AdjustInvoice(ctx, &AdjustInvoiceRequest{})
	check("AdjustInvoice", nilIfNil(res2), err)
	res3, err := client.ReplaceInvoice(ctx, &ReplaceInvoiceRequest{})
	check("ReplaceInvoice", nilIfNil(res3), err)
	res4, err := client.IssueInvoice(ctx, "ref")
	check("IssueInvoice", nilIfNil(res4), err)
	res5, err := client.GetCreateStatus(ctx, "trk")
	check("GetCreateStatus", nilIfNil(res5), err)
	res6, err := client.GetIssueStatus(ctx, "trk")
	check("GetIssueStatus", nilIfNil(res6), err)
	res7, err := client.ListInvoices(ctx, nil)
	check("ListInvoices", nilIfNil(res7), err)
	res8, err := client.DownloadInvoice(ctx, "trk", FileXML)
	check("DownloadInvoice", nilIfNil(res8), err)
	res9, err := client.GetUsage(ctx)
	check("GetUsage", nilIfNil(res9), err)
	check("DeleteDraftInvoice", nil, client.DeleteDraftInvoice(ctx, "ref"))
}

// nilIfNil turns a typed nil pointer into an untyped nil so it compares equal to nil as any.
func nilIfNil[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}
