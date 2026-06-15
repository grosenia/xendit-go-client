package xenditgo

// XenditCreateInvoiceReq is JSON request to be sent to Xendit when an Invoice Created
type XenditCreateInvoiceReq struct {
	ExternalID               string   `json:"external_id"`
	Amount                   float64  `json:"amount"`
	PayerEmail               string   `json:"payer_email"`
	Description              string   `json:"description"`
	ShouldSendEmail          bool     `json:"should_send_email"`
	CallbackVirtualAccountID string   `json:"callback_virtual_account_id,omitempty"`
	InvoiceDuration          int      `json:"invoice_duration"`
	PaymentMethod            []string `json:"payment_methods"`
}

// XenditCreateFixedVaReq  is JSON request to be sent to Xendit to create Callback Fixed VA
type XenditCreateFixedVaReq struct {
	ExternalID           string  `json:"external_id"`
	BankCode             string  `json:"bank_code"`
	Name                 string  `json:"name"`
	VirtualAccountNumber string  `json:"virtual_account_number,omitempty"`
	IsClosed             *bool   `json:"is_closed,omitempty"`
	ExpectedAmount       float64 `json:"expected_amount,omitempty"`
	IsSingleUse          *bool   `json:"is_single_use,omitempty"`
}

type XenditUpdateFixedVaReq struct {
	Name                 string  `json:"name,omitempty"`
	ExpectedAmount       float64 `json:"expected_amount,omitempty"`
	IsClosed             *bool   `json:"is_closed,omitempty"`
	IsSingleUse          *bool   `json:"is_single_use,omitempty"`
	ExpirationDate       string  `json:"expiration_date,omitempty"`
	VirtualAccountNumber string  `json:"virtual_account_number,omitempty"`
}

// XenditCreatePayoutReq is JSON request for Payout feature
type XenditCreatePayoutReq struct {
	ExternalID string  `json:"external_id"`
	Amount     float64 `json:"amount"`
}

type XenditCreateBatchReq struct {
	HeaderID      string             `json:"reference"`
	Disbursements []DisbursementItem `json:"disbursements"`
}

type DisbursementItem struct {
	Amount            float64 `json:"amount"`
	ExternalID        string  `json:"external_id"`
	BankCode          string  `json:"bank_code"`
	BankAccountName   string  `json:"bank_account_name"`
	BankAccountNumber string  `json:"bank_account_number"`
	Description       string  `json:"description"`
}

type XenditCreteQrcodeReq struct {
	ReferenceId string  `json:"reference_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	ExpiredAt   string  `json:"expires_at"`
	Currency    string  `json:"currency"`
}

type XenditCreatePaymentMethodReq struct {
	ReferenceID       string                               `json:"reference_id"`
	Type              string                               `json:"type"`
	Reusability       string                               `json:"reusability"`
	Country           string                               `json:"country"`
	Currency          string                               `json:"currency"`
	ChannelCode       string                               `json:"channel_code"`
	ChannelProperties XenditPaymentMethodChannelProperties `json:"channel_properties"`
	Metadata          map[string]interface{}               `json:"metadata,omitempty"`
	CustomerID        string                               `json:"customer_id,omitempty"`
	Description       string                               `json:"description,omitempty"`
}

type XenditPaymentMethodChannelProperties struct {
	ExpiresAt            string                               `json:"expires_at,omitempty"`
	DisplayName          string                               `json:"display_name,omitempty"`
	VerificationData     *XenditPaymentMethodVerificationData `json:"verification_data,omitempty"`
	VirtualAccountNumber string                               `json:"virtual_account_number,omitempty"`
}

type XenditPaymentMethodVerificationData struct {
	CustomerName           string                                  `json:"customer_name"`
	AcceptedNameVariations []string                                `json:"accepted_name_variations,omitempty"`
	AllowedBankAccounts    []XenditPaymentMethodAllowedBankAccount `json:"allowed_bank_accounts,omitempty"`
}

type XenditPaymentMethodAllowedBankAccount struct {
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

// XenditCreateCreditCardChargeReq is JSON request for credit card charge API.
type XenditCreateCreditCardChargeReq struct {
	TokenID          string  `json:"token_id"`
	ExternalID       string  `json:"external_id"`
	Amount           float64 `json:"amount"`
	AuthenticationID string  `json:"authentication_id,omitempty"`
	Capture          bool    `json:"capture"`
}
