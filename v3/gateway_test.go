package xenditv3_test

import (
	"testing"

	"github.com/cheekybits/is"
	xenditv3 "github.com/grosenia/xendit-go-client/v3"
)

func TestPaymentCaptureCallback(t *testing.T) {
	is := is.New(t)
	payload := xenditv3.PaymentCaptureCallback{
		Event: xenditv3.WebhookEventPaymentCapture,
	}
	payload.Data.PaymentRequestID = "pr-1"
	payload.Data.ReferenceID = "ORD-1"
	payload.Data.Status = "SUCCEEDED"

	is.Equal("pr-1", payload.PaymentRequestID())
	is.Equal("ORD-1", payload.ReferenceID())
	is.True(payload.IsPaymentSucceeded())
}

func TestPaymentTokenCallback(t *testing.T) {
	is := is.New(t)
	payload := xenditv3.PaymentTokenCallback{Event: xenditv3.WebhookEventPaymentTokenActivated}
	payload.Data.PaymentTokenID = "pt-1"
	payload.Data.TokenDetails = &xenditv3.PaymentTokenCardDetails{MaskedCardNumber: "411111XXXXXX1111"}

	is.Equal("pt-1", payload.PaymentTokenID())
	is.Equal("411111XXXXXX1111", payload.MaskedCardNumber())
}

func TestPaymentRequestAuthenticationURL(t *testing.T) {
	is := is.New(t)
	resp := xenditv3.PaymentRequestResponse{
		Actions: []xenditv3.PaymentRequestAction{
			{Type: "REDIRECT_CUSTOMER", Value: "https://3ds.example/auth"},
		},
	}
	is.Equal("https://3ds.example/auth", resp.AuthenticationURL())
}

func TestPaymentSessionCompletedCallback(t *testing.T) {
	is := is.New(t)
	payload := xenditv3.PaymentSessionCompletedCallback{Event: xenditv3.WebhookEventPaymentSessionCompleted}
	payload.Data.PaymentSessionID = "ps-1"
	payload.Data.ReferenceID = "ORD-1"
	payload.Data.PaymentTokenID = "pt-1"
	payload.Data.Status = "COMPLETED"

	is.Equal("ps-1", payload.PaymentSessionID())
	is.Equal("ORD-1", payload.ReferenceID())
	is.Equal("pt-1", payload.PaymentTokenID())
	is.True(payload.IsCompleted())
}

func TestNormalizeCustomerID(t *testing.T) {
	is := is.New(t)
	is.Equal("cust-abc", xenditv3.NormalizeCustomerID("abc"))
	is.Equal("cust-abc", xenditv3.NormalizeCustomerID("cust-abc"))
}
