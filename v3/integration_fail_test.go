//go:build integration

package xenditv3_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cheekybits/is"
	xenditv3 "github.com/grosenia/xendit-go-client/v3"
)

// Integration fail — harapan: HTTP call sukses, tapi ErrorStatus=true + error_code dari Xendit.

func TestIntegrationFail_InvalidAPIKey(t *testing.T) {
	is := is.New(t)

	client := xenditv3.NewClient("xnd_development_INVALID_KEY_FOR_TEST")
	client.APIVersion = "2024-11-11"
	gw := xenditv3.NewGateway(client)

	session, err := gw.CreatePaymentLinkSession(&xenditv3.PaymentLinkSessionRequest{
		ReferenceID:    fmt.Sprintf("fail-key-%d", time.Now().Unix()),
		Amount:         50000,
		PayerEmail:     "fail@grosenia.co.id",
		PaymentMethods: []string{"BCA"},
	})
	is.NoErr(err)
	is.True(session != nil)
	is.True(session.ErrorStatus)
	is.True(session.ErrorCode != "")
	t.Logf("expected fail: [%s] %s", session.ErrorCode, session.ErrorMessage)
}

func TestIntegrationFail_InvalidPaymentToken(t *testing.T) {
	is := is.New(t)
	gw, _ := integrationClient(t)

	token, err := gw.GetPaymentToken("pt_INVALID_TOKEN_DOES_NOT_EXIST")
	is.NoErr(err)
	is.True(token != nil)
	is.True(token.ErrorStatus)
	is.True(token.ErrorCode != "")
	t.Logf("expected fail: [%s] %s", token.ErrorCode, token.ErrorMessage)
}

func TestIntegrationFail_PaymentRequestMissingCVN(t *testing.T) {
	is := is.New(t)
	gw, _ := integrationClient(t)

	pr, err := gw.CreatePaymentRequest(&xenditv3.CreatePaymentRequestRequest{
		ReferenceID:    fmt.Sprintf("fail-cvn-%d", time.Now().Unix()),
		PaymentTokenID: "pt_INVALID_TOKEN_DOES_NOT_EXIST",
		Type:           xenditv3.PaymentRequestTypePay,
		Country:        "ID",
		Currency:       "IDR",
		RequestAmount:  50000,
	})
	is.NoErr(err)
	is.True(pr != nil)
	is.True(pr.ErrorStatus)
	is.True(pr.ErrorCode != "")
	t.Logf("expected fail: [%s] %s", pr.ErrorCode, pr.ErrorMessage)
}

func TestIntegrationFail_ReusableVA_InvalidDisplayName(t *testing.T) {
	is := is.New(t)
	gw, _ := integrationClient(t)

	pr, err := gw.CreateReusableVirtualAccount(&xenditv3.ReusableVARequest{
		ReferenceID: fmt.Sprintf("fail-va-%d", time.Now().Unix()),
		BankCode:    "BCA",
		DisplayName: "",
	})
	is.NoErr(err)
	is.True(pr != nil)
	is.True(pr.ErrorStatus)
	is.True(pr.ErrorCode != "")
	t.Logf("expected fail: [%s] %s", pr.ErrorCode, pr.ErrorMessage)
}
