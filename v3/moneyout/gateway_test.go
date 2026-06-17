package moneyout_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/grosenia/xendit-go-client/v3/moneyout"
)

func TestCreateBatchDisbursementHTTP(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/batch_disbursements", r.URL.Path)
		is.Equal("POST", r.Method)
		is.Equal("idem-key-1", r.Header.Get("X-IDEMPOTENCY-KEY"))

		body, _ := ioutil.ReadAll(r.Body)
		var req moneyout.CreateBatchDisbursementRequest
		is.NoErr(json.Unmarshal(body, &req))
		is.Equal("batch-ref-1", req.Reference)
		is.Equal(1, len(req.Disbursements))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"bd-1","reference":"batch-ref-1","status":"UPLOADING"}`))
	}))
	defer srv.Close()

	client := moneyout.NewClient("test-secret")
	client.BaseURL = srv.URL
	gw := moneyout.NewGateway(client)

	resp, err := gw.CreateBatchDisbursement("idem-key-1", &moneyout.CreateBatchDisbursementRequest{
		Reference: "batch-ref-1",
		Disbursements: []moneyout.DisbursementItem{{
			ExternalID:        "item-1",
			Amount:            20000,
			BankCode:          "BCA",
			BankAccountName:   "Test User",
			BankAccountNumber: "1234567890",
			Description:       "test",
		}},
	})
	is.NoErr(err)
	is.True(!resp.ErrorStatus)
	is.Equal("UPLOADING", resp.Status)
}

func TestCreateBatchDisbursementHTTP_Fail_Unauthorized(t *testing.T) {
	is := is.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error_code":"INVALID_API_KEY","message":"API key format is invalid"}`))
	}))
	defer srv.Close()

	client := moneyout.NewClient("bad-key")
	client.BaseURL = srv.URL
	gw := moneyout.NewGateway(client)

	resp, err := gw.CreateBatchDisbursement("idem-fail", &moneyout.CreateBatchDisbursementRequest{
		Reference: "batch-fail",
		Disbursements: []moneyout.DisbursementItem{{
			ExternalID:        "item-1",
			Amount:            20000,
			BankCode:          "BCA",
			BankAccountName:   "Test",
			BankAccountNumber: "123",
		}},
	})
	is.NoErr(err)
	is.True(resp.ErrorStatus)
	is.Equal("INVALID_API_KEY", resp.ErrorCode)
}

