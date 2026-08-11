package xenditv3

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestErrorResponseCapturesValidationDetails(t *testing.T) {
	raw := `{
		"error_code": "API_VALIDATION_ERROR",
		"message": "Inputs are failing validation. The errors field contains details about which fields are violating validation",
		"errors": [
			{"path": "recipient.address.postal_code", "message": "must be a valid postal code"}
		]
	}`

	var resp ErrorResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(resp.Errors) != 1 {
		t.Fatalf("expected 1 validation detail, got %d", len(resp.Errors))
	}
	if !strings.Contains(resp.Error(), "recipient.address.postal_code") {
		t.Errorf("Error() should surface the failing field, got %q", resp.Error())
	}
}

func TestErrorResponseWithoutDetailsFallsBackToMessage(t *testing.T) {
	resp := ErrorResponse{ErrorCode: "SOME_ERROR", ErrorMessage: "something broke"}
	if got := resp.Error(); got != "[SOME_ERROR] something broke" {
		t.Errorf("got %q", got)
	}
}
