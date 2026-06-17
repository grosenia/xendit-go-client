package xenditv3

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

func TestCreatePaymentSessionHTTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/sessions", r.URL.Path)
		is.Equal("POST", r.Method)
		is.Equal(DefaultAPIVersion, r.Header.Get("api-version"))
		user, pass, ok := r.BasicAuth()
		is.True(ok)
		is.Equal("test-secret", user)
		is.Equal("", pass)

		body, _ := ioutil.ReadAll(r.Body)
		var req CreatePaymentSessionRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal("ORD-001", req.ReferenceID)
		is.Equal(SessionTypePay, req.SessionType)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_session_id":"ps-mock","status":"ACTIVE","reference_id":"ORD-001"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	session, err := gw.CreatePaymentSession(&CreatePaymentSessionRequest{
		ReferenceID:   "ORD-001",
		SessionType:   SessionTypePay,
		Mode:          SessionModeCardsSessionJS,
		Amount:        50000,
		Currency:      "IDR",
		Country:       "ID",
		CustomerID:    "cust-mock",
		CaptureMethod: CaptureMethodAutomatic,
		CardsSessionJS: &CardsSessionJS{
			SuccessReturnURL: "https://example.com/success",
			FailureReturnURL: "https://example.com/failure",
		},
	})
	is.NoErr(err)
	is.True(!session.ErrorStatus)
	is.Equal("ps-mock", session.PaymentSessionID)
	is.Equal("ACTIVE", session.Status)
}

func TestCreateCustomerHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/customers", r.URL.Path)
		is.Equal("POST", r.Method)

		body, _ := ioutil.ReadAll(r.Body)
		is.True(strings.Contains(string(body), `"given_names":"Reyvin"`))
		is.True(!strings.Contains(string(body), "individual_detail"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"uuid-123","reference_id":"buyer-1","email":"a@b.c"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	customer, err := gw.CreateCustomer(&SessionCustomer{
		ReferenceID:  "buyer-1",
		Type:         "INDIVIDUAL",
		Email:        "a@b.c",
		MobileNumber: "+6281234567890",
		IndividualDetail: &SessionIndividualDetail{
			GivenNames: "Reyvin",
			Surname:    "Test",
		},
	})
	is.NoErr(err)
	is.True(!customer.ErrorStatus)
	is.Equal("cust-uuid-123", customer.ID)
}

func TestCreatePaymentRequestHTTPError(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error_code":"INVALID_TOKEN","message":"token invalid"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	pr, err := gw.CreatePaymentRequest(&CreatePaymentRequestRequest{
		ReferenceID:    "ORD-002",
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

func TestGetPaymentSessionHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/sessions/ps-mock", r.URL.Path)
		is.Equal("GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payment_session_id":"ps-mock","status":"ACTIVE","reference_id":"ORD-001"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	session, err := gw.GetPaymentSession("ps-mock")
	is.NoErr(err)
	is.True(!session.ErrorStatus)
	is.Equal("ps-mock", session.PaymentSessionID)
}

func TestCancelPaymentSessionHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/sessions/ps-mock/cancel", r.URL.Path)
		is.Equal("POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payment_session_id":"ps-mock","status":"CANCELED"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	session, err := gw.CancelPaymentSession("ps-mock")
	is.NoErr(err)
	is.True(!session.ErrorStatus)
	is.Equal("CANCELED", session.Status)
}

func TestGetPaymentRequestHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/v3/payment_requests/pr-mock", r.URL.Path)
		is.Equal("GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payment_request_id":"pr-mock","status":"SUCCEEDED","reference_id":"ORD-002"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	pr, err := gw.GetPaymentRequest("pr-mock")
	is.NoErr(err)
	is.True(!pr.ErrorStatus)
	is.Equal("pr-mock", pr.PaymentRequestID)
}

func TestGetPaymentTokenHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/v3/payment_tokens/pt-mock", r.URL.Path)
		is.Equal("GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payment_token_id":"pt-mock","status":"ACTIVE","token_details":{"masked_card_number":"411111XXXXXX1111"}}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	token, err := gw.GetPaymentToken("pt-mock")
	is.NoErr(err)
	is.True(!token.ErrorStatus)
	is.Equal("411111XXXXXX1111", token.MaskedCardNumber())
}

func TestCancelPaymentTokenHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/v3/payment_tokens/pt-mock/cancel", r.URL.Path)
		is.Equal("POST", r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"payment_token_id":"pt-mock","status":"INACTIVE"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	token, err := gw.CancelPaymentToken("pt-mock")
	is.NoErr(err)
	is.True(!token.ErrorStatus)
	is.Equal("INACTIVE", token.Status)
}

func TestGetCustomerByReferenceIDHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("GET", r.Method)
		is.True(strings.Contains(r.URL.RawQuery, "reference_id=buyer-1"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"id":"uuid-99","reference_id":"buyer-1","email":"a@b.c"}]}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	customer, err := gw.GetCustomerByReferenceID("buyer-1")
	is.NoErr(err)
	is.True(customer != nil)
	is.Equal("cust-uuid-99", customer.ID)
}

