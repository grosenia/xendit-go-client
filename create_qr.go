package xenditgo

import (
	"bytes"
	"encoding/json"

	"github.com/nbs-go/clog"
	_ "github.com/nbs-go/clogrus"
)

// InvoiceGateway struct
type QrcodeGateway struct {
	Client Client
}

// CreateInvoice call create invoice API
func (gateway *QrcodeGateway) CreateQrCode(req *XenditCreteQrcodeReq) (*XenditQrCodeResp, error) {
	log := clog.Get()
	resp := XenditQrCodeResp{}
	jsonReq, _ := json.Marshal(req)

	path := gateway.Client.APIEnvType.String() + "/qr_codes"
	httpRequest, err := gateway.Client.NewRequest("POST", path, bytes.NewBuffer(jsonReq))

	if err != nil {
		return nil, err
	}

	httpStatus, err := gateway.Client.ExecuteRequest(httpRequest, &resp)
	if err != nil {
		log.Error("Error charging ", err)
		return nil, err
	}

	if httpStatus != 201 {
		resp.ErrorStatus = true
	} else {
		resp.ErrorStatus = false
	}

	return &resp, nil
}

func (gateway *QrcodeGateway) PayQrCode(qrId string) (*XenditQrCodeResp, error) {
	log := clog.Get()
	resp := XenditQrCodeResp{}

	path := gateway.Client.APIEnvType.String() + "/qr_codes/" + qrId + "/payments/simulate"
	httpRequest, err := gateway.Client.NewRequest("POST", path, nil)

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
