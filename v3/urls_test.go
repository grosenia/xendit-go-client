package xenditv3

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestURLs(t *testing.T) {
	is := is.New(t)
	client := Client{BaseURL: DefaultBaseURL}

	is.Equal("https://api.xendit.co/sessions", paymentSessionURL(client))
	is.Equal("https://api.xendit.co/sessions/ps-test", getPaymentSessionURL(client, "ps-test"))
	is.Equal("https://api.xendit.co/sessions/ps-test/cancel", cancelPaymentSessionURL(client, "ps-test"))
	is.Equal("https://api.xendit.co/v3/payment_requests", createPaymentRequestURL(client))
	is.Equal("https://api.xendit.co/v3/payment_requests/pr-test", getPaymentRequestURL(client, "pr-test"))
	is.Equal("https://api.xendit.co/v3/payment_tokens/pt-test", getPaymentTokenURL(client, "pt-test"))
	is.Equal("https://api.xendit.co/v3/payment_tokens/pt-test/cancel", cancelPaymentTokenURL(client, "pt-test"))
}
