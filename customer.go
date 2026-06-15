package xenditgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nbs-go/clog"
	_ "github.com/nbs-go/clogrus"
)

// NormalizeXenditCustomerID ensures Payments API v3 customer_id format (cust-{uuid}, min 41 chars).
// Customer API may return bare UUID in the "id" field without the cust- prefix.
func NormalizeXenditCustomerID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	if strings.HasPrefix(id, "cust-") {
		return id
	}
	return "cust-" + id
}

func (gateway *InvoiceGateway) customerRequestBytes(method, path string, reqBody interface{}) (int, []byte, error) {
	log := clog.Get()

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

	httpRequest, err := gateway.Client.NewRequestWithoutAPIVersion(method, path, bodyReader)
	if err != nil {
		return 0, nil, err
	}

	res, err := httpClient.Do(httpRequest)
	if err != nil {
		log.Error("customer api request failed ", err)
		return 0, nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resBody, nil
}

func parseCustomerListResponse(body []byte) ([]XenditCustomerResp, error) {
	if len(body) == 0 {
		return nil, nil
	}

	var wrapped XenditCustomerListResp
	if err := json.Unmarshal(body, &wrapped); err == nil && len(wrapped.Data) > 0 {
		return wrapped.Data, nil
	}

	var list []XenditCustomerResp
	if err := json.Unmarshal(body, &list); err == nil {
		return list, nil
	}

	var single XenditCustomerResp
	if err := json.Unmarshal(body, &single); err == nil && single.ID != "" {
		return []XenditCustomerResp{single}, nil
	}

	return nil, fmt.Errorf("unable to parse customer list response")
}

func markCustomerAPIError(resp *XenditCustomerResp, httpStatus int) {
	if resp == nil {
		return
	}
	if httpStatus != 200 && httpStatus != 201 {
		resp.ErrorStatus = true
		if resp.ErrorCode == "" && resp.ErrorMessage == "" {
			resp.ErrorMessage = fmt.Sprintf("unexpected http status %d", httpStatus)
		}
	} else {
		resp.ErrorStatus = false
	}
}

// GetCustomerByReferenceID returns an existing customer or nil if not found.
// Docs: https://docs.xendit.co/apidocs/get-customers-list
func (gateway *InvoiceGateway) GetCustomerByReferenceID(referenceID string) (*XenditCustomerResp, error) {
	path := gateway.Client.APIEnvType.GetCustomersByReferenceIDURL(referenceID)
	httpStatus, body, err := gateway.customerRequestBytes("GET", path, nil)
	if err != nil {
		return nil, err
	}
	if httpStatus == 404 {
		return nil, nil
	}
	if httpStatus != 200 {
		return nil, fmt.Errorf("get customer by reference_id failed: http %d", httpStatus)
	}

	customers, err := parseCustomerListResponse(body)
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, nil
	}
	customers[0].ID = NormalizeXenditCustomerID(customers[0].ID)
	return &customers[0], nil
}

// CreateCustomer creates an end-customer resource before payment session.
// Docs: https://docs.xendit.co/apidocs/create-customer-request
func (gateway *InvoiceGateway) CreateCustomer(req *XenditPaymentSessionCustomer) (*XenditCustomerResp, error) {
	path := gateway.Client.APIEnvType.CreateCustomerURL()
	httpStatus, body, err := gateway.customerRequestBytes("POST", path, req)
	if err != nil {
		return nil, err
	}

	resp := &XenditCustomerResp{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, resp); err != nil {
			return nil, err
		}
	}
	markCustomerAPIError(resp, httpStatus)
	if resp != nil && !resp.ErrorStatus {
		resp.ID = NormalizeXenditCustomerID(resp.ID)
	}
	return resp, nil
}
