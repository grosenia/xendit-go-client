package xenditv3

// CreateReusableVirtualAccount creates a fixed/reusable VA (replaces legacy POST /callback_virtual_accounts).
func (g *Gateway) CreateReusableVirtualAccount(in *ReusableVARequest) (*PaymentRequestResponse, error) {
	req, err := BuildReusableVARequest(*in)
	if err != nil {
		return nil, err
	}
	return g.CreatePaymentRequest(req)
}
