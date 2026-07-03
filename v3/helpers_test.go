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
	is.True(req.ChannelProperties.VerificationData == nil)
}

func TestBuildReusableVARequest_WithVerificationData(t *testing.T) {
	is := is.New(t)

	req, err := BuildReusableVARequest(ReusableVARequest{
		ReferenceID: "seller-1",
		BankCode:    "BRI",
		DisplayName: "Toko ABC",
		Verification: &ReusableVAVerification{
			CustomerName:           "Toko ABC",
			AcceptedNameVariations: []string{"Toko A B C", "TOKO ABC"},
			AllowedBankAccounts: []PaymentRequestBankAccount{
				{BankName: "BCA", AccountNumber: "1234567890", AccountName: "Toko ABC"},
			},
		},
	})
	is.NoErr(err)
	is.Equal(ChannelCodeBRIVirtualAccount, req.ChannelCode)
	is.True(req.ChannelProperties.VerificationData != nil)
	is.Equal("Toko ABC", req.ChannelProperties.VerificationData.CustomerName)
	is.Equal(2, len(req.ChannelProperties.VerificationData.AcceptedNameVariations))
	is.Equal(1, len(req.ChannelProperties.VerificationData.AllowedBankAccounts))
	is.Equal("BCA", req.ChannelProperties.VerificationData.AllowedBankAccounts[0].BankName)
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

func TestBuildPoolVARequest(t *testing.T) {
	is := is.New(t)

	req, err := BuildPoolVARequest(PoolVARequest{
		ReferenceID: "ORD-POOL-1",
		BankCode:    "BCA",
		Amount:      93560,
		DisplayName: "Grosenia Niaga Indonesia",
		Description: "Test order",
	})
	is.NoErr(err)
	is.Equal(PaymentRequestTypePay, req.Type)
	is.Equal(ChannelCodeBCAVirtualAccount, req.ChannelCode)
	is.Equal(93560.0, req.RequestAmount)
	is.Equal("Grosenia Niaga Indonesia", req.ChannelProperties.DisplayName)
}

func TestVirtualAccountNumber(t *testing.T) {
	is := is.New(t)

	resp := PaymentRequestResponse{
		Actions: []PaymentRequestAction{
			{Type: "PRESENT_TO_CUSTOMER", Descriptor: "VIRTUAL_ACCOUNT_NUMBER", Value: "8812345678"},
		},
	}
	is.Equal("8812345678", resp.VirtualAccountNumber())
}
