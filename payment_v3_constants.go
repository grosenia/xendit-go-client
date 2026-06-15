package xenditgo

// Payments API v3 constants.
// Docs: https://docs.xendit.co/docs/payments-via-api-overview
const (
	// PaymentsAPIVersion is the api-version header value for Payments API v3.
	PaymentsAPIVersion = "2024-11-11"

	SessionTypePay  = "PAY"
	SessionTypeSave = "SAVE"

	SessionModeCardsSessionJS = "CARDS_SESSION_JS"
	SessionModePaymentLink    = "PAYMENT_LINK"
	SessionModeComponents     = "COMPONENTS"

	PaymentRequestTypePay        = "PAY"
	PaymentRequestTypePayAndSave = "PAY_AND_SAVE"

	CaptureMethodAutomatic = "AUTOMATIC"
	CaptureMethodManual    = "MANUAL"

	CardOnFileCustomerUnscheduled = "CUSTOMER_UNSCHEDULED"
	CardOnFileRecurring           = "RECURRING"
	CardOnFileMerchantUnscheduled = "MERCHANT_UNSCHEDULED"

	TransactionSequenceSubsequent = "SUBSEQUENT"
	TransactionSequenceInitial    = "INITIAL"

	AllowSavePaymentMethodForced   = "FORCED"
	AllowSavePaymentMethodOptional = "OPTIONAL"
	AllowSavePaymentMethodDisabled = "DISABLED"

	WebhookEventPaymentCapture          = "payment.capture"
	WebhookEventPaymentTokenActivated   = "payment_token.activated"
	WebhookEventPaymentSessionCompleted = "payment_session.completed"
)
