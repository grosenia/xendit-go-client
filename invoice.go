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

// ExpireInvoice expires an invoice by invoice ID
// Endpoint: POST https://api.xendit.co/invoices/{invoice_id}/expire!
// According to Xendit docs, the endpoint is /invoices/{invoice_id}/expire! (with !, without /v2/)
// We'll try the documented endpoint first, then fallback to v2 if needed
func (gateway *InvoiceGateway) ExpireInvoice(invoiceID string) (*XenditCreateInvoiceResp, error) {
	resp := XenditCreateInvoiceResp{}
	log := clog.Get()

	// Try endpoint exactly as per Xendit documentation: /invoices/{invoice_id}/expire! (with !)
	// Note: This is different from other invoice APIs which use /v2/invoices
	path := gateway.Client.APIEnvType.String() + "/invoices/" + invoiceID + "/expire!"

	if gateway.Client.LogLevel > 1 {
		// keep minimal visibility when log level > info
		log.Infof("ExpireInvoice: POST %s invoiceID=%s", path, invoiceID)
	}

	httpRequest, err := gateway.Client.NewRequest("POST", path, nil)
	if err != nil {
		log.Errorf("ExpireInvoice: Failed to create request for invoice ID: %s, error: %v", invoiceID, err)
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Errorf("ExpireInvoice: ExecuteRequest failed for invoice ID: %s, error: %v", invoiceID, err)
		return nil, err
	}

	// If endpoint with ! returns 404, try v2 endpoint (for consistency with other invoice APIs)
	if httpStatus == 404 {
		log.Warnf("ExpireInvoice: endpoint with ! returned 404, trying v2 endpoint for invoice ID: %s", invoiceID)
		pathV2 := gateway.Client.APIEnvType.String() + "/v2/invoices/" + invoiceID + "/expire"
		httpRequestV2, errV2 := gateway.Client.NewRequest("POST", pathV2, nil)
		if errV2 == nil {
			httpStatus, err = gateway.Client.ExecuteRequest(httpRequestV2, &resp)
			if err != nil {
				log.Errorf("ExpireInvoice: v2 ExecuteRequest also failed for invoice ID: %s, error: %v", invoiceID, err)
				return nil, err
			}
		}
	}

	// Accept both 200 (OK) and 204 (No Content) as success status codes
	// Some APIs return 204 for successful operations without response body
	if httpStatus != 200 && httpStatus != 204 {
		resp.ErrorStatus = true
		log.Warnf("ExpireInvoice: non-success status %d for invoice ID: %s, response: %+v", httpStatus, invoiceID, resp)
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}
