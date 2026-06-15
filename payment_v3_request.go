package xenditgo

// XenditPaymentSessionCustomer is customer payload for POST /sessions.
type XenditPaymentSessionCustomer struct {
	ReferenceID      string                              `json:"reference_id"`
	Type             string                              `json:"type"`
	Email            string                              `json:"email,omitempty"`
	MobileNumber     string                              `json:"mobile_number,omitempty"`
	IndividualDetail *XenditPaymentSessionIndividualDetail `json:"individual_detail,omitempty"`
}

// XenditPaymentSessionIndividualDetail is individual customer detail.
type XenditPaymentSessionIndividualDetail struct {
	GivenNames string `json:"given_names"`
	Surname    string `json:"surname"`
}

// XenditPaymentSessionCardsChannelProperties is cards-specific channel properties.
type XenditPaymentSessionCardsChannelProperties struct {
	CardOnFileType      string `json:"card_on_file_type,omitempty"`
	TransactionSequence string `json:"transaction_sequence,omitempty"`
	SkipThreeDS         *bool  `json:"skip_three_ds,omitempty"`
}

// XenditPaymentSessionChannelProperties groups channel properties for sessions.
type XenditPaymentSessionChannelProperties struct {
	Cards               *XenditPaymentSessionCardsChannelProperties `json:"cards,omitempty"`
	TransactionSequence string                                      `json:"transaction_sequence,omitempty"`
	CardOnFileType      string                                      `json:"card_on_file_type,omitempty"`
}

// XenditPaymentSessionCardsSessionJS configures CARDS_SESSION_JS mode.
type XenditPaymentSessionCardsSessionJS struct {
	SuccessReturnURL string `json:"success_return_url"`
	FailureReturnURL string `json:"failure_return_url"`
}

// XenditCreatePaymentSessionReq is JSON request for POST /sessions.
// Docs: https://docs.xendit.co/apidocs/create-session
type XenditCreatePaymentSessionReq struct {
	ReferenceID              string                                `json:"reference_id"`
	SessionType              string                                `json:"session_type"`
	Mode                     string                                `json:"mode"`
	Amount                   float64                               `json:"amount"`
	Currency                 string                                `json:"currency"`
	Country                  string                                `json:"country"`
	CustomerID               string                                `json:"customer_id,omitempty"`
	Customer                 *XenditPaymentSessionCustomer         `json:"customer,omitempty"`
	PaymentTokenID           string                                `json:"payment_token_id,omitempty"`
	ChannelProperties        *XenditPaymentSessionChannelProperties `json:"channel_properties,omitempty"`
	CardsSessionJS           *XenditPaymentSessionCardsSessionJS   `json:"cards_session_js,omitempty"`
	AllowSavePaymentMethod   string                                `json:"allow_save_payment_method,omitempty"`
	CaptureMethod            string                                `json:"capture_method,omitempty"`
	SuccessReturnURL         string                                `json:"success_return_url,omitempty"`
	CancelReturnURL          string                                `json:"cancel_return_url,omitempty"`
	Description              string                                `json:"description,omitempty"`
	Metadata                 map[string]interface{}                `json:"metadata,omitempty"`
}

// XenditPaymentRequestCardDetails is card input for pay-with-token requests.
type XenditPaymentRequestCardDetails struct {
	Cvn string `json:"cvn,omitempty"`
}

// XenditPaymentRequestChannelProperties is channel properties for POST /v3/payment_requests.
type XenditPaymentRequestChannelProperties struct {
	SkipThreeDS         *bool                            `json:"skip_three_ds,omitempty"`
	CardOnFileType      string                           `json:"card_on_file_type,omitempty"`
	TransactionSequence string                           `json:"transaction_sequence,omitempty"`
	SuccessReturnURL    string                           `json:"success_return_url,omitempty"`
	FailureReturnURL    string                           `json:"failure_return_url,omitempty"`
	StatementDescriptor string                           `json:"statement_descriptor,omitempty"`
	CardDetails         *XenditPaymentRequestCardDetails `json:"card_details,omitempty"`
}

// XenditCreatePaymentRequestReq is JSON request for POST /v3/payment_requests.
// Docs: https://docs.xendit.co/docs/pay-with-tokens
type XenditCreatePaymentRequestReq struct {
	ReferenceID        string                                 `json:"reference_id"`
	PaymentTokenID     string                                 `json:"payment_token_id,omitempty"`
	Type               string                                 `json:"type"`
	Country            string                                 `json:"country"`
	Currency           string                                 `json:"currency"`
	RequestAmount      float64                                `json:"request_amount"`
	CaptureMethod      string                                 `json:"capture_method,omitempty"`
	ChannelProperties  *XenditPaymentRequestChannelProperties `json:"channel_properties,omitempty"`
	Description        string                                 `json:"description,omitempty"`
	Metadata           map[string]interface{}                 `json:"metadata,omitempty"`
}
