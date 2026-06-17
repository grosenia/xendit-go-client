//go:build integration

package xenditv3_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cheekybits/is"
	xenditv3 "github.com/grosenia/xendit-go-client/v3"
)

func integrationClient(t *testing.T) (*xenditv3.Gateway, string) {
	t.Helper()
	key := strings.TrimSpace(os.Getenv("XENDIT_SECRET_KEY"))
	if key == "" {
		key = strings.TrimSpace(os.Getenv("KEY_WRITE_MONEY_IN"))
	}
	if key == "" || strings.Contains(key, "xxx") || strings.Contains(key, "REPLACE") {
		t.Skip("set XENDIT_SECRET_KEY (sandbox secret) to run integration tests")
	}

	client := xenditv3.NewClient(key)
	if v := strings.TrimSpace(os.Getenv("XENDIT_API_VERSION")); v != "" {
		client.APIVersion = v
	}
	return xenditv3.NewGateway(client), key
}

func integrationBuyerRef() string {
	if ref := strings.TrimSpace(os.Getenv("XENDIT_BUYER_REFERENCE_ID")); ref != "" {
		return ref
	}
	return "v3-test-" + strconv.FormatInt(time.Now().Unix(), 10)
}

func TestIntegrationCreateCustomerAndPaymentSession(t *testing.T) {
	is := is.New(t)
	gw, _ := integrationClient(t)

	buyerRef := integrationBuyerRef()
	email := strings.TrimSpace(os.Getenv("XENDIT_PAYER_EMAIL"))
	if email == "" {
		email = "v3-integration@grosenia.co.id"
	}

	customer, err := gw.GetCustomerByReferenceID(buyerRef)
	is.NoErr(err)
	if customer == nil {
		customer, err = gw.CreateCustomer(&xenditv3.SessionCustomer{
			ReferenceID:  buyerRef,
			Type:         "INDIVIDUAL",
			Email:        email,
			MobileNumber: "+6281234567890",
			IndividualDetail: &xenditv3.SessionIndividualDetail{
				GivenNames: "V3",
				Surname:    "Integration",
			},
		})
		is.NoErr(err)
		if customer.ErrorStatus {
			t.Fatalf("create customer: %s", customer.Error())
		}
		t.Logf("created customer: %s", customer.ID)
	} else {
		t.Logf("reused customer: %s", customer.ID)
	}

	refID := fmt.Sprintf("v3-it-%d", time.Now().Unix())
	session, err := gw.CreatePaymentSession(&xenditv3.CreatePaymentSessionRequest{
		ReferenceID:            refID,
		SessionType:            xenditv3.SessionTypePay,
		Mode:                   xenditv3.SessionModeCardsSessionJS,
		Amount:                 50000,
		Currency:               "IDR",
		Country:                "ID",
		CustomerID:             customer.ID,
		CaptureMethod:          xenditv3.CaptureMethodAutomatic,
		AllowSavePaymentMethod: xenditv3.AllowSavePaymentMethodOptional,
		CardsSessionJS: &xenditv3.CardsSessionJS{
			SuccessReturnURL: "https://grosenia.co.id/payment/success",
			FailureReturnURL: "https://grosenia.co.id/payment/failure",
		},
	})
	is.NoErr(err)
	if session.ErrorStatus {
		t.Fatalf("create session: %s", session.Error())
	}
	is.True(session.PaymentSessionID != "")
	is.Equal("ACTIVE", session.Status)
	t.Logf("payment_session_id=%s reference_id=%s", session.PaymentSessionID, session.ReferenceID)

	got, err := gw.GetPaymentSession(session.PaymentSessionID)
	is.NoErr(err)
	if got.ErrorStatus {
		t.Fatalf("get session: %s", got.Error())
	}
	is.Equal(session.PaymentSessionID, got.PaymentSessionID)

	canceled, err := gw.CancelPaymentSession(session.PaymentSessionID)
	is.NoErr(err)
	if canceled.ErrorStatus {
		t.Fatalf("cancel session: %s", canceled.Error())
	}
	t.Logf("canceled session status=%s", canceled.Status)
}

