package xenditgo

import (
	"bytes"
	"encoding/json"

	"github.com/nbs-go/clog"
	_ "github.com/nbs-go/clogrus"
)

// PayoutGateway struct
type PayoutGateway struct {
	Client Client
}

// CreatePayout call create Payout API
func (gateway *InvoiceGateway) CreatePayout(req *XenditCreatePayoutReq) (*XenditCreatePayoutResp, error) {
	log := clog.Get()
	resp := XenditCreatePayoutResp{}
	jsonReq, _ := json.Marshal(req)

	path := gateway.Client.APIEnvType.String() + "/payouts"
	httpRequest, err := gateway.Client.NewRequest("POST", path, bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error create payout ", err)
		return nil, err
	}

	if httpStatus != 200 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}

// GetPayout call create Payout API
func (gateway *InvoiceGateway) GetPayout(payoutID string) (*XenditCreatePayoutResp, error) {
	log := clog.Get()
	resp := XenditCreatePayoutResp{}

	path := gateway.Client.APIEnvType.String() + "/payouts/" + payoutID
	httpRequest, err := gateway.Client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error get payout ", err)
		return nil, err
	}

	if httpStatus != 200 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}

// VoidPayout call create Payout API
func (gateway *InvoiceGateway) VoidPayout(payoutID string) (*XenditCreatePayoutResp, error) {
	log := clog.Get()
	resp := XenditCreatePayoutResp{}

	path := gateway.Client.APIEnvType.String() + "/payouts/" + payoutID + "/void"
	httpRequest, err := gateway.Client.NewRequest("POST", path, nil)

	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error get payout ", err)
		return nil, err
	}

	if httpStatus != 200 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}
