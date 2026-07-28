package xenditv3

import "net/url"

func baseURL(client Client) string {
	if client.BaseURL != "" {
		return client.BaseURL
	}
	return DefaultBaseURL
}

func paymentSessionURL(client Client) string {
	return baseURL(client) + "/sessions"
}

func getPaymentSessionURL(client Client, paymentSessionID string) string {
	return baseURL(client) + "/sessions/" + paymentSessionID
}

func cancelPaymentSessionURL(client Client, paymentSessionID string) string {
	return baseURL(client) + "/sessions/" + paymentSessionID + "/cancel"
}

func createPaymentRequestURL(client Client) string {
	return baseURL(client) + "/v3/payment_requests"
}

func getPaymentRequestURL(client Client, paymentRequestID string) string {
	return baseURL(client) + "/v3/payment_requests/" + paymentRequestID
}

func updatePaymentRequestURL(client Client, paymentRequestID string) string {
	return baseURL(client) + "/v3/payment_requests/" + paymentRequestID
}

func cancelPaymentRequestURL(client Client, paymentRequestID string) string {
	return baseURL(client) + "/v3/payment_requests/" + paymentRequestID + "/cancel"
}

func getPaymentTokenURL(client Client, paymentTokenID string) string {
	return baseURL(client) + "/v3/payment_tokens/" + paymentTokenID
}

func cancelPaymentTokenURL(client Client, paymentTokenID string) string {
	return baseURL(client) + "/v3/payment_tokens/" + paymentTokenID + "/cancel"
}

func getCustomersByReferenceIDURL(client Client, referenceID string) string {
	return baseURL(client) + "/customers?reference_id=" + url.QueryEscape(referenceID)
}

func createCustomerURL(client Client) string {
	return baseURL(client) + "/customers"
}

func createPayoutURL(client Client) string {
	return baseURL(client) + "/v3/payouts"
}

func getPayoutURL(client Client, payoutID string) string {
	return baseURL(client) + "/v3/payouts/" + payoutID
}
