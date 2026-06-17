package xenditv3

import "strings"

// CustomerResponse is a customer object from GET/POST /customers.
type CustomerResponse struct {
	ID          string `json:"id"`
	ReferenceID string `json:"reference_id"`
	Email       string `json:"email"`
	ErrorResponse
}

type customerListResponse struct {
	Data    []CustomerResponse `json:"data"`
	HasMore bool               `json:"has_more"`
}

// PaymentSessionResponse is JSON from session APIs.
type PaymentSessionResponse struct {
	PaymentSessionID       string                    `json:"payment_session_id"`
	ReferenceID            string                    `json:"reference_id"`
	CustomerID             string                    `json:"customer_id"`
	SessionType            string                    `json:"session_type"`
	Mode                   string                    `json:"mode"`
	Currency               string                    `json:"currency"`
	Country                string                    `json:"country"`
	Amount                 interface{}               `json:"amount"`
	Status                 string                    `json:"status"`
	PaymentLinkURL         string                    `json:"payment_link_url"`
	PaymentTokenID         string                    `json:"payment_token_id"`
	PaymentRequestID       string                    `json:"payment_request_id"`
	BusinessID             string                    `json:"business_id"`
	AllowSavePaymentMethod string                    `json:"allow_save_payment_method"`
	CaptureMethod          string                    `json:"capture_method"`
	ChannelProperties      *SessionChannelProperties `json:"channel_properties"`
	CardsSessionJS         *CardsSessionJS           `json:"cards_session_js"`
	SuccessReturnURL       string                    `json:"success_return_url"`
	CancelReturnURL        string                    `json:"cancel_return_url"`
	ExpiresAt              string                    `json:"expires_at"`
	Created                string                    `json:"created"`
	Updated                string                    `json:"updated"`
	ErrorResponse
}

// PaymentURL returns hosted payment link URL (legacy invoice_url equivalent).
func (r *PaymentSessionResponse) PaymentURL() string {
	if r == nil {
		return ""
	}
	return strings.TrimSpace(r.PaymentLinkURL)
}

// PaymentRequestAction is a redirect or auth action from payment request.
type PaymentRequestAction struct {
	Type       string `json:"type"`
	Descriptor string `json:"descriptor"`
	Value      string `json:"value"`
}

// PaymentRequestResponse is JSON from payment request APIs.
type PaymentRequestResponse struct {
	PaymentRequestID  string                           `json:"payment_request_id"`
	ReferenceID       string                           `json:"reference_id"`
	PaymentTokenID    string                           `json:"payment_token_id"`
	CustomerID        string                           `json:"customer_id"`
	BusinessID        string                           `json:"business_id"`
	Type              string                           `json:"type"`
	Country           string                           `json:"country"`
	Currency          string                           `json:"currency"`
	ChannelCode       string                           `json:"channel_code"`
	RequestAmount     float64                          `json:"request_amount"`
	CaptureMethod     string                           `json:"capture_method"`
	Status            string                           `json:"status"`
	ChannelProperties *PaymentRequestChannelProperties `json:"channel_properties"`
	Actions           []PaymentRequestAction           `json:"actions"`
	Description       string                           `json:"description"`
	Created           string                           `json:"created"`
	Updated           string                           `json:"updated"`
	ErrorResponse
}

// PaymentTokenCardDetails is masked card info from payment token.
type PaymentTokenCardDetails struct {
	MaskedCardNumber string `json:"masked_card_number"`
	CardBrand        string `json:"card_brand"`
	ExpiryMonth      string `json:"expiry_month"`
	ExpiryYear       string `json:"expiry_year"`
}

// PaymentTokenChannelProperties is channel_properties on payment token responses.
type PaymentTokenChannelProperties struct {
	CardDetails *PaymentTokenCardDetails `json:"card_details"`
}

// PaymentTokenResponse is JSON from payment token APIs.
type PaymentTokenResponse struct {
	PaymentTokenID    string                         `json:"payment_token_id"`
	ReferenceID       string                         `json:"reference_id"`
	CustomerID        string                         `json:"customer_id"`
	BusinessID        string                         `json:"business_id"`
	Country           string                         `json:"country"`
	Currency          string                         `json:"currency"`
	ChannelCode       string                         `json:"channel_code"`
	Status            string                         `json:"status"`
	maskedCardNumber  string                         `json:"masked_card_number"`
	TokenDetails      *PaymentTokenCardDetails       `json:"token_details"`
	CardDetails       *PaymentTokenCardDetails       `json:"card_details"`
	ChannelProperties *PaymentTokenChannelProperties `json:"channel_properties"`
	Created           string                         `json:"created"`
	Updated           string                         `json:"updated"`
	ErrorResponse
}

// MaskedCardNumber returns masked PAN from token response.
func (r *PaymentTokenResponse) MaskedCardNumber() string {
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
func (r *PaymentRequestResponse) AuthenticationURL() string {
	for _, action := range r.Actions {
		if action.Type == "REDIRECT_CUSTOMER" && action.Value != "" {
			return action.Value
		}
	}
	return ""
}

// PresentToCustomerValue returns VA/QR presentation payload from actions.
func (r *PaymentRequestResponse) PresentToCustomerValue() string {
	for _, action := range r.Actions {
		if action.Type == "PRESENT_TO_CUSTOMER" && action.Value != "" {
			return action.Value
		}
	}
	return ""
}
