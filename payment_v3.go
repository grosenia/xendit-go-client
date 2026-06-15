package xenditgo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/nbs-go/clog"
	_ "github.com/nbs-go/clogrus"
)

func (gateway *InvoiceGateway) executePaymentsV3(method, path string, reqBody interface{}, resp interface{}) (int, error) {
	log := clog.Get()

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

	httpRequest, err := gateway.Client.NewRequest(method, path, bodyReader)
	if err != nil {
		return 0, err
	}
	// Payments API v3 uses Client.PaymentsAPIVersion (credit_card_api_version in Grosenia config),
	// not the legacy invoice api-version (e.g. 2022-07-31).
	apiVersion := gateway.Client.PaymentsAPIVersion
	if apiVersion == "" {
		apiVersion = PaymentsAPIVersion
	}
	httpRequest.Header.Set("api-version", apiVersion)

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, resp)
	if err != nil {
		log.Error("payments v3 request failed ", err)
		return httpStatus, err
	}
	return httpStatus, nil
}

func markPaymentsV3Error(resp interface{}, httpStatus int) {
	switch r := resp.(type) {
	case *XenditPaymentSessionResp:
		if httpStatus != 200 && httpStatus != 201 {
			r.ErrorStatus = true
			if r.ErrorCode == "" && r.ErrorMessage == "" {
				r.ErrorMessage = fmt.Sprintf("unexpected http status %d", httpStatus)
			}
		} else {
			r.ErrorStatus = false
		}
	case *XenditPaymentRequestResp:
		if httpStatus != 200 && httpStatus != 201 {
			r.ErrorStatus = true
			if r.ErrorCode == "" && r.ErrorMessage == "" {
				r.ErrorMessage = fmt.Sprintf("unexpected http status %d", httpStatus)
			}
		} else {
			r.ErrorStatus = false
		}
	case *XenditPaymentTokenResp:
		if httpStatus != 200 && httpStatus != 201 {
			r.ErrorStatus = true
			if r.ErrorCode == "" && r.ErrorMessage == "" {
				r.ErrorMessage = fmt.Sprintf("unexpected http status %d", httpStatus)
			}
		} else {
			r.ErrorStatus = false
		}
	}
}

// CreatePaymentSession creates a Payments API v3 session (PAY, SAVE, one-click CVN).
// Docs: https://docs.xendit.co/apidocs/create-session
func (gateway *InvoiceGateway) CreatePaymentSession(req *XenditCreatePaymentSessionReq) (*XenditPaymentSessionResp, error) {
	resp := &XenditPaymentSessionResp{}
	path := gateway.Client.APIEnvType.CreatePaymentSessionURL()
	httpStatus, err := gateway.executePaymentsV3("POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	markPaymentsV3Error(resp, httpStatus)
	return resp, nil
}

// GetPaymentSession retrieves session status by payment_session_id.
func (gateway *InvoiceGateway) GetPaymentSession(paymentSessionID string) (*XenditPaymentSessionResp, error) {
	resp := &XenditPaymentSessionResp{}
	path := gateway.Client.APIEnvType.GetPaymentSessionURL(paymentSessionID)
	httpStatus, err := gateway.executePaymentsV3("GET", path, nil, resp)
	if err != nil {
		return nil, err
	}
	markPaymentsV3Error(resp, httpStatus)
	return resp, nil
}

// CancelPaymentSession cancels an ACTIVE session.
// Docs: https://docs.xendit.co/apidocs/cancel-session
func (gateway *InvoiceGateway) CancelPaymentSession(paymentSessionID string) (*XenditPaymentSessionResp, error) {
	resp := &XenditPaymentSessionResp{}
	path := gateway.Client.APIEnvType.CancelPaymentSessionURL(paymentSessionID)
	httpStatus, err := gateway.executePaymentsV3("POST", path, nil, resp)
	if err != nil {
		return nil, err
	}
	markPaymentsV3Error(resp, httpStatus)
	return resp, nil
}

// CreatePaymentRequest creates a payment request (one-time, pay-and-save, pay with token).
// Docs: https://docs.xendit.co/apidocs/create-payment-request
func (gateway *InvoiceGateway) CreatePaymentRequest(req *XenditCreatePaymentRequestReq) (*XenditPaymentRequestResp, error) {
	resp := &XenditPaymentRequestResp{}
	path := gateway.Client.APIEnvType.CreatePaymentRequestURL()
	httpStatus, err := gateway.executePaymentsV3("POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	markPaymentsV3Error(resp, httpStatus)
	return resp, nil
}

// GetPaymentRequest retrieves payment request status.
func (gateway *InvoiceGateway) GetPaymentRequest(paymentRequestID string) (*XenditPaymentRequestResp, error) {
	resp := &XenditPaymentRequestResp{}
	path := gateway.Client.APIEnvType.GetPaymentRequestURL(paymentRequestID)
	httpStatus, err := gateway.executePaymentsV3("GET", path, nil, resp)
	if err != nil {
		return nil, err
	}
	markPaymentsV3Error(resp, httpStatus)
	return resp, nil
}

// GetPaymentToken retrieves saved payment token details (masked card).
func (gateway *InvoiceGateway) GetPaymentToken(paymentTokenID string) (*XenditPaymentTokenResp, error) {
	resp := &XenditPaymentTokenResp{}
	path := gateway.Client.APIEnvType.GetPaymentTokenURL(paymentTokenID)
	httpStatus, err := gateway.executePaymentsV3("GET", path, nil, resp)
	if err != nil {
		return nil, err
	}
	markPaymentsV3Error(resp, httpStatus)
	return resp, nil
}

// CancelPaymentToken deactivates a saved payment token.
// Docs: https://docs.xendit.co/apidocs/cancel-payment-token
func (gateway *InvoiceGateway) CancelPaymentToken(paymentTokenID string) (*XenditPaymentTokenResp, error) {
	resp := &XenditPaymentTokenResp{}
	path := gateway.Client.APIEnvType.CancelPaymentTokenURL(paymentTokenID)
	httpStatus, err := gateway.executePaymentsV3("POST", path, nil, resp)
	if err != nil {
		return nil, err
	}
	markPaymentsV3Error(resp, httpStatus)
	return resp, nil
}
