package xenditv3

// CreatePaymentLinkSession creates a hosted payment page (replaces legacy POST /v2/invoices).
// Docs: https://docs.xendit.co/docs/migrate-to-payment-session
func (g *Gateway) CreatePaymentLinkSession(in *PaymentLinkSessionRequest) (*PaymentSessionResponse, error) {
	req, err := BuildPaymentLinkSessionRequest(*in)
	if err != nil {
		return nil, err
	}
	return g.CreatePaymentSession(req)
}