func TestGetCustomerByReferenceIDNotFound(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	customer, err := gw.GetCustomerByReferenceID("missing")
	is.NoErr(err)
	is.True(customer == nil)
}

func TestCreatePaymentRequestHTTPSuccess(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/v3/payment_requests", r.URL.Path)
		is.Equal("POST", r.Method)

		body, _ := ioutil.ReadAll(r.Body)
		var req CreatePaymentRequestRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal("pt-mock", req.PaymentTokenID)
		is.Equal(CardOnFileCustomerUnscheduled, req.ChannelProperties.CardOnFileType)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_request_id":"pr-mock","status":"REQUIRES_ACTION","actions":[{"type":"REDIRECT_CUSTOMER","value":"https://3ds.example/auth"}]}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	pr, err := gw.CreatePaymentRequest(&CreatePaymentRequestRequest{
		ReferenceID:    "ORD-002",
		PaymentTokenID: "pt-mock",
		Type:           PaymentRequestTypePay,
		Country:        "ID",
		Currency:       "IDR",
		RequestAmount:  50000,
		ChannelProperties: &PaymentRequestChannelProperties{
			CardOnFileType:      CardOnFileCustomerUnscheduled,
			TransactionSequence: TransactionSequenceSubsequent,
			CardDetails:         &PaymentRequestCardDetails{Cvn: "123"},
		},
	})
	is.NoErr(err)
	is.True(!pr.ErrorStatus)
	is.Equal("pr-mock", pr.PaymentRequestID)
	is.Equal("https://3ds.example/auth", pr.AuthenticationURL())
}

func TestCreatePaymentLinkSessionHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/sessions", r.URL.Path)
		is.Equal("POST", r.Method)

		body, _ := ioutil.ReadAll(r.Body)
		var req CreatePaymentSessionRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal(SessionModePaymentLink, req.Mode)
		is.True(len(req.AllowedPaymentChannels) >= 1)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_session_id":"ps-link","payment_link_url":"https://checkout.xendit.co/web/ps-link","status":"ACTIVE"}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	session, err := gw.CreatePaymentLinkSession(&PaymentLinkSessionRequest{
		ReferenceID:    "ORD-LINK",
		Amount:         62000,
		PayerEmail:     "buyer@example.com",
		PaymentMethods: []string{"BCA", "BNI"},
	})
	is.NoErr(err)
	is.True(!session.ErrorStatus)
	is.Equal("https://checkout.xendit.co/web/ps-link", session.PaymentURL())
}

func TestCreateReusableVirtualAccountHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/v3/payment_requests", r.URL.Path)
		body, _ := ioutil.ReadAll(r.Body)
		var req CreatePaymentRequestRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal(PaymentRequestTypeReusablePaymentCode, req.Type)
		is.Equal(ChannelCodeBCAVirtualAccount, req.ChannelCode)

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"payment_request_id":"pr-va","status":"ACTIVE","actions":[{"type":"PRESENT_TO_CUSTOMER","value":"va-info"}]}`))
	}))
	defer srv.Close()

	client := NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := NewGateway(client)

	pr, err := gw.CreateReusableVirtualAccount(&ReusableVARequest{
		ReferenceID: "seller-va-1",
		BankCode:    "BCA",
		DisplayName: "Toko Seller",
	})
	is.NoErr(err)
	is.True(!pr.ErrorStatus)
	is.Equal("va-info", pr.PresentToCustomerValue())
}
