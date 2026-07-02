package xenditv3

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
)

// Skenario 1–6 mengikuti urutan v3/examples (card-session → batch-disburse).
// Setiap skenario punya versi sukses dan gagal agar harapan ErrorStatus / field jelas.

func mockGateway(t *testing.T, handler http.HandlerFunc) *Gateway {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	return NewGateway(client)
}

func TestScenario01_CardSession_Success(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		is.Equal("POST", r.Method)
		is.Equal("/sessions", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_session_id":"ps-ok","status":"ACTIVE","reference_id":"ord-1"}`))
	})

	session, err := gw.CreatePaymentSession(&CreatePaymentSessionRequest{
		ReferenceID:   "ord-1",
		SessionType:   SessionTypePay,
		Mode:          SessionModeCardsSessionJS,
		Amount:        50000,
		Currency:      "IDR",
		Country:       "ID",
		CustomerID:    "cust-abc",
		CaptureMethod: CaptureMethodAutomatic,
		CardsSessionJS: &CardsSessionJS{
			SuccessReturnURL: "https://example.com/success",
			FailureReturnURL: "https://example.com/failure",
		},
	})
	is.NoErr(err)
	is.True(!session.ErrorStatus)
	is.Equal("ps-ok", session.PaymentSessionID)
	is.Equal("ACTIVE", session.Status)
}

func TestScenario01_CardSession_Fail_InvalidCustomer(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":"INVALID_CUSTOMER_ID","message":"customer not found"}`))
	})

	session, err := gw.CreatePaymentSession(&CreatePaymentSessionRequest{
		ReferenceID: "ord-1",
		SessionType: SessionTypePay,
		Mode:        SessionModeCardsSessionJS,
		Amount:      50000,
		Currency:    "IDR",
		Country:     "ID",
		CustomerID:  "cust-bad",
	})
	is.NoErr(err)
	is.True(session.ErrorStatus)
	is.Equal("INVALID_CUSTOMER_ID", session.ErrorCode)
	is.Equal("", session.PaymentSessionID)
}

func TestScenario02_PaymentLink_Success(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req CreatePaymentSessionRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal(SessionModePaymentLink, req.Mode)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_session_id":"ps-link","payment_link_url":"https://checkout.xendit.co/web/ps-link","status":"ACTIVE"}`))
	})

	session, err := gw.CreatePaymentLinkSession(&PaymentLinkSessionRequest{
		ReferenceID:    "inv-1",
		Amount:         62000,
		PayerEmail:     "buyer@example.com",
		PaymentMethods: []string{"BCA"},
	})
	is.NoErr(err)
	is.True(!session.ErrorStatus)
	is.Equal("ps-link", session.PaymentSessionID)
	is.Equal("https://checkout.xendit.co/web/ps-link", session.PaymentURL())
}

func TestScenario02_PaymentLink_Fail_MissingPaymentMethods(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not call API when payment_methods empty")
	})

	_, err := gw.CreatePaymentLinkSession(&PaymentLinkSessionRequest{
		ReferenceID: "inv-1",
		Amount:      62000,
		PayerEmail:  "buyer@example.com",
	})
	is.True(err != nil)
	is.Equal("payment_methods is required", err.Error())
}

func TestScenario02_PaymentLink_Fail_InvalidAmount(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":"API_VALIDATION_ERROR","message":"amount must be positive"}`))
	})

	session, err := gw.CreatePaymentLinkSession(&PaymentLinkSessionRequest{
		ReferenceID:    "inv-1",
		Amount:         -100,
		PayerEmail:     "buyer@example.com",
		PaymentMethods: []string{"BCA"},
	})
	is.NoErr(err)
	is.True(session.ErrorStatus)
	is.Equal("API_VALIDATION_ERROR", session.ErrorCode)
}

func TestScenario02b_PoolVA_Success(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req CreatePaymentRequestRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal(PaymentRequestTypePay, req.Type)
		is.Equal(ChannelCodeBCAVirtualAccount, req.ChannelCode)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_request_id":"pr-pool","status":"REQUIRES_ACTION","channel_properties":{"display_name":"Grosenia Niaga Indonesia"},"actions":[{"type":"PRESENT_TO_CUSTOMER","descriptor":"VIRTUAL_ACCOUNT_NUMBER","value":"8808123456789"}]}`))
	})

	pr, err := gw.CreatePoolVirtualAccount(&PoolVARequest{
		ReferenceID: "ORD-POOL-1",
		BankCode:    "BCA",
		Amount:      93560,
		DisplayName: "Grosenia Niaga Indonesia",
	})
	is.NoErr(err)
	is.True(!pr.ErrorStatus)
	is.Equal("pr-pool", pr.PaymentRequestID)
	is.Equal("8808123456789", pr.VirtualAccountNumber())
}

func TestScenario03_ReusableVA_Success(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		var req CreatePaymentRequestRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal(PaymentRequestTypeReusablePaymentCode, req.Type)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_request_id":"pr-va","status":"ACTIVE","actions":[{"type":"PRESENT_TO_CUSTOMER","value":"8812345678"}]}`))
	})

	pr, err := gw.CreateReusableVirtualAccount(&ReusableVARequest{
		ReferenceID: "va-1",
		BankCode:    "BCA",
		DisplayName: "Toko Seller",
	})
	is.NoErr(err)
	is.True(!pr.ErrorStatus)
	is.Equal("pr-va", pr.PaymentRequestID)
	is.Equal("8812345678", pr.PresentToCustomerValue())
}

