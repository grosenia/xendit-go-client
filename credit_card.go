package xenditgo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/nbs-go/clog"
	_ "github.com/nbs-go/clogrus"
)

// CreateCreditCardCharge calls Xendit credit card charge API.
func (gateway *InvoiceGateway) CreateCreditCardCharge(req *XenditCreateCreditCardChargeReq) (*XenditCreateCreditCardChargeResp, error) {
	log := clog.Get()
	resp := XenditCreateCreditCardChargeResp{}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	path := gateway.Client.APIEnvType.CreateCreditCardChargeURL()
	httpRequest, err := gateway.Client.NewRequest("POST", path, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error creating credit card charge ", err)
		return nil, err
	}

	if httpStatus != 200 && httpStatus != 201 {
		resp.ErrorStatus = true
		if resp.ErrorCode == "" && resp.ErrorMessage == "" {
			resp.ErrorMessage = fmt.Sprintf("unexpected http status %d", httpStatus)
		}
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}
