package xenditv3

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestMapLegacyBankCode(t *testing.T) {
	is := is.New(t)

	ch, err := MapLegacyBankCode("BCA")
	is.NoErr(err)
	is.Equal(ChannelCodeBCAVirtualAccount, ch)

	ch, err = MapLegacyBankCode("BCA_VIRTUAL_ACCOUNT")
	is.NoErr(err)
	is.Equal(ChannelCodeBCAVirtualAccount, ch)

	_, err = MapLegacyBankCode("UNKNOWN")
	is.True(err != nil)
}

func TestBuildPaymentLinkSessionRequest(t *testing.T) {
	is := is.New(t)

	req, err := BuildPaymentLinkSessionRequest(PaymentLinkSessionRequest{
		ReferenceID:    "ORD-1",
		Amount:         62000,
		PayerEmail:     "buyer@example.com",
		Description:    "Test order",
		PaymentMethods: []string{"BCA", "MANDIRI"},
		InvoiceDurationSec: 3600,
	})
	is.NoErr(err)
	is.Equal(SessionModePaymentLink, req.Mode)
	is.Equal(2, len(req.AllowedPaymentChannels))
	is.True(req.ExpiresAt != "")
	is.Equal("buyer@example.com", req.Customer.Email)
}

func TestBuildReusableVARequest(t *testing.T) {
	is := is.New(t)

	req, err := BuildReusableVARequest(ReusableVARequest{
		ReferenceID: "seller-1",
		BankCode:    "BCA",
		DisplayName: "Toko ABC",
	})
	is.NoErr(err)
	is.Equal(PaymentRequestTypeReusablePaymentCode, req.Type)
	is.Equal(ChannelCodeBCAVirtualAccount, req.ChannelCode)
	is.Equal("Toko ABC", req.ChannelProperties.DisplayName)
}

func TestPaymentSessionPaymentURL(t *testing.T) {
	is := is.New(t)
	resp := &PaymentSessionResponse{PaymentLinkURL: "https://checkout.xendit.co/web/abc"}
	is.Equal("https://checkout.xendit.co/web/abc", resp.PaymentURL())
}

func TestPresentToCustomerValue(t *testing.T) {
	is := is.New(t)
	resp := PaymentRequestResponse{
		Actions: []PaymentRequestAction{{Type: "PRESENT_TO_CUSTOMER", Value: "va-payload"}},
	}
	is.Equal("va-payload", resp.PresentToCustomerValue())
}
