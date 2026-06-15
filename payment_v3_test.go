package xenditgo_test

import (
	"testing"

	"github.com/cheekybits/is"
	xenditgo "github.com/grosenia/xendit-go-client"
)

func TestPaymentV3URLs(t *testing.T) {
	is := is.New(t)
	env := xenditgo.Sandbox

	is.Equal("https://api.xendit.co/sessions", env.CreatePaymentSessionURL())
	is.Equal("https://api.xendit.co/sessions/ps-test123", env.GetPaymentSessionURL("ps-test123"))
	is.Equal("https://api.xendit.co/sessions/ps-test123/cancel", env.CancelPaymentSessionURL("ps-test123"))
	is.Equal("https://api.xendit.co/v3/payment_requests", env.CreatePaymentRequestURL())
	is.Equal("https://api.xendit.co/v3/payment_requests/pr-test", env.GetPaymentRequestURL("pr-test"))
	is.Equal("https://api.xendit.co/v3/payment_tokens/pt-test", env.GetPaymentTokenURL("pt-test"))
	is.Equal("https://api.xendit.co/v3/payment_tokens/pt-test/cancel", env.CancelPaymentTokenURL("pt-test"))
}

func TestPaymentCaptureCallbackHelpers(t *testing.T) {
	is := is.New(t)

	payload := xenditgo.XenditPaymentCaptureCallback{
		Event: xenditgo.WebhookEventPaymentCapture,
		Data: struct {
			PaymentID        string  `json:"payment_id"`
			PaymentRequestID string  `json:"payment_request_id"`
			ReferenceID      string  `json:"reference_id"`
			Status           string  `json:"status"`
			Currency         string  `json:"currency"`
			RequestAmount    float64 `json:"request_amount"`
			CaptureAmount    float64 `json:"capture_amount"`
			PaymentTokenID   string  `json:"payment_token_id"`
		}{
			PaymentRequestID: "pr-123",
			ReferenceID:      "ORD-001",
			Status:           "SUCCEEDED",
		},
	}

	is.Equal("pr-123", payload.PaymentRequestID())
	is.Equal("ORD-001", payload.OrderNo())
	is.Equal("SUCCEEDED", payload.PaymentStatus())
	is.True(payload.IsPaymentSucceeded())
}

func TestPaymentTokenCallbackHelpers(t *testing.T) {
	is := is.New(t)

	payload := xenditgo.XenditPaymentTokenCallback{
		Event: xenditgo.WebhookEventPaymentTokenActivated,
		Data: struct {
			PaymentTokenID string                         `json:"payment_token_id"`
			ReferenceID    string                         `json:"reference_id"`
			Status         string                         `json:"status"`
			CustomerID     string                         `json:"customer_id"`
			ChannelCode    string                         `json:"channel_code"`
			TokenDetails   *xenditgo.XenditPaymentTokenCardDetails `json:"token_details"`
			CardDetails    *xenditgo.XenditPaymentTokenCardDetails `json:"card_details"`
		}{
			PaymentTokenID: "pt-abc",
			TokenDetails: &xenditgo.XenditPaymentTokenCardDetails{
				MaskedCardNumber: "400000XXXXXX1091",
				CardBrand:        "VISA",
			},
		},
	}

	is.Equal("pt-abc", payload.PaymentTokenID())
	is.Equal("400000XXXXXX1091", payload.MaskedCardNumber())
}

func TestPaymentRequestAuthenticationURL(t *testing.T) {
	is := is.New(t)

	resp := xenditgo.XenditPaymentRequestResp{
		Actions: []xenditgo.XenditPaymentRequestAction{
			{Type: "REDIRECT_CUSTOMER", Value: "https://3ds.example/auth"},
		},
	}
	is.Equal("https://3ds.example/auth", resp.AuthenticationURL())
}

func TestPaymentTokenMaskedCardNumber(t *testing.T) {
	is := is.New(t)

	resp := &xenditgo.XenditPaymentTokenResp{
		TokenDetails: &xenditgo.XenditPaymentTokenCardDetails{
			MaskedCardNumber: "411111XXXXXX1111",
		},
	}
	is.Equal("411111XXXXXX1111", resp.MaskedCardNumber())
}
