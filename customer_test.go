package xenditgo

import "testing"

func TestParseCustomerListResponseArray(t *testing.T) {
	body := []byte(`[{"id":"0b7bd4ef-b054-43a6-ae4d-0e6b7d4d1d53","reference_id":"1640672049749823488","email":"a@b.com"}]`)
	customers, err := parseCustomerListResponse(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(customers) != 1 {
		t.Fatalf("expected 1 customer, got %d", len(customers))
	}
	if customers[0].ID != "0b7bd4ef-b054-43a6-ae4d-0e6b7d4d1d53" {
		t.Fatalf("unexpected id: %s", customers[0].ID)
	}
}

func TestNormalizeXenditCustomerID(t *testing.T) {
	raw := "0b7bd4ef-b054-43a6-ae4d-0e6b7d4d1d53"
	normalized := NormalizeXenditCustomerID(raw)
	if len(normalized) < 41 {
		t.Fatalf("customer_id too short: %q (%d)", normalized, len(normalized))
	}
	if normalized != "cust-"+raw {
		t.Fatalf("unexpected normalized id: %s", normalized)
	}
	if NormalizeXenditCustomerID(normalized) != normalized {
		t.Fatalf("should be idempotent")
	}
}

func TestParseCustomerListResponseWrapped(t *testing.T) {
	body := []byte(`{"data":[{"id":"cust-1","reference_id":"user-1","email":"a@b.com"}],"has_more":false}`)
	customers, err := parseCustomerListResponse(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(customers) != 1 || customers[0].ID != "cust-1" {
		t.Fatalf("unexpected customers: %+v", customers)
	}
}
