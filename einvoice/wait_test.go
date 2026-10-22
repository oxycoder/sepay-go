package einvoice

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestWaitForCreate(t *testing.T) {
	var calls atomic.Int32
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/invoices/create/check/trk" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if calls.Add(1) < 3 {
			writeJSON(w, map[string]any{"success": true, "data": map[string]any{"status": "Pending", "message": "..."}})
			return
		}
		writeJSON(w, map[string]any{"success": true, "data": map[string]any{
			"reference_code": "ref", "status": "Success", "message": "ok",
			"invoice": map[string]any{"reference_code": "ref", "status": "issued"},
		}})
	})

	res, err := client.WaitForCreate(context.Background(), "trk")
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != StatusSuccess || res.Invoice == nil || res.Invoice.Status != InvoiceIssued {
		t.Errorf("res = %+v", res)
	}
	if n := calls.Load(); n != 3 {
		t.Errorf("polled %d times, want 3", n)
	}
}

func TestWaitForIssueStopsOnContext(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"success": true, "data": map[string]any{
			"status": "Pending", "retry": map[string]any{"retries_count": 1, "max_retries": 7, "next_retry_at": nil},
		}})
	})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	res, err := client.WaitForIssue(ctx, "trk")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v", err)
	}
	if res == nil || res.Status != StatusPending {
		t.Errorf("res = %+v", res)
	}
}

func TestPollDelay(t *testing.T) {
	c := New("", "", WithPollInterval(2*time.Second))
	next := time.Now().Add(90 * time.Second).In(vnZone).Format(dateTimeLayout)
	past := time.Now().Add(-time.Hour).In(vnZone).Format(dateTimeLayout)

	tests := []struct {
		name     string
		retry    *RetryInfo
		min, max time.Duration
	}{
		{"no retry block", nil, 2 * time.Second, 2 * time.Second},
		{"retry scheduled", &RetryInfo{NextRetryAt: &next}, 80 * time.Second, 90 * time.Second},
		{"retry in past", &RetryInfo{NextRetryAt: &past}, 2 * time.Second, 2 * time.Second},
		{"settling", &RetryInfo{}, time.Minute, time.Minute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := c.pollDelay(&ProcessResult{Status: StatusPending, Retry: tt.retry})
			if d < tt.min || d > tt.max {
				t.Errorf("delay = %v, want [%v, %v]", d, tt.min, tt.max)
			}
		})
	}
}

func TestUnconfirmed(t *testing.T) {
	r := &ProcessResult{Status: StatusFailed, Message: "Hệ thống chưa xác nhận được kết quả phát hành với nhà cung cấp hóa đơn"}
	if !r.Unconfirmed() {
		t.Error("Unconfirmed() = false")
	}
}

func TestWaitReturnsCheckError(t *testing.T) {
	client, _ := newTestServer(t, notFoundHandler)
	res, err := client.WaitForIssue(context.Background(), "trk")
	if !IsNotFound(err) || res != nil {
		t.Fatalf("res = %v, err = %v", res, err)
	}
}

func TestWaitReturnsFailedAsResult(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/invoices/issue/check/trk" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeJSON(w, map[string]any{"success": true, "data": map[string]any{"status": "Failed", "message": "x"}})
	})
	res, err := client.WaitForIssue(context.Background(), "trk")
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != StatusFailed || res.Unconfirmed() {
		t.Errorf("res = %+v", res)
	}
}

func TestNextRetryTime(t *testing.T) {
	bad := "not a time"
	good := "2026-01-02 03:04:05"
	if _, ok := (*RetryInfo)(nil).NextRetryTime(); ok {
		t.Error("nil RetryInfo: ok = true")
	}
	if _, ok := (&RetryInfo{NextRetryAt: &bad}).NextRetryTime(); ok {
		t.Error("malformed: ok = true")
	}
	got, ok := (&RetryInfo{NextRetryAt: &good}).NextRetryTime()
	if want := time.Date(2026, 1, 1, 20, 4, 5, 0, time.UTC); !ok || !got.Equal(want) {
		t.Errorf("NextRetryTime = %v, %v; want %v", got, ok, want)
	}
}
