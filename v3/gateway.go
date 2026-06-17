package xenditv3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (g *Gateway) doPayments(method, path string, reqBody interface{}, resp interface{}) (int, error) {
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

	httpRequest, err := g.Client.newPaymentsRequest(method, path, bodyReader)
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

func (g *Gateway) doCustomer(method, path string, reqBody interface{}) (int, []byte, error) {
	var bodyReader *bytes.Buffer
	if reqBody != nil {
		jsonReq, err := json.Marshal(reqBody)
		if err != nil {
			return 0, nil, err
		}
		bodyReader = bytes.NewBuffer(jsonReq)
	} else {
		bodyReader = bytes.NewBuffer(nil)
	}

	httpRequest, err := g.Client.newCustomerRequest(method, path, bodyReader)
	if err != nil {
		return 0, nil, err
	}

	res, err := g.Client.httpClient().Do(httpRequest)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resBody, nil
}

func markSessionError(resp *PaymentSessionResponse, httpStatus int) {
	if resp != nil {
		resp.ErrorResponse.markHTTPError(httpStatus)
	}
}

func markPaymentRequestError(resp *PaymentRequestResponse, httpStatus int) {
	if resp != nil {
		resp.ErrorResponse.markHTTPError(httpStatus)
	}
}

func markPaymentTokenError(resp *PaymentTokenResponse, httpStatus int) {
	if resp != nil {
		resp.ErrorResponse.markHTTPError(httpStatus)
	}
}

// CreatePaymentSession creates a Payments API v3 session.
func (g *Gateway) CreatePaymentSession(req *CreatePaymentSessionRequest) (*PaymentSessionResponse, error) {
	resp := &PaymentSessionResponse{}
	httpStatus, err := g.doPayments("POST", paymentSessionURL(g.Client), req, resp)
	if err != nil {
		return nil, err
	}
	markSessionError(resp, httpStatus)
	return resp, nil
}

// GetPaymentSession retrieves session status.
func (g *Gateway) GetPaymentSession(paymentSessionID string) (*PaymentSessionResponse, error) {
	resp := &PaymentSessionResponse{}
	httpStatus, err := g.doPayments("GET", getPaymentSessionURL(g.Client, paymentSessionID), nil, resp)
	if err != nil {
		return nil, err
	}
	markSessionError(resp, httpStatus)
	return resp, nil
}

// CancelPaymentSession cancels an ACTIVE session.
func (g *Gateway) CancelPaymentSession(paymentSessionID string) (*PaymentSessionResponse, error) {
	resp := &PaymentSessionResponse{}
	httpStatus, err := g.doPayments("POST", cancelPaymentSessionURL(g.Client, paymentSessionID), nil, resp)
	if err != nil {
		return nil, err
	}
	markSessionError(resp, httpStatus)
	return resp, nil
}

// CreatePaymentRequest creates a payment request (one-time, pay-and-save, pay with token).
func (g *Gateway) CreatePaymentRequest(req *CreatePaymentRequestRequest) (*PaymentRequestResponse, error) {
	resp := &PaymentRequestResponse{}
	httpStatus, err := g.doPayments("POST", createPaymentRequestURL(g.Client), req, resp)
	if err != nil {
		return nil, err
	}
	markPaymentRequestError(resp, httpStatus)
	return resp, nil
}

// GetPaymentRequest retrieves payment request status.
func (g *Gateway) GetPaymentRequest(paymentRequestID string) (*PaymentRequestResponse, error) {
	resp := &PaymentRequestResponse{}
	httpStatus, err := g.doPayments("GET", getPaymentRequestURL(g.Client, paymentRequestID), nil, resp)
	if err != nil {
		return nil, err
	}
	markPaymentRequestError(resp, httpStatus)
	return resp, nil
}

// GetPaymentToken retrieves saved payment token details.
func (g *Gateway) GetPaymentToken(paymentTokenID string) (*PaymentTokenResponse, error) {
	resp := &PaymentTokenResponse{}
	httpStatus, err := g.doPayments("GET", getPaymentTokenURL(g.Client, paymentTokenID), nil, resp)
	if err != nil {
		return nil, err
	}
	markPaymentTokenError(resp, httpStatus)
	return resp, nil
}

// CancelPaymentToken deactivates a saved payment token.
func (g *Gateway) CancelPaymentToken(paymentTokenID string) (*PaymentTokenResponse, error) {
	resp := &PaymentTokenResponse{}
	httpStatus, err := g.doPayments("POST", cancelPaymentTokenURL(g.Client, paymentTokenID), nil, resp)
	if err != nil {
		return nil, err
	}
	markPaymentTokenError(resp, httpStatus)
	return resp, nil
}

func parseCustomerList(body []byte) ([]CustomerResponse, error) {
	if len(body) == 0 {
		return nil, nil
	}
	var wrapped customerListResponse
	if err := json.Unmarshal(body, &wrapped); err == nil && len(wrapped.Data) > 0 {
		return wrapped.Data, nil
	}
	var list []CustomerResponse
	if err := json.Unmarshal(body, &list); err == nil {
		return list, nil
	}
	var single CustomerResponse
	if err := json.Unmarshal(body, &single); err == nil && single.ID != "" {
		return []CustomerResponse{single}, nil
	}
	return nil, fmt.Errorf("unable to parse customer list response")
}
