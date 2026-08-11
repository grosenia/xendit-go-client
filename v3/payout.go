package xenditv3

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// Payout API v3 — replaces legacy POST /batch_disbursements. Unlike the batch endpoint, this is
// one recipient per call; there is no batch/multi-recipient variant of /v3/payouts.
//
// Field shapes below were confirmed against a live Xendit sandbox on 2026-07-28 (public docs did
// not show the Indonesia-specific routing mapping). Domestic Indonesia bank transfers use
// AccountDetails.RoutingType1 = RoutingTypeSWIFT + RoutingValue1 = the destination bank's SWIFT/BIC
// code — there is no bank-code-style routing type for Indonesia in this API.

type PayoutRequest struct {
	ReferenceID   string          `json:"reference_id"`
	Recipient     PayoutRecipient `json:"recipient"`
	PayoutDetails PayoutDetails   `json:"payout_details"`
	SourceOfFund  string          `json:"source_of_fund"`
	PurposeCode   string          `json:"purpose_code"`
	Description   string          `json:"description,omitempty"`
}

type PayoutRecipient struct {
	Type           string               `json:"type"`
	GivenName      string               `json:"given_name,omitempty"`
	Surname        string               `json:"surname,omitempty"`
	BusinessName   string               `json:"business_name,omitempty"`
	Relationship   string               `json:"relationship"`
	Address        PayoutAddress        `json:"address"`
	AccountDetails PayoutAccountDetails `json:"account_details"`
}

type PayoutAddress struct {
	Country       string `json:"country"`
	ProvinceState string `json:"province_state,omitempty"`
	City          string `json:"city"`
	StreetLine1   string `json:"street_line_1"`
	StreetLine2   string `json:"street_line_2,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
}

type PayoutAccountDetails struct {
	Currency          string `json:"currency"`
	AccountCountry    string `json:"account_country"`
	AccountHolderName string `json:"account_holder_name"`
	AccountNumber     string `json:"account_number"`
	RoutingType1      string `json:"routing_type_1"`
	RoutingValue1     string `json:"routing_value_1"`
}

type PayoutDetails struct {
	SourceCurrency      string `json:"source_currency"`
	SourceAmount        int64  `json:"source_amount,omitempty"`
	DestinationCurrency string `json:"destination_currency"`
	DestinationAmount   int64  `json:"destination_amount,omitempty"`
}

type PayoutResponse struct {
	PayoutID             string          `json:"payout_id"`
	Status               string          `json:"status"`
	ReferenceID          string          `json:"reference_id"`
	Type                 string          `json:"type"`
	SourceCurrency       string          `json:"source_currency"`
	SourceAmount         int64           `json:"source_amount"`
	DestinationCurrency  string          `json:"destination_currency"`
	DestinationAmount    int64           `json:"destination_amount"`
	Recipient            PayoutRecipient `json:"recipient"`
	SourceOfFund         string          `json:"source_of_fund"`
	PurposeCode          string          `json:"purpose_code"`
	Description          string          `json:"description"`
	Created              string          `json:"created"`
	Updated              string          `json:"updated"`
	EstimatedArrivalTime string          `json:"estimated_arrival_time"`
	FailureCode          string          `json:"failure_code,omitempty"`
	BusinessID           string          `json:"business_id"`

	ErrorResponse
}

func (g *Gateway) doPayouts(method, path, idempotencyKey string, reqBody interface{}, resp *PayoutResponse) (int, error) {
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

	httpRequest, err := g.Client.newPayoutsRequest(method, path, idempotencyKey, bodyReader)
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
	if resp != nil {
		resp.markHTTPError(res.StatusCode)
	}
	return res.StatusCode, nil
}

// CreatePayout sends a single payout via POST /v3/payouts. idempotencyKey should be a stable
// per-payout identifier (e.g. the caller's disbursement_detail row ID) so retries don't double-pay.
func (g *Gateway) CreatePayout(in *PayoutRequest, idempotencyKey string) (*PayoutResponse, error) {
	resp := &PayoutResponse{}
	_, err := g.doPayouts("POST", createPayoutURL(g.Client), idempotencyKey, in, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPayout fetches current status via GET /v3/payouts/{id}.
func (g *Gateway) GetPayout(payoutID string) (*PayoutResponse, error) {
	resp := &PayoutResponse{}
	_, err := g.doPayouts("GET", getPayoutURL(g.Client, payoutID), "", nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
