package xenditv3

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
)

func TestCreatePayoutHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/v3/payouts", r.URL.Path)
		is.Equal("POST", r.Method)
		is.Equal(PayoutAPIVersion, r.Header.Get("api-version"))
		is.Equal("detail-123", r.Header.Get("idempotency-key"))
		user, pass, ok := r.BasicAuth()
		is.True(ok)
		is.Equal("test-secret", user)
		is.Equal("", pass)

		body, _ := ioutil.ReadAll(r.Body)
		var req PayoutRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal("detail-123", req.ReferenceID)
		is.Equal(RoutingTypeSWIFT, req.Recipient.AccountDetails.RoutingType1)
		is.Equal("CENAIDJA", req.Recipient.AccountDetails.RoutingValue1)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payout_id":"po-mock","status":"ACCEPTED","reference_id":"detail-123"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	resp, err := gw.CreatePayout(&PayoutRequest{
		ReferenceID: "detail-123",
		Recipient: PayoutRecipient{
			Type:         PayoutEntityTypeIndividual,
			GivenName:    "Seller",
			Surname:      "Test",
			Relationship: "SUPPLIER",
			Address: PayoutAddress{
				Country:     "ID",
				City:        "Jakarta",
				StreetLine1: "Jl. Test No. 1",
			},
			AccountDetails: PayoutAccountDetails{
				Currency:          "IDR",
				AccountCountry:    "ID",
				AccountHolderName: "Seller Test",
				AccountNumber:     "1234567890",
				RoutingType1:      RoutingTypeSWIFT,
				RoutingValue1:     "CENAIDJA",
			},
		},
		PayoutDetails: PayoutDetails{
			SourceCurrency:      "IDR",
			SourceAmount:        10000,
			DestinationCurrency: "IDR",
		},
		SourceOfFund: "BUSINESS_REVENUE",
		PurposeCode:  "TRADES",
	}, "detail-123")

	is.NoErr(err)
	is.Equal("po-mock", resp.PayoutID)
	is.Equal("ACCEPTED", resp.Status)
	is.False(resp.ErrorStatus)
}

func TestCreatePayoutErrorStatus(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":"API_VALIDATION_ERROR","message":"recipient.address.country is required"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	resp, err := gw.CreatePayout(&PayoutRequest{ReferenceID: "detail-456"}, "detail-456")
	is.NoErr(err)
	is.True(resp.ErrorStatus)
	is.Equal("API_VALIDATION_ERROR", resp.ErrorCode)
}

// TestCreatePayoutErrorDetails confirms PayoutResponse (previously duplicating ErrorCode/
// ErrorMessage instead of embedding ErrorResponse) now captures the errors[] array too — the
// field-level detail behind the generic "Inputs are failing validation" sentence.
func TestCreatePayoutErrorDetails(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{
			"error_code": "API_VALIDATION_ERROR",
			"message": "Inputs are failing validation. The errors field contains details about which fields are violating validation",
			"errors": [{"path": "recipient.address.postal_code", "message": "must be a valid postal code"}]
		}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	resp, err := gw.CreatePayout(&PayoutRequest{ReferenceID: "detail-789"}, "detail-789")
	is.NoErr(err)
	is.True(resp.ErrorStatus)
	is.Equal(1, len(resp.Errors))
	is.Equal("recipient.address.postal_code", resp.Errors[0].Path)
}