func TestScenario03_ReusableVA_Fail_InvalidBank(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":"CHANNEL_UNAVAILABLE","message":"channel not available"}`))
	})

	pr, err := gw.CreateReusableVirtualAccount(&ReusableVARequest{
		ReferenceID: "va-1",
		BankCode:    "BCA",
		DisplayName: "Toko Seller",
	})
	is.NoErr(err)
	is.True(pr.ErrorStatus)
	is.Equal("CHANNEL_UNAVAILABLE", pr.ErrorCode)
}

func TestScenario03_ReusableVA_Fail_UnknownBankCode_Build(t *testing.T) {
	is := is.New(t)
	_, err := BuildReusableVARequest(ReusableVARequest{
		ReferenceID: "va-1",
		BankCode:    "UNKNOWN_BANK",
		DisplayName: "Toko",
	})
	is.True(err != nil)
}

func TestScenario04_PaymentRequest_Success(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_request_id":"pr-ok","status":"REQUIRES_ACTION","actions":[{"type":"REDIRECT_CUSTOMER","value":"https://3ds.example/auth"}]}`))
	})

	pr, err := gw.CreatePaymentRequest(&CreatePaymentRequestRequest{
		ReferenceID:    "ord-2",
		PaymentTokenID: "pt-ok",
		Type:           PaymentRequestTypePay,
		Country:        "ID",
		Currency:       "IDR",
		RequestAmount:  50000,
		ChannelProperties: &PaymentRequestChannelProperties{
			CardDetails: &PaymentRequestCardDetails{Cvn: "123"},
		},
	})
	is.NoErr(err)
	is.True(!pr.ErrorStatus)
	is.Equal("pr-ok", pr.PaymentRequestID)
	is.Equal("https://3ds.example/auth", pr.AuthenticationURL())
}

func TestScenario04_PaymentRequest_Fail_InvalidToken(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":"INVALID_TOKEN","message":"token invalid or expired"}`))
	})

	pr, err := gw.CreatePaymentRequest(&CreatePaymentRequestRequest{
		ReferenceID:    "ord-2",
		PaymentTokenID: "pt-bad",
		Type:           PaymentRequestTypePay,
		Country:        "ID",
		Currency:       "IDR",
		RequestAmount:  50000,
	})
	is.NoErr(err)
	is.True(pr.ErrorStatus)
	is.Equal("INVALID_TOKEN", pr.ErrorCode)
}

func TestScenario05_GetToken_Success(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		is.Equal("GET", r.Method)
		is.Equal("/v3/payment_tokens/pt-ok", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payment_token_id":"pt-ok","status":"ACTIVE","token_details":{"masked_card_number":"411111XXXXXX1111"}}`))
	})

	token, err := gw.GetPaymentToken("pt-ok")
	is.NoErr(err)
	is.True(!token.ErrorStatus)
	is.Equal("pt-ok", token.PaymentTokenID)
	is.Equal("411111XXXXXX1111", token.MaskedCardNumber())
}

func TestScenario05_GetToken_Fail_NotFound(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error_code":"DATA_NOT_FOUND","message":"payment token not found"}`))
	})

	token, err := gw.GetPaymentToken("pt-missing")
	is.NoErr(err)
	is.True(token.ErrorStatus)
	is.Equal("DATA_NOT_FOUND", token.ErrorCode)
}

func TestScenario06_CreateCustomer_Success(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		is.Equal("POST", r.Method)
		is.Equal("/customers", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"uuid-1","reference_id":"buyer-1","email":"a@b.c"}`))
	})

	customer, err := gw.CreateCustomer(&SessionCustomer{
		ReferenceID: "buyer-1",
		Type:        "INDIVIDUAL",
		Email:       "a@b.c",
		IndividualDetail: &SessionIndividualDetail{
			GivenNames: "Test",
			Surname:    "User",
		},
	})
	is.NoErr(err)
	is.True(!customer.ErrorStatus)
	is.Equal("cust-uuid-1", customer.ID)
}

func TestScenario06_CreateCustomer_Fail_Validation(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":"API_VALIDATION_ERROR","message":"given_names is required"}`))
	})

	customer, err := gw.CreateCustomer(&SessionCustomer{
		ReferenceID: "buyer-1",
		Type:        "INDIVIDUAL",
		Email:       "a@b.c",
	})
	is.NoErr(err)
	is.True(customer.ErrorStatus)
	is.Equal("API_VALIDATION_ERROR", customer.ErrorCode)
}

func TestScenario_GetSession_Fail_NotFound(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error_code":"DATA_NOT_FOUND","message":"session not found"}`))
	})

	session, err := gw.GetPaymentSession("ps-missing")
	is.NoErr(err)
	is.True(session.ErrorStatus)
	is.Equal("DATA_NOT_FOUND", session.ErrorCode)
}

func TestScenario_CancelSession_Fail_AlreadyCanceled(t *testing.T) {
	is := is.New(t)
	gw := mockGateway(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":"INVALID_PAYMENT_SESSION","message":"session already canceled"}`))
	})

	session, err := gw.CancelPaymentSession("ps-canceled")
	is.NoErr(err)
	is.True(session.ErrorStatus)
	is.Equal("INVALID_PAYMENT_SESSION", session.ErrorCode)
}