func TestIntegrationGetPaymentToken(t *testing.T) {
	tokenID := strings.TrimSpace(os.Getenv("XENDIT_PAYMENT_TOKEN_ID"))
	if tokenID == "" || strings.Contains(tokenID, "REPLACE") {
		t.Skip("set XENDIT_PAYMENT_TOKEN_ID to run saved-card token test")
	}

	is := is.New(t)
	gw, _ := integrationClient(t)

	token, err := gw.GetPaymentToken(tokenID)
	is.NoErr(err)
	if token.ErrorStatus {
		t.Fatalf("get token: %s", token.Error())
	}
	is.Equal(tokenID, token.PaymentTokenID)
	t.Logf("token status=%s masked=%s customer=%s", token.Status, token.MaskedCardNumber(), token.CustomerID)
}

func TestIntegrationCreatePaymentRequestWithSavedCard(t *testing.T) {
	tokenID := strings.TrimSpace(os.Getenv("XENDIT_PAYMENT_TOKEN_ID"))
	cvn := strings.TrimSpace(os.Getenv("XENDIT_CVN"))
	if tokenID == "" || strings.Contains(tokenID, "REPLACE") || cvn == "" {
		t.Skip("set XENDIT_PAYMENT_TOKEN_ID and XENDIT_CVN to run saved-card payment test")
	}

	is := is.New(t)
	gw, _ := integrationClient(t)

	refID := fmt.Sprintf("v3-pr-%d", time.Now().Unix())
	pr, err := gw.CreatePaymentRequest(&xenditv3.CreatePaymentRequestRequest{
		ReferenceID:    refID,
		PaymentTokenID: tokenID,
		Type:           xenditv3.PaymentRequestTypePay,
		Country:        "ID",
		Currency:       "IDR",
		RequestAmount:  50000,
		CaptureMethod:  xenditv3.CaptureMethodAutomatic,
		ChannelProperties: &xenditv3.PaymentRequestChannelProperties{
			CardOnFileType:      xenditv3.CardOnFileCustomerUnscheduled,
			TransactionSequence: xenditv3.TransactionSequenceSubsequent,
			CardDetails:         &xenditv3.PaymentRequestCardDetails{Cvn: cvn},
		},
	})
	is.NoErr(err)
	if pr.ErrorStatus {
		t.Fatalf("create payment request: %s", pr.Error())
	}
	is.True(pr.PaymentRequestID != "")
	t.Logf("payment_request_id=%s status=%s auth=%s", pr.PaymentRequestID, pr.Status, pr.AuthenticationURL())
}

func TestIntegrationCreatePaymentLinkSession(t *testing.T) {
	is := is.New(t)
	gw, _ := integrationClient(t)

	refID := fmt.Sprintf("v3-link-%d", time.Now().Unix())
	session, err := gw.CreatePaymentLinkSession(&xenditv3.PaymentLinkSessionRequest{
		ReferenceID:        refID,
		Amount:             62000,
		PayerEmail:         "v3-link@grosenia.co.id",
		Description:        "v3 payment link test",
		InvoiceDurationSec: 86400,
		PaymentMethods:     []string{"BCA", "MANDIRI"},
	})
	is.NoErr(err)
	if session.ErrorStatus {
		t.Fatalf("create payment link session: %s", session.Error())
	}
	is.True(session.PaymentSessionID != "")
	t.Logf("payment_url=%s session_id=%s", session.PaymentURL(), session.PaymentSessionID)
}

func TestIntegrationCreateReusableVA(t *testing.T) {
	is := is.New(t)
	gw, _ := integrationClient(t)

	refID := fmt.Sprintf("v3-va-%d", time.Now().Unix())
	pr, err := gw.CreateReusableVirtualAccount(&xenditv3.ReusableVARequest{
		ReferenceID: refID,
		BankCode:    "BCA",
		DisplayName: "Grosenia VA Test",
	})
	is.NoErr(err)
	if pr.ErrorStatus {
		t.Fatalf("create reusable va: %s", pr.Error())
	}
	is.True(pr.PaymentRequestID != "")
	t.Logf("payment_request_id=%s channel=%s status=%s", pr.PaymentRequestID, pr.ChannelCode, pr.Status)
}
