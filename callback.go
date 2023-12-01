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
