package moneyout

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (g *Gateway) doJSON(method, path string, idempotencyKey string, reqBody interface{}, resp interface{}) (int, error) {
	var bodyReader *bytes.Buffer
	if reqBody != nil {
		jsonReq, err := json.Marshal(reqBody)
		if err != nil {
			return 0, err
		}
		bodyReader = bytes.NewBuffer(jsonReq)
	} else {
		bodyReader = bytes.NewBuffer(nil)
	}

	var httpRequest *http.Request
	var err error
	if idempotencyKey != "" {
		httpRequest, err = g.Client.newBatchRequest(idempotencyKey, method, path, bodyReader)
	} else {
		httpRequest, err = g.Client.newRequest(method, path, bodyReader)
	}
	if err != nil {
		return 0, err
	}

	res, err := g.Client.httpClient().Do(httpRequest)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, err
	}
	if len(resBody) > 0 && resp != nil {
		if err := json.Unmarshal(resBody, resp); err != nil {
			return res.StatusCode, err
		}
	}
	return res.StatusCode, nil
}

// CreateBatchDisbursement uploads a batch disbursement file.
func (g *Gateway) CreateBatchDisbursement(idempotencyKey string, req *CreateBatchDisbursementRequest) (*CreateBatchDisbursementResponse, error) {
	resp := &CreateBatchDisbursementResponse{}
	httpStatus, err := g.doJSON("POST", "/batch_disbursements", idempotencyKey, req, resp)
	if err != nil {
		return nil, err
	}
	resp.ErrorResponse.markHTTPError(httpStatus)
	return resp, nil
}

// CreatePayout creates a single payout.
func (g *Gateway) CreatePayout(req *CreatePayoutRequest) (*PayoutResponse, error) {
	resp := &PayoutResponse{}
	httpStatus, err := g.doJSON("POST", "/payouts", "", req, resp)
	if err != nil {
		return nil, err
	}
	resp.ErrorResponse.markHTTPError(httpStatus)
	return resp, nil
}

// GetPayout retrieves payout status.
func (g *Gateway) GetPayout(payoutID string) (*PayoutResponse, error) {
	resp := &PayoutResponse{}
	httpStatus, err := g.doJSON("GET", "/payouts/"+payoutID, "", nil, resp)
	if err != nil {
		return nil, err
	}
	resp.ErrorResponse.markHTTPError(httpStatus)
	return resp, nil
}

// VoidPayout voids a payout.
func (g *Gateway) VoidPayout(payoutID string) (*PayoutResponse, error) {
	resp := &PayoutResponse{}
	httpStatus, err := g.doJSON("POST", "/payouts/"+payoutID+"/void", "", nil, resp)
	if err != nil {
		return nil, err
	}
	resp.ErrorResponse.markHTTPError(httpStatus)
	return resp, nil
}
