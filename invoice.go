package xenditgo

import (
	"bytes"
	"encoding/json"

	"github.com/nbs-go/clog"
	_ "github.com/nbs-go/clogrus"
)

// InvoiceGateway struct
type InvoiceGateway struct {
	Client Client
}

// CreateInvoice call create invoice API
func (gateway *InvoiceGateway) CreateInvoice(req *XenditCreateInvoiceReq) (*XenditCreateInvoiceResp, error) {
	log := clog.Get()
	resp := XenditCreateInvoiceResp{}
	jsonReq, _ := json.Marshal(req)

	path := gateway.Client.APIEnvType.String() + "/v2/invoices"
	httpRequest, err := gateway.Client.NewRequest("POST", path, bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error charging ", err)
		return nil, err
	}

	if httpStatus != 200 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}

// CreateFixedVa call create fixed va API
func (gateway *InvoiceGateway) CreateFixedVa(req *XenditCreateFixedVaReq) (*XenditCreateFixedVaResp, error) {
	log := clog.Get()
	resp := XenditCreateFixedVaResp{}
	jsonReq, _ := json.Marshal(req)

	path := gateway.Client.APIEnvType.String() + "/callback_virtual_accounts"
	httpRequest, err := gateway.Client.NewRequest("POST", path, bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error charging ", err)
		return nil, err
	}

	if httpStatus != 200 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}

func (gateway *InvoiceGateway) UpdateFixedVa(id string, req *XenditUpdateFixedVaReq) (*XenditCreateFixedVaResp, error) {
	log := clog.Get()
	resp := XenditCreateFixedVaResp{}
	jsonReq, _ := json.Marshal(req)

	path := gateway.Client.APIEnvType.String() + "/callback_virtual_accounts/" + id
	httpRequest, err := gateway.Client.NewRequest("PATCH", path, bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error charging ", err)
		return nil, err
	}

	if httpStatus != 200 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}

func (gateway *InvoiceGateway) GetFixedVa(id string) (*XenditCreateFixedVaResp, error) {
	log := clog.Get()
	resp := XenditCreateFixedVaResp{}

	path := gateway.Client.APIEnvType.String() + "/callback_virtual_accounts/" + id
	httpRequest, err := gateway.Client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	_, err = gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error charging ", err)
		return nil, err
	}

	return &resp, nil
}

func (gateway *InvoiceGateway) CreatePaymentMethod(req *XenditCreatePaymentMethodReq) (*XenditPaymentMethodResp, error) {
	log := clog.Get()
	resp := XenditPaymentMethodResp{}
	jsonReq, _ := json.Marshal(req)

	path := gateway.Client.APIEnvType.String() + "/v2/payment_methods"
	httpRequest, err := gateway.Client.NewRequest("POST", path, bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error charging ", err)
		return nil, err
	}

	if httpStatus != 200 && httpStatus != 201 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}
