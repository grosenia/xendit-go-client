package xenditgo

import "strings"

// XenditCustomerResp is a customer object from GET/POST /customers.
type XenditCustomerResp struct {
	ID          string `json:"id"`
	ReferenceID string `json:"reference_id"`
	Email       string `json:"email"`

	XenditErrorResponse
}

// XenditCustomerListResp is the list wrapper for GET /customers?reference_id=...
type XenditCustomerListResp struct {
	Data    []XenditCustomerResp `json:"data"`
	HasMore bool                 `json:"has_more"`
}

// XenditPaymentSessionResp is JSON response from session APIs.
type XenditPaymentSessionResp struct {
	PaymentSessionID       string                                 `json:"payment_session_id"`
	ReferenceID            string                                 `json:"reference_id"`
	CustomerID             string                                 `json:"customer_id"`
	SessionType            string                                 `json:"session_type"`
	Mode                   string                                 `json:"mode"`
	Currency               string                                 `json:"currency"`
	Country                string                                 `json:"country"`
	Amount                 interface{}                            `json:"amount"`
	Status                 string                                 `json:"status"`
	PaymentLinkURL         string                                 `json:"payment_link_url"`
	PaymentTokenID         string                                 `json:"payment_token_id"`
	PaymentRequestID       string                                 `json:"payment_request_id"`
	BusinessID             string                                 `json:"business_id"`
	AllowSavePaymentMethod string                                 `json:"allow_save_payment_method"`
	CaptureMethod          string                                 `json:"capture_method"`
	ChannelProperties      *XenditPaymentSessionChannelProperties `json:"channel_properties"`
	CardsSessionJS         *XenditPaymentSessionCardsSessionJS    `json:"cards_session_js"`
	SuccessReturnURL       string                                 `json:"success_return_url"`
	CancelReturnURL        string                                 `json:"cancel_return_url"`
	ExpiresAt              string                                 `json:"expires_at"`
	Created                string                                 `json:"created"`
	Updated                string                                 `json:"updated"`

	XenditErrorResponse
}

// XenditPaymentRequestAction is a redirect or auth action from payment request.
type XenditPaymentRequestAction struct {
	Type       string `json:"type"`
	Descriptor string `json:"descriptor"`
	Value      string `json:"value"`
}

// XenditPaymentRequestResp is JSON response from payment request APIs.
type XenditPaymentRequestResp struct {
	PaymentRequestID  string                                 `json:"payment_request_id"`
	ReferenceID       string                                 `json:"reference_id"`
	PaymentTokenID    string                                 `json:"payment_token_id"`
	CustomerID        string                                 `json:"customer_id"`
	BusinessID        string                                 `json:"business_id"`
	Type              string                                 `json:"type"`
	Country           string                                 `json:"country"`
	Currency          string                                 `json:"currency"`
	ChannelCode       string                                 `json:"channel_code"`
	RequestAmount     float64                                `json:"request_amount"`
	CaptureMethod     string                                 `json:"capture_method"`
	Status            string                                 `json:"status"`
	ChannelProperties *XenditPaymentRequestChannelProperties `json:"channel_properties"`
	Actions           []XenditPaymentRequestAction           `json:"actions"`
	Description       string                                 `json:"description"`
	Created           string                                 `json:"created"`
	Updated           string                                 `json:"updated"`

	XenditErrorResponse
}

// XenditPaymentTokenCardDetails is masked card info from payment token.
type XenditPaymentTokenCardDetails struct {
	MaskedCardNumber string `json:"masked_card_number"`
	CardBrand        string `json:"card_brand"`
	ExpiryMonth      string `json:"expiry_month"`
	ExpiryYear       string `json:"expiry_year"`
}

// XenditPaymentTokenChannelProperties is channel_properties on payment token responses.
type XenditPaymentTokenChannelProperties struct {
	CardDetails *XenditPaymentTokenCardDetails `json:"card_details"`
}

// XenditPaymentTokenResp is JSON response from payment token APIs.
type XenditPaymentTokenResp struct {
	PaymentTokenID    string                               `json:"payment_token_id"`
	ReferenceID       string                               `json:"reference_id"`
	CustomerID        string                               `json:"customer_id"`
	BusinessID        string                               `json:"business_id"`
	Country           string                               `json:"country"`
	Currency          string                               `json:"currency"`
	ChannelCode       string                               `json:"channel_code"`
	Status            string                               `json:"status"`
	maskedCardNumber  string                               `json:"masked_card_number"`
	TokenDetails      *XenditPaymentTokenCardDetails       `json:"token_details"`
	CardDetails       *XenditPaymentTokenCardDetails       `json:"card_details"`
	ChannelProperties *XenditPaymentTokenChannelProperties `json:"channel_properties"`
	Created           string                               `json:"created"`
	Updated           string                               `json:"updated"`

	XenditErrorResponse
}

// MaskedCardNumber returns masked PAN from token response.
func (r *XenditPaymentTokenResp) MaskedCardNumber() string {
	if r == nil {
		return ""
	}
	if masked := strings.TrimSpace(r.maskedCardNumber); masked != "" {
		return masked
	}
	if r.TokenDetails != nil && r.TokenDetails.MaskedCardNumber != "" {
		return r.TokenDetails.MaskedCardNumber
	}
	if r.CardDetails != nil && r.CardDetails.MaskedCardNumber != "" {
		return r.CardDetails.MaskedCardNumber
	}
	if r.ChannelProperties != nil && r.ChannelProperties.CardDetails != nil {
		return r.ChannelProperties.CardDetails.MaskedCardNumber
	}
	return ""
}

// AuthenticationURL returns 3DS redirect URL from payment request actions.
func (r *XenditPaymentRequestResp) AuthenticationURL() string {
	for _, action := range r.Actions {
		if action.Type == "REDIRECT_CUSTOMER" && action.Value != "" {
			return action.Value
		}
	}
	return ""
}
