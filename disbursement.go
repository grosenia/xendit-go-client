package xenditgo

import (
	"bytes"
	"encoding/json"

	"github.com/nbs-go/clog"
	_ "github.com/nbs-go/clogrus"
)

// PayoutGateway struct
type DisbursementGateway struct {
	Client Client
}

// CreatePayout call create Payout API
func (gateway *DisbursementGateway) CreateBatchDisbursement(key string, req *XenditCreateBatchReq) (*XenditCreateBatchResp, error) {
	log := clog.Get()
	resp := XenditCreateBatchResp{}
	jsonReq, _ := json.Marshal(req)

	path := gateway.Client.APIEnvType.String() + "/batch_disbursements"
	httpRequest, err := gateway.Client.NewDisbBatchRequest(key, "POST", path, bytes.NewBuffer(jsonReq))
	// fmt.Println(bytes.NewBuffer(jsonReq).String())

	if err != nil {
		return nil, err
	}
	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error create batch disbursement ", err)
		return nil, err
	}

	if httpStatus != 200 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}
