package einvoice

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestListBusinessLocations(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/business-locations" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		io.WriteString(w, `{"success":true,"data":[
			{"xid":"x1","code":"HN","name":"Hà Nội","address":"1 Phố Huế","is_used":true,"is_active":true},
			{"xid":"x2","code":"SG","name":"Sài Gòn","address":"2 Lê Lợi","is_used":false,"is_active":false}]}`)
	})
	locs, err := client.ListBusinessLocations(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(locs) != 2 || locs[0].XID != "x1" || !locs[0].IsUsed || locs[1].IsActive {
		t.Errorf("locations = %+v", locs)
	}
}

func TestCreateBusinessLocation(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/business-locations/create" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		var got map[string]any
		json.NewDecoder(r.Body).Decode(&got)
		if got["code"] != "HN" || got["name"] != "Hà Nội" || got["address"] != "1 Phố Huế" {
			t.Errorf("body = %v", got)
		}
		io.WriteString(w, `{"success":true,"data":{"xid":"x1","code":"HN","name":"Hà Nội","address":"1 Phố Huế","is_active":true}}`)
	})
	loc, err := client.CreateBusinessLocation(context.Background(), &CreateBusinessLocationRequest{
		Code: "HN", Name: "Hà Nội", Address: "1 Phố Huế",
	})
	if err != nil {
		t.Fatal(err)
	}
	if loc.XID != "x1" || !loc.IsActive {
		t.Errorf("location = %+v", loc)
	}
}

func TestUpdateBusinessLocation(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/business-locations/update" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		var got map[string]any
		json.NewDecoder(r.Body).Decode(&got)
		if got["xid"] != "x1" || got["is_active"] != false {
			t.Errorf("body = %v", got)
		}
		for _, k := range []string{"name", "address"} {
			if _, ok := got[k]; ok {
				t.Errorf("%s should be omitted when empty", k)
			}
		}
		io.WriteString(w, `{"success":true,"data":{"xid":"x1","is_active":false}}`)
	})
	loc, err := client.UpdateBusinessLocation(context.Background(), &UpdateBusinessLocationRequest{
		XID: "x1", IsActive: Ptr(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	if loc.XID != "x1" || loc.IsActive {
		t.Errorf("location = %+v", loc)
	}
}

func TestDeleteBusinessLocation(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/business-locations/delete" {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"xid":"x1"}` {
			t.Errorf("body = %s", body)
		}
		io.WriteString(w, `{"success":true,"data":null}`)
	})
	if err := client.DeleteBusinessLocation(context.Background(), "x1"); err != nil {
		t.Fatal(err)
	}
}

func TestLocationEndpointsReturnErrors(t *testing.T) {
	client, _ := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"success":false,"error":{"code":"SELLER_STORE_IN_USE","message":"in use"}}`)
	})
	ctx := context.Background()

	if res, err := client.ListBusinessLocations(ctx); !HasCode(err, CodeSellerStoreInUse) || res != nil {
		t.Errorf("ListBusinessLocations = %v, %v", res, err)
	}
	if res, err := client.CreateBusinessLocation(ctx, &CreateBusinessLocationRequest{}); !HasCode(err, CodeSellerStoreInUse) || res != nil {
		t.Errorf("CreateBusinessLocation = %v, %v", res, err)
	}
	if res, err := client.UpdateBusinessLocation(ctx, &UpdateBusinessLocationRequest{}); !HasCode(err, CodeSellerStoreInUse) || res != nil {
		t.Errorf("UpdateBusinessLocation = %v, %v", res, err)
	}
	if err := client.DeleteBusinessLocation(ctx, "x1"); !HasCode(err, CodeSellerStoreInUse) {
		t.Errorf("DeleteBusinessLocation = %v", err)
	}
}
