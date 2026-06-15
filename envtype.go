package xenditgo

import "net/url"

// EnvironmentType value
type EnvironmentType int8

const (
	_ EnvironmentType = iota

	// Sandbox : represent sandbox environment
	Sandbox

	// Production : represent production environment
	Production
)

var typeString = map[EnvironmentType]string{
	Sandbox:    "https://api.xendit.co",
	Production: "https://api.xendit.co",
}

// implement stringer
func (e EnvironmentType) String() string {
	for k, v := range typeString {
		if k == e {
			return v
		}
	}
	return "undefined"
}

// CreateInvoiceURL : Create invoice for accepting payment
func (e EnvironmentType) CreateInvoiceURL() string {
	return e.String() + "/v2/invoices"
}

// CreateDisbursementURL : Create a disbursement
func (e EnvironmentType) CreateDisbursementURL() string {
	return e.String() + "/disbursements"
}

// GetVirtualAccountBanksURL : Get available virtual account banks
func (e EnvironmentType) GetVirtualAccountBanksURL() string {
	return e.String() + "/available_virtual_account_banks"
}

// CreateCallbackVirtualAccountURL : is used to create FixedVA
func (e EnvironmentType) CreateCallbackVirtualAccountURL() string {
	return e.String() + "/callback_virtual_accounts"
}

// CreateCreditCardChargeURL creates a credit card charge (legacy v2).
func (e EnvironmentType) CreateCreditCardChargeURL() string {
	return e.String() + "/credit_card_charges"
}

// CreatePaymentSessionURL creates a Payments API v3 session.
func (e EnvironmentType) CreatePaymentSessionURL() string {
	return e.String() + "/sessions"
}

// GetCustomersByReferenceIDURL looks up an existing Xendit customer by merchant reference_id.
func (e EnvironmentType) GetCustomersByReferenceIDURL(referenceID string) string {
	return e.String() + "/customers?reference_id=" + url.QueryEscape(referenceID)
}

// CreateCustomerURL creates an end-customer resource.
func (e EnvironmentType) CreateCustomerURL() string {
	return e.String() + "/customers"
}

// GetPaymentSessionURL retrieves a Payments API v3 session.
func (e EnvironmentType) GetPaymentSessionURL(paymentSessionID string) string {
	return e.String() + "/sessions/" + paymentSessionID
}

// CancelPaymentSessionURL cancels a Payments API v3 session.
func (e EnvironmentType) CancelPaymentSessionURL(paymentSessionID string) string {
	return e.String() + "/sessions/" + paymentSessionID + "/cancel"
}

// CreatePaymentRequestURL creates a Payments API v3 payment request.
func (e EnvironmentType) CreatePaymentRequestURL() string {
	return e.String() + "/v3/payment_requests"
}

// GetPaymentRequestURL retrieves a Payments API v3 payment request.
func (e EnvironmentType) GetPaymentRequestURL(paymentRequestID string) string {
	return e.String() + "/v3/payment_requests/" + paymentRequestID
}

// GetPaymentTokenURL retrieves a Payments API v3 payment token.
func (e EnvironmentType) GetPaymentTokenURL(paymentTokenID string) string {
	return e.String() + "/v3/payment_tokens/" + paymentTokenID
}

// CancelPaymentTokenURL deactivates a Payments API v3 payment token.
func (e EnvironmentType) CancelPaymentTokenURL(paymentTokenID string) string {
	return e.String() + "/v3/payment_tokens/" + paymentTokenID + "/cancel"
}
