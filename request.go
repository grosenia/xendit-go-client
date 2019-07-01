package xenditgo

// XenditCreateInvoiceReq is JSON request to be sent to Xendit when an Invoice Created
type XenditCreateInvoiceReq struct {
	ExternalID               string   `json:"external_id"`
	Amount                   float64  `json:"amount"`
	PayerEmail               string   `json:"payer_email"`
	Description              string   `json:"description"`
	ShouldSendEmail          bool     `json:"should_send_email"`
	CallbackVirtualAccountID string   `json:"callback_virtual_account_id"`
	InvoiceDuration          int      `json:"invoice_duration"`
	PaymentMethod            []string `json:"payment_methods"`
}