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
