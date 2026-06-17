package xenditv3

import (
	"encoding/json"
	"fmt"
	"strings"
)

// NormalizeCustomerID ensures Payments API v3 customer_id format (cust-{uuid}).
func NormalizeCustomerID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	if strings.HasPrefix(id, "cust-") {
		return id
	}
	return "cust-" + id
}

// GetCustomerByReferenceID returns an existing customer or nil if not found.
func (g *Gateway) GetCustomerByReferenceID(referenceID string) (*CustomerResponse, error) {
	httpStatus, body, err := g.doCustomer("GET", getCustomersByReferenceIDURL(g.Client, referenceID), nil)
	if err != nil {
		return nil, err
	}
	if httpStatus == 404 {
		return nil, nil
	}
	if httpStatus != 200 {
		return nil, fmt.Errorf("get customer by reference_id failed: http %d", httpStatus)
	}

	customers, err := parseCustomerList(body)
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, nil
	}
	customers[0].ID = NormalizeCustomerID(customers[0].ID)
	return &customers[0], nil
}

// CreateCustomer creates an end-customer resource.
func (g *Gateway) CreateCustomer(req *SessionCustomer) (*CustomerResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("customer request is required")
	}
	createReq := CustomerCreateRequest{
		ReferenceID:  req.ReferenceID,
		Type:         req.Type,
		Email:        req.Email,
		MobileNumber: req.MobileNumber,
	}
	if req.IndividualDetail != nil {
		createReq.GivenNames = req.IndividualDetail.GivenNames
		createReq.Surname = req.IndividualDetail.Surname
	}

	httpStatus, body, err := g.doCustomer("POST", createCustomerURL(g.Client), createReq)
	if err != nil {
		return nil, err
	}

	resp := &CustomerResponse{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, resp); err != nil {
			return nil, err
		}
	}
	resp.ErrorResponse.markHTTPError(httpStatus)
	if !resp.ErrorStatus {
		resp.ID = NormalizeCustomerID(resp.ID)
	}
	return resp, nil
}
