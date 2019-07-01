package xenditgo

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

// InvoiceGateway struct
type InvoiceGateway struct {
	Client Client
}

// Call : basic call
func (gateway *InvoiceGateway) Call(method, path string, body io.Reader, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = gateway.Client.APIEnvType.String() + path
	return gateway.Client.Call(method, path, body, v)
}

// CreateInvoice call create invoice API
func (gateway *InvoiceGateway) CreateInvoice(req *XenditCreateInvoiceReq) (XenditCreateInvoiceResp, error) {
	resp := XenditCreateInvoiceResp{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "/v2/invoices", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error charging: ", err)
		return resp, err
	}

	if resp.Status != "" {
		gateway.Client.Logger.Println(resp.Status)
	}
	return resp, nil
}

/*
// Charge : Perform transaction using ChargeReq
func (gateway *InvoiceGateway) Charge(req *ChargeReq) (Response, error) {
	resp := Response{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "v2/charge", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error charging: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// PreauthCard : Perform authorized transactions using ChargeReq
func (gateway *InvoiceGateway) PreauthCard(req *ChargeReq) (Response, error) {
	req.CreditCard.Type = "authorize"
	return gateway.Charge(req)
}

// CaptureCard : Capture an authorized transaction for card payment
func (gateway *InvoiceGateway) CaptureCard(req *CaptureReq) (Response, error) {
	resp := Response{}
	jsonReq, _ := json.Marshal(req)

	err := gateway.Call("POST", "v2/capture", bytes.NewBuffer(jsonReq), &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error capturing: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Approve : Approve order using order ID
func (gateway *InvoiceGateway) Approve(orderID string) (Response, error) {
	resp := Response{}

	err := gateway.Call("POST", "v2/"+orderID+"/approve", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Cancel : Cancel order using order ID
func (gateway *InvoiceGateway) Cancel(orderID string) (Response, error) {
	resp := Response{}

	err := gateway.Call("POST", "v2/"+orderID+"/cancel", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Expire : change order status to expired using order ID
func (gateway *InvoiceGateway) Expire(orderID string) (Response, error) {
	resp := Response{}

	err := gateway.Call("POST", "v2/"+orderID+"/expire", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}

// Status : get order status using order ID
func (gateway *InvoiceGateway) Status(orderID string) (Response, error) {
	resp := Response{}

	err := gateway.Call("GET", "v2/"+orderID+"/status", nil, &resp)
	if err != nil {
		gateway.Client.Logger.Println("Error approving: ", err)
		return resp, err
	}

	if resp.StatusMessage != "" {
		gateway.Client.Logger.Println(resp.StatusMessage)
	}

	return resp, nil
}
*/
