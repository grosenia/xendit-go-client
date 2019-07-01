package xenditgo

// XenditCreateInvoiceResp is JSON response returned by Xendit when an Invoice Created
type XenditCreateInvoiceResp struct {
	InvoiceID               string                 `json:"id"`
	UserID                  string                 `json:"user_id"`
	ExternalID              string                 `json:"external_id"`
	IsHigh                  bool                   `json:"is_high"`
	Status                  string                 `json:"status"`
	MerchantName            string                 `json:"merchant_name"`
	Amount                  float64                `json:"amount"`
	PayerEmail              string                 `json:"payer_email"`
	Description             string                 `json:"description"`
	ExpiryDate              string                 `json:"expiry_date"`
	InvoiceURL              string                 `json:"invoice_url"`
	ShouldExcludeCreditCard bool                   `json:"should_exclude_credit_card" `
	ShouldSendEmail         bool                   `json:"should_send_email"`
	CreatedDateTime         string                 `json:"created"`
	UpdatedDateTime         string                 `json:"updated"`
	AvailableBanks          []InvoiceAvailableBank `json:"available_banks"`
}

// InvoiceAvailableBank is options of invoice available bank
type InvoiceAvailableBank struct {
	BankCode          string  `json:"bank_code"`
	CollectionType    string  `json:"collection_type"`
	BankAccountNumber string  `json:"bank_account_number"`
	TransferAmount    float64 `json:"transfer_amount"`
	BankBranch        string  `json:"bank_branch"`
	AccountHolderName string  `json:"account_holder_name"`
	IdentityAmount    float64 `json:"identity_amount"`
}

// XenditNotificationResp is standard response for Xendit
type XenditNotificationResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// XenditPaymentResp is standard response to Payment
type XenditPaymentResp struct {
	PaymentURL string `json:"payment_url"`
}
