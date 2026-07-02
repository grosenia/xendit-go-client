package xenditv3

// CreatePoolVirtualAccount creates a one-time pool VA (replaces legacy POST /v2/invoices POOL).
// Response includes PRESENT_TO_CUSTOMER with virtual account number for native UI.
func (g *Gateway) CreatePoolVirtualAccount(in *PoolVARequest) (*PaymentRequestResponse, error) {
	req, err := BuildPoolVARequest(*in)
	if err != nil {
		return nil, err
	}
	return g.CreatePaymentRequest(req)
}
