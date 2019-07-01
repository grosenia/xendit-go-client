package xenditgo

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
