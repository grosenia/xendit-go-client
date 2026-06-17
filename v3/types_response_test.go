package xenditv3_test

import (
	"testing"

	"github.com/cheekybits/is"
	xenditv3 "github.com/grosenia/xendit-go-client/v3"
)

func TestPaymentTokenMaskedCardNumber(t *testing.T) {
	is := is.New(t)

	fromTokenDetails := &xenditv3.PaymentTokenResponse{
		TokenDetails: &xenditv3.PaymentTokenCardDetails{MaskedCardNumber: "411111XXXXXX1111"},
	}
	is.Equal("411111XXXXXX1111", fromTokenDetails.MaskedCardNumber())

	fromCardDetails := &xenditv3.PaymentTokenResponse{
		CardDetails: &xenditv3.PaymentTokenCardDetails{MaskedCardNumber: "400000XXXXXX1091"},
	}
	is.Equal("400000XXXXXX1091", fromCardDetails.MaskedCardNumber())

	fromChannel := &xenditv3.PaymentTokenResponse{
		ChannelProperties: &xenditv3.PaymentTokenChannelProperties{
			CardDetails: &xenditv3.PaymentTokenCardDetails{MaskedCardNumber: "520000XXXXXX0007"},
		},
	}
	is.Equal("520000XXXXXX0007", fromChannel.MaskedCardNumber())

	is.Equal("", (*xenditv3.PaymentTokenResponse)(nil).MaskedCardNumber())
}

func TestPaymentCaptureCallbackPaymentStatus(t *testing.T) {
	is := is.New(t)

	ok := xenditv3.PaymentCaptureCallback{}
	ok.Data.Status = "CAPTURED"
	is.True(ok.IsPaymentSucceeded())

	fail := xenditv3.PaymentCaptureCallback{}
	fail.Data.Status = "FAILED"
	is.True(!fail.IsPaymentSucceeded())
}
