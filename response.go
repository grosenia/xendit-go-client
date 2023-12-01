package xenditgo

import (
	"fmt"
)

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

	XenditErrorResponse
}

// XenditErrorResponse xendit error response
type XenditErrorResponse struct {
	// For Errors
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"message"`
	ErrorStatus  bool   `json:"-"`
}

func (e XenditErrorResponse) Error() string {
	return fmt.Sprintf("[%s] %s", e.ErrorCode, e.ErrorMessage)
}

// XenditCreateFixedVaResp is JSON response returned by Xendit when cal Create Callback Fixed VA
type XenditCreateFixedVaResp struct {
	ID              string `json:"id"`
	OwnerID         string `json:"owner_id"`
	ExternalID      string `json:"external_id"`
	MerchantCode    string `json:"merchant_code"`
	AccountNumber   string `json:"account_number"`
	BankCode        string `json:"bank_code"`
	Name            string `json:"name"`
	IsClosed        bool   `json:"is_closed"`
	ExpirationDate  string `json:"expiration_date"`
	IsSingleUse     bool   `json:"is_single_use"`
	Status          string `json:"status"`
	Currency        string `json:"currency"`
	CreatedDateTime string `json:"created"`
	UpdatedDateTime string `json:"updated"`

	XenditErrorResponse
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

// XenditVerifyPaymentResponse type response when verifying payment notification
type XenditVerifyPaymentResponse struct {
	ID                       string  `json:"id"`
	PaymentID                string  `json:"payment_id"`
	CallbackVirtualAccountID string  `json:"callback_virtual_account_id"`
	ExternalID               string  `json:"external_id"`
	MerchantCode             string  `json:"merchant_code"`
	AccountNumber            string  `json:"account_number"`
	BankCode                 string  `json:"bank_code"`
	Amount                   float64 `json:"amount"`
	TransactionTimestamp     string  `json:"transaction_timestamp"`
}

// XenditCreatePayoutResp response from payout
type XenditCreatePayoutResp struct {
	ID               string  `json:"id"`
	ExternalID       string  `json:"external_id"`
	Amount           float64 `json:"amount"`
	PassCode         string  `json:"passcode"`
	MerchantName     string  `json:"merchant_name"`
	Status           string  `json:"status"`
	ExpiredTimestamp string  `json:"expiration_timestamp"`
	CreatedTimestamp string  `json:"created"`
	PayoutURL        string  `json:"payout_url,omitempty"`
	XenditErrorResponse
}

type XenditCreateBatchResp struct {
	Created             string `json:"created"`
	Reference           string `json:"reference"`
	TotalUpdatedCount   int64  `json:"total_uploaded_count"`
	TotalUploadedAmount int64  `json:"total_uploaded_amount"`
	Status              string `json:"status"`
	ID                  string `json:"id"`
	XenditErrorResponse
}

type XenditQrCodeResp struct {
	Id              string         `json:"id"`
	ReferenceId     string         `json:"reference_id"`
	BusinessId      string         `json:"business_id"`
	Type            string         `json:"type"`
	Currency        string         `json:"currency"`
	Amount          float64        `json:"amount"`
	ChannelCode     string         `json:"channel_code"`
	Status          string         `json:"status"`
	QrString        string         `json:"qr_string"`
	Description     string         `json:"description"`
	ExpiryDate      string         `json:"expires_at"`
	CreatedDateTime string         `json:"created"`
	UpdatedDateTime string         `json:"updated" `
	Basket          []BasketQrCode `json:"basket"`
	PaymentDetails

	XenditErrorResponse
}

type BasketQrCode struct {
	ReferenceId string  `json:"reference_id"`
	name        string  `json:"name"`
	Category    string  `json:"category"`
	Currency    string  `json:"currency"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
	Type        string  `json:"type"`
	Url         string  `json:"url"`
	Description string  `json:"description"`
	SubCategory string  `json:"sub-category"`
}

type PaymentDetails struct {
	ReceiptId     string `json:"reference_id"`
	Source        string `json:"source"`
	Name          string `json:"name"`
	AccountDetail string `json:"account_details"`
}
