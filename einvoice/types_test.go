package einvoice

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNumberUnmarshal(t *testing.T) {
	tests := []struct {
		in      string
		want    Number
		wantErr bool
	}{
		{`12.5`, 12.5, false},
		{`"2.00"`, 2, false},
		{`"-3"`, -3, false},
		{`null`, 0, false},
		{`""`, 0, false},
		{`"abc"`, 0, true},
		{`true`, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var n Number
			err := json.Unmarshal([]byte(tt.in), &n)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if n != tt.want {
				t.Errorf("n = %v, want %v", n, tt.want)
			}
		})
	}
}

func TestFormatDateTime(t *testing.T) {
	if got := formatDateTime(time.Time{}); got != "" {
		t.Errorf("zero = %q", got)
	}
	if got := formatDateTime(time.Date(2026, 12, 31, 17, 0, 0, 0, time.UTC)); got != "2027-01-01 00:00:00" {
		t.Errorf("got %q", got)
	}
}

func TestPtr(t *testing.T) {
	p := Ptr(TaxRate0)
	if p == nil || *p != TaxRate0 {
		t.Errorf("Ptr = %v", p)
	}
}
