package xenditgo

// XenditNotification to store notification payment
type XenditNotification struct {
	TransactionID          string `json:"id"`
	ExternalID             string `json:"external_id"`
	UserID                 string `json:"user_id"`
	IsHigh                 bool   `json:"is_high"`
	PaymentMethod          string `json:"payment_method"`
	Status                 string `json:"status"`
	MerchantName           string `json:"merchant_name"`
	Amount                 string `json:"amount"`
	PaidAmount             string `json:"paid_amount"`
	BankCode               string `json:"bank_code"`
	PayerEmail             string `json:"payer_email"`
	Description            string `json:"description"`
	AdjustedReceivedAmount string `json:"adjusted_received_amount"`
	FeesPaidAmount         string `json:"feeds_paid_amount"`
	CreatedDateTime        string `json:"created_datetime"`
	UpdatedDateTime        string `json:"updated_datetime"`
}

// XenditFixedVaCreatedNotification is standard callback when FixedVaCreated
type XenditFixedVaCreatedNotification struct {
	VaID            string `json:"id"`
	OwnerID         string `json:"owner_id"`
	ExternalID      string `json:"extenral_id"`
	MerchantCode    string `json:"merchant_code"`
	AccountNumber   string `json:"account_number"`
	BankCode        string `json:"bank_code"`
	Name            string `json:"name"`
	IsClosed        bool   `json:"is_closed"`
	ExpirationDate  string `json:"expiration_date"`
	IsSingleUse     bool   `json:"is_single_use"`
	Status          string `json:"active"`
	CreatedDateTime string `json:"created_datetime"`
	UpdatedDateTime string `json:"updated_datetime"`
}

// XenditQrCodeCallback is standard callback when payment
type XenditQrCodeCallback struct {
	Event           string             `json:"event"`
	ApiVersion      string             `json:"api_version"`
	BusinessId      string             `json:"business_id"`
	CreatedDateTime string             `json:"created"`
	Data            []XenditQrCodeResp `json:"data"`
}

// XenditCreditCardCallback is webhook payload for credit card charge events.
type XenditCreditCardCallback struct {
	Event      string `json:"event"`
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	Data       struct {
		ID         string `json:"id"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	} `json:"data"`
}

// ChargeID returns the Xendit charge ID from callback payload.
func (p XenditCreditCardCallback) ChargeID() string {
	if p.Data.ID != "" {
		return p.Data.ID
	}
	return p.ID
}

// OrderNo returns the order reference (external_id) from callback payload.
func (p XenditCreditCardCallback) OrderNo() string {
	if p.Data.ExternalID != "" {
		return p.Data.ExternalID
	}
	return p.ExternalID
}

// ChargeStatus returns normalized charge status from callback payload.
func (p XenditCreditCardCallback) ChargeStatus() string {
	if p.Data.Status != "" {
		return p.Data.Status
	}
	return p.Status
}

// XenditPaymentCaptureCallback is webhook payload for payment.capture (Payments API v3).
type XenditPaymentCaptureCallback struct {
	Event      string `json:"event"`
	BusinessID string `json:"business_id"`
	Created    string `json:"created"`
	Data       struct {
		PaymentID        string `json:"payment_id"`
		PaymentRequestID string `json:"payment_request_id"`
		ReferenceID      string `json:"reference_id"`
		Status           string `json:"status"`
		Currency         string `json:"currency"`
		RequestAmount    float64 `json:"request_amount"`
		CaptureAmount    float64 `json:"capture_amount"`
		PaymentTokenID   string `json:"payment_token_id"`
	} `json:"data"`
}

// PaymentRequestID returns payment_request_id from webhook payload.
func (p XenditPaymentCaptureCallback) PaymentRequestID() string {
	return p.Data.PaymentRequestID
}

// OrderNo returns merchant reference_id from webhook payload.
func (p XenditPaymentCaptureCallback) OrderNo() string {
	return p.Data.ReferenceID
}

// PaymentStatus returns normalized payment status from webhook payload.
func (p XenditPaymentCaptureCallback) PaymentStatus() string {
	return p.Data.Status
}

// IsPaymentSucceeded reports whether capture webhook indicates success.
func (p XenditPaymentCaptureCallback) IsPaymentSucceeded() bool {
	switch p.Data.Status {
	case "SUCCEEDED", "CAPTURED", "SUCCESS", "PAID":
		return true
	default:
		return false
	}
}

// XenditPaymentTokenCallback is webhook payload for payment_token.activated.
type XenditPaymentTokenCallback struct {
	Event      string `json:"event"`
	BusinessID string `json:"business_id"`
	Created    string `json:"created"`
	Data       struct {
		PaymentTokenID string `json:"payment_token_id"`
		ReferenceID    string `json:"reference_id"`
		Status         string `json:"status"`
		CustomerID     string `json:"customer_id"`
		ChannelCode    string `json:"channel_code"`
		TokenDetails   *XenditPaymentTokenCardDetails `json:"token_details"`
		CardDetails    *XenditPaymentTokenCardDetails `json:"card_details"`
	} `json:"data"`
}

// PaymentTokenID returns payment_token_id from webhook payload.
func (p XenditPaymentTokenCallback) PaymentTokenID() string {
	return p.Data.PaymentTokenID
}

// MaskedCardNumber returns masked PAN from token activation webhook.
func (p XenditPaymentTokenCallback) MaskedCardNumber() string {
	if p.Data.TokenDetails != nil && p.Data.TokenDetails.MaskedCardNumber != "" {
		return p.Data.TokenDetails.MaskedCardNumber
	}
	if p.Data.CardDetails != nil {
		return p.Data.CardDetails.MaskedCardNumber
	}
	return ""
}
