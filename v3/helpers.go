package xenditv3

import (
	"fmt"
	"strings"
	"time"
)

// MapLegacyBankCode converts legacy invoice/fixed-VA bank code to Payments API v3 channel_code.
// Accepts legacy codes (BCA, MANDIRI) or full channel codes (BCA_VIRTUAL_ACCOUNT).
func MapLegacyBankCode(bankCode string) (string, error) {
	code := strings.ToUpper(strings.TrimSpace(bankCode))
	switch code {
	case "BCA", "BCAX":
		return ChannelCodeBCAVirtualAccount, nil
	case "BRI":
		return ChannelCodeBRIVirtualAccount, nil
	case "BNI":
		return ChannelCodeBNIVirtualAccount, nil
	case "MANDIRI":
		return ChannelCodeMandiriVirtualAccount, nil
	case "PERMATA":
		return ChannelCodePermataVirtualAccount, nil
	case ChannelCodeBCAVirtualAccount, ChannelCodeBRIVirtualAccount, ChannelCodeBNIVirtualAccount,
		ChannelCodeMandiriVirtualAccount, ChannelCodePermataVirtualAccount, ChannelCodeQRIS:
		return code, nil
	default:
		return "", fmt.Errorf("unsupported bank_code: %s", bankCode)
	}
}

// MapLegacyPaymentMethods converts legacy invoice payment_methods to allowed_payment_channels.
func MapLegacyPaymentMethods(methods []string) ([]string, error) {
	if len(methods) == 0 {
		return nil, fmt.Errorf("payment_methods is required")
	}
	out := make([]string, 0, len(methods))
	for _, m := range methods {
		ch, err := MapLegacyBankCode(m)
		if err != nil {
			return nil, err
		}
		out = append(out, ch)
	}
	return out, nil
}

// PaymentLinkSessionRequest is input for CreatePaymentLinkSession (replaces legacy invoice).
type PaymentLinkSessionRequest struct {
	ReferenceID    string
	Amount         float64
	Currency       string
	Country        string
	PayerEmail     string
	Description    string
	ExpiresAt      string
	InvoiceDurationSec int
	PaymentMethods []string // legacy: BCA, MANDIRI, …
	CustomerID     string
	Customer       *SessionCustomer
	SuccessReturnURL string
	FailureReturnURL string
}

func (r PaymentLinkSessionRequest) expiresAtRFC3339() string {
	if ts := strings.TrimSpace(r.ExpiresAt); ts != "" {
		return ts
	}
	sec := r.InvoiceDurationSec
	if sec <= 0 {
		sec = 86400
	}
	return time.Now().UTC().Add(time.Duration(sec) * time.Second).Format(time.RFC3339)
}

// BuildPaymentLinkSessionRequest maps legacy invoice fields to POST /sessions PAYMENT_LINK.
func BuildPaymentLinkSessionRequest(in PaymentLinkSessionRequest) (*CreatePaymentSessionRequest, error) {
	channels, err := MapLegacyPaymentMethods(in.PaymentMethods)
	if err != nil {
		return nil, err
	}
	currency := in.Currency
	if currency == "" {
		currency = "IDR"
	}
	country := in.Country
	if country == "" {
		country = "ID"
	}

	req := &CreatePaymentSessionRequest{
		ReferenceID:            in.ReferenceID,
		SessionType:            SessionTypePay,
		Mode:                   SessionModePaymentLink,
		Amount:                 in.Amount,
		Currency:               currency,
		Country:                country,
		CustomerID:             strings.TrimSpace(in.CustomerID),
		Customer:               in.Customer,
		AllowedPaymentChannels: channels,
		ExpiresAt:              in.expiresAtRFC3339(),
		CaptureMethod:          CaptureMethodAutomatic,
		Description:            in.Description,
	}
	if in.SuccessReturnURL != "" || in.FailureReturnURL != "" {
		req.PaymentLink = &PaymentLinkSession{
			SuccessReturnURL: in.SuccessReturnURL,
			FailureReturnURL: in.FailureReturnURL,
		}
	}
	if req.Customer == nil && req.CustomerID == "" && strings.TrimSpace(in.PayerEmail) != "" {
		req.Customer = &SessionCustomer{
			ReferenceID: "buyer-" + in.ReferenceID,
			Type:        "INDIVIDUAL",
			Email:       strings.TrimSpace(in.PayerEmail),
			IndividualDetail: &SessionIndividualDetail{
				GivenNames: "Customer",
				Surname:    "Grosenia",
			},
		}
	}
	return req, nil
}

// ReusableVARequest is input for CreateReusableVirtualAccount (replaces legacy fixed VA).
type ReusableVARequest struct {
	ReferenceID          string
	BankCode             string // legacy BCA / BRI / … or full channel_code
	DisplayName          string
	ExpiresAt            string
	ExpectedAmount       float64
	VirtualAccountNumber string
}

// BuildReusableVARequest maps legacy fixed VA to POST /v3/payment_requests REUSABLE_PAYMENT_CODE.
func BuildReusableVARequest(in ReusableVARequest) (*CreatePaymentRequestRequest, error) {
	channel, err := MapLegacyBankCode(in.BankCode)
	if err != nil {
		return nil, err
	}
	expiresAt := strings.TrimSpace(in.ExpiresAt)
	if expiresAt == "" {
		expiresAt = time.Now().UTC().AddDate(31, 0, 0).Format(time.RFC3339)
	}
	props := &PaymentRequestChannelProperties{
		ExpiresAt:   expiresAt,
		DisplayName: strings.TrimSpace(in.DisplayName),
	}
	if in.VirtualAccountNumber != "" {
		props.ReusablePaymentCode = &ReusablePaymentCodeChannelProperties{
			DisplayName: props.DisplayName,
			ExpiresAt:   expiresAt,
		}
	}

	req := &CreatePaymentRequestRequest{
		ReferenceID:   in.ReferenceID,
		Type:          PaymentRequestTypeReusablePaymentCode,
		Country:       "ID",
		Currency:      "IDR",
		RequestAmount: in.ExpectedAmount,
		ChannelCode:   channel,
		ChannelProperties: props,
	}
	return req, nil
}
