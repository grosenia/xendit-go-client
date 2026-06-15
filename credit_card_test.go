package xenditgo_test

import (
	"testing"

	"github.com/cheekybits/is"
	xenditgo "github.com/grosenia/xendit-go-client"
)

func TestCreateCreditCardChargeURL(t *testing.T) {
	is := is.New(t)
	is.Equal("https://api.xendit.co/credit_card_charges", xenditgo.Sandbox.CreateCreditCardChargeURL())
}

func TestCreditCardCallbackHelpers(t *testing.T) {
	is := is.New(t)

	payload := xenditgo.XenditCreditCardCallback{
		Event: "credit_card.succeeded",
		Data: struct {
			ID         string `json:"id"`
			ExternalID string `json:"external_id"`
			Status     string `json:"status"`
		}{
			ID:         "charge-123",
			ExternalID: "ORD-001",
			Status:     "CAPTURED",
		},
	}

	is.Equal("charge-123", payload.ChargeID())
	is.Equal("ORD-001", payload.OrderNo())
	is.Equal("CAPTURED", payload.ChargeStatus())
}
