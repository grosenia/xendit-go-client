package xenditv3

// PaymentCaptureCallback is webhook payload for payment.capture.
type PaymentCaptureCallback struct {
	Event      string `json:"event"`
	BusinessID string `json:"business_id"`
	Created    string `json:"created"`
	Data       struct {
		PaymentID        string  `json:"payment_id"`
		PaymentRequestID string  `json:"payment_request_id"`
		ReferenceID      string  `json:"reference_id"`
		Status           string  `json:"status"`
		Currency         string  `json:"currency"`
		RequestAmount    float64 `json:"request_amount"`
		CaptureAmount    float64 `json:"capture_amount"`
		PaymentTokenID   string  `json:"payment_token_id"`
	} `json:"data"`
}

func (p PaymentCaptureCallback) PaymentRequestID() string { return p.Data.PaymentRequestID }
func (p PaymentCaptureCallback) ReferenceID() string     { return p.Data.ReferenceID }
func (p PaymentCaptureCallback) PaymentStatus() string     { return p.Data.Status }

func (p PaymentCaptureCallback) IsPaymentSucceeded() bool {
	switch p.Data.Status {
	case "SUCCEEDED", "CAPTURED", "SUCCESS", "PAID":
		return true
	default:
		return false
	}
}

// PaymentTokenCallback is webhook payload for payment_token.activated.
type PaymentTokenCallback struct {
	Event      string `json:"event"`
	BusinessID string `json:"business_id"`
	Created    string `json:"created"`
	Data       struct {
		PaymentTokenID string                   `json:"payment_token_id"`
		ReferenceID    string                   `json:"reference_id"`
		Status         string                   `json:"status"`
		CustomerID     string                   `json:"customer_id"`
		ChannelCode    string                   `json:"channel_code"`
		TokenDetails   *PaymentTokenCardDetails `json:"token_details"`
		CardDetails    *PaymentTokenCardDetails `json:"card_details"`
	} `json:"data"`
}

func (p PaymentTokenCallback) PaymentTokenID() string { return p.Data.PaymentTokenID }

func (p PaymentTokenCallback) MaskedCardNumber() string {
	if p.Data.TokenDetails != nil && p.Data.TokenDetails.MaskedCardNumber != "" {
		return p.Data.TokenDetails.MaskedCardNumber
	}
	if p.Data.CardDetails != nil {
		return p.Data.CardDetails.MaskedCardNumber
	}
	return ""
}

// PaymentSessionCompletedCallback is webhook payload for payment_session.completed.
type PaymentSessionCompletedCallback struct {
	Event string `json:"event"`
	Data  struct {
		PaymentSessionID string `json:"payment_session_id"`
		PaymentRequestID string `json:"payment_request_id"`
		ReferenceID      string `json:"reference_id"`
		Status           string `json:"status"`
		PaymentTokenID   string `json:"payment_token_id"`
	} `json:"data"`
}

func (p PaymentSessionCompletedCallback) PaymentSessionID() string { return p.Data.PaymentSessionID }
func (p PaymentSessionCompletedCallback) ReferenceID() string      { return p.Data.ReferenceID }
func (p PaymentSessionCompletedCallback) PaymentTokenID() string   { return p.Data.PaymentTokenID }

func (p PaymentSessionCompletedCallback) IsCompleted() bool {
	switch p.Data.Status {
	case "COMPLETED", "SUCCEEDED", "SUCCESS":
		return true
	default:
		return false
	}
}
