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
	ExternalID           string `json:"external_id"`
	BankCode             string `json:"bank_code"`
	Name                 string `json:"name"`
	VirtualAccountNumber string `json:"virtual_account_number,omitempty"`
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
