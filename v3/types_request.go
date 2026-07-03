package xenditv3

// CustomerCreateRequest is JSON for POST /customers (given_names at root, not nested).
type CustomerCreateRequest struct {
	ReferenceID  string `json:"reference_id"`
	Type         string `json:"type"`
	Email        string `json:"email,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
	GivenNames   string `json:"given_names"`
	Surname      string `json:"surname,omitempty"`
}

// SessionCustomer is customer payload embedded in POST /sessions.
type SessionCustomer struct {
	ReferenceID      string                   `json:"reference_id"`
	Type             string                   `json:"type"`
	Email            string                   `json:"email,omitempty"`
	MobileNumber     string                   `json:"mobile_number,omitempty"`
	IndividualDetail *SessionIndividualDetail `json:"individual_detail,omitempty"`
}

// SessionIndividualDetail is individual customer detail.
type SessionIndividualDetail struct {
	GivenNames string `json:"given_names"`
	Surname    string `json:"surname"`
}

// SessionCardsChannelProperties is cards-specific channel properties for sessions.
type SessionCardsChannelProperties struct {
	CardOnFileType      string `json:"card_on_file_type,omitempty"`
	TransactionSequence string `json:"transaction_sequence,omitempty"`
	SkipThreeDS         *bool  `json:"skip_three_ds,omitempty"`
}

// SessionChannelProperties groups channel properties for sessions.
type SessionChannelProperties struct {
	Cards               *SessionCardsChannelProperties `json:"cards,omitempty"`
	TransactionSequence string                         `json:"transaction_sequence,omitempty"`
	CardOnFileType      string                         `json:"card_on_file_type,omitempty"`
}

// CardsSessionJS configures CARDS_SESSION_JS mode.
type CardsSessionJS struct {
	SuccessReturnURL string `json:"success_return_url"`
	FailureReturnURL string `json:"failure_return_url"`
}

// PaymentLinkSession configures PAYMENT_LINK mode (pengganti /v2/invoices).
type PaymentLinkSession struct {
	SuccessReturnURL string `json:"success_return_url,omitempty"`
	FailureReturnURL string `json:"failure_return_url,omitempty"`
}

// ReusablePaymentCodeChannelProperties — fixed VA / static QRIS.
type ReusablePaymentCodeChannelProperties struct {
	ExpiresAt   string `json:"expires_at,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// CreatePaymentSessionRequest is JSON for POST /sessions.
type CreatePaymentSessionRequest struct {
	ReferenceID            string                    `json:"reference_id"`
	SessionType            string                    `json:"session_type"`
	Mode                   string                    `json:"mode"`
	Amount                 float64                   `json:"amount"`
	Currency               string                    `json:"currency"`
	Country                string                    `json:"country"`
	CustomerID             string                    `json:"customer_id,omitempty"`
	Customer               *SessionCustomer          `json:"customer,omitempty"`
	PaymentTokenID         string                    `json:"payment_token_id,omitempty"`
	ChannelProperties      *SessionChannelProperties `json:"channel_properties,omitempty"`
	CardsSessionJS         *CardsSessionJS           `json:"cards_session_js,omitempty"`
	PaymentLink            *PaymentLinkSession       `json:"payment_link,omitempty"`
	AllowedPaymentChannels []string                  `json:"allowed_payment_channels,omitempty"`
	ExpiresAt              string                    `json:"expires_at,omitempty"`
	AllowSavePaymentMethod string                    `json:"allow_save_payment_method,omitempty"`
	CaptureMethod          string                    `json:"capture_method,omitempty"`
	SuccessReturnURL       string                    `json:"success_return_url,omitempty"`
	CancelReturnURL        string                    `json:"cancel_return_url,omitempty"`
	Description            string                    `json:"description,omitempty"`
	Metadata               map[string]interface{}    `json:"metadata,omitempty"`
}

// PaymentRequestCardDetails is card input for pay-with-token requests.
type PaymentRequestCardDetails struct {
	Cvn string `json:"cvn,omitempty"`
}

// PaymentRequestChannelProperties is channel properties for POST /v3/payment_requests.
type PaymentRequestChannelProperties struct {
	SkipThreeDS         *bool                      `json:"skip_three_ds,omitempty"`
	CardOnFileType      string                     `json:"card_on_file_type,omitempty"`
	TransactionSequence string                     `json:"transaction_sequence,omitempty"`
	SuccessReturnURL    string                     `json:"success_return_url,omitempty"`
	FailureReturnURL    string                     `json:"failure_return_url,omitempty"`
	StatementDescriptor string                     `json:"statement_descriptor,omitempty"`
	CardDetails         *PaymentRequestCardDetails `json:"card_details,omitempty"`
	ExpiresAt           string                     `json:"expires_at,omitempty"`
	DisplayName         string                     `json:"display_name,omitempty"`
	// VirtualAccountNumber — custom VA number request (channel harus support "Custom Payment
	// Code", lihat https://docs.xendit.co/xenpayments/virtual-account/). Xendit menolak dengan
	// INVALID_VALUE_ERROR "outside of allowed range" kalau prefix/rentang belum di-approve untuk
	// endpoint ini (dikonfirmasi live, terpisah dari approval range di endpoint legacy).
	VirtualAccountNumber string                                `json:"virtual_account_number,omitempty"`
	ReusablePaymentCode  *ReusablePaymentCodeChannelProperties `json:"reusable_payment_code,omitempty"`
	// VerificationData — wajib untuk BRI_VIRTUAL_ACCOUNT sejak Aug 2025 (verifikasi nama/rekening
	// pengirim). Referensi: https://docs.xendit.co/docs/bri-virtual-account
	VerificationData *PaymentRequestVerificationData `json:"verification_data,omitempty"`
}

// PaymentRequestVerificationData — verifikasi identitas pengirim untuk channel VA yang
// mewajibkannya (mis. BRI_VIRTUAL_ACCOUNT sejak Aug 2025).
type PaymentRequestVerificationData struct {
	CustomerName           string                      `json:"customer_name"`
	AcceptedNameVariations []string                    `json:"accepted_name_variations,omitempty"`
	AllowedBankAccounts    []PaymentRequestBankAccount `json:"allowed_bank_accounts,omitempty"`
}

// PaymentRequestBankAccount — satu entri allowed_bank_accounts di VerificationData.
type PaymentRequestBankAccount struct {
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

// CreatePaymentRequestRequest is JSON for POST /v3/payment_requests.
type CreatePaymentRequestRequest struct {
	ReferenceID       string                           `json:"reference_id"`
	PaymentTokenID    string                           `json:"payment_token_id,omitempty"`
	Type              string                           `json:"type"`
	Country           string                           `json:"country"`
	Currency          string                           `json:"currency"`
	RequestAmount     float64                          `json:"request_amount"`
	ChannelCode       string                           `json:"channel_code,omitempty"`
	CaptureMethod     string                           `json:"capture_method,omitempty"`
	ChannelProperties *PaymentRequestChannelProperties `json:"channel_properties,omitempty"`
	Description       string                           `json:"description,omitempty"`
	Metadata          map[string]interface{}           `json:"metadata,omitempty"`
}

// UpdatePaymentRequestRequest is JSON for PATCH /v3/payment_requests/{id}. Hanya berlaku saat
// status REQUIRES_ACTION/ACCEPTING_PAYMENTS. Dikonfirmasi lewat pengujian: request_amount
// benar-benar ter-update tanpa mengubah nomor VA; field lain di channel_properties (mis.
// display_name) belum terkonfirmasi ikut ter-update lewat endpoint ini.
type UpdatePaymentRequestRequest struct {
	RequestAmount float64 `json:"request_amount,omitempty"`
}
