package xenditv3

// Payments API v3 constants.
// Docs: https://docs.xendit.co/docs/payments-via-api-overview
const (
	DefaultAPIVersion = "2024-11-11"
	DefaultBaseURL    = "https://api.xendit.co"

	SessionTypePay  = "PAY"
	SessionTypeSave = "SAVE"

	SessionModeCardsSessionJS = "CARDS_SESSION_JS"
	SessionModePaymentLink    = "PAYMENT_LINK"
	SessionModeComponents     = "COMPONENTS"

	PaymentRequestTypePay        = "PAY"
	PaymentRequestTypePayAndSave = "PAY_AND_SAVE"
	// PaymentRequestTypeReusablePaymentCode — fixed VA, static QRIS (multi payment).
	PaymentRequestTypeReusablePaymentCode = "REUSABLE_PAYMENT_CODE"

	// Channel codes — VA & QRIS (Payments API v3). Lihat MIGRATION.md fase 2–4.
	ChannelCodeBCAVirtualAccount     = "BCA_VIRTUAL_ACCOUNT"
	ChannelCodeBRIVirtualAccount     = "BRI_VIRTUAL_ACCOUNT"
	ChannelCodeBNIVirtualAccount     = "BNI_VIRTUAL_ACCOUNT"
	ChannelCodeMandiriVirtualAccount = "MANDIRI_VIRTUAL_ACCOUNT"
	ChannelCodePermataVirtualAccount = "PERMATA_VIRTUAL_ACCOUNT"
	ChannelCodeQRIS                  = "QRIS"

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

// Payout API v3 constants.
// Docs: https://docs.xendit.co/apidocs/create-payout-v3
// Confirmed via live sandbox testing (2026-07-28): domestic Indonesia bank transfers route via
// RoutingTypeSWIFT + the destination bank's SWIFT/BIC code (routing_value_1) — there is no
// bank-code-style routing type for Indonesia in this API.
const (
	PayoutAPIVersion = "2025-09-01"

	RoutingTypeSWIFT = "SWIFT"

	PayoutEntityTypeIndividual = "INDIVIDUAL"
	PayoutEntityTypeBusiness   = "BUSINESS"

	WebhookEventPayoutSucceeded         = "v3_payout.succeeded"
	WebhookEventPayoutFailed            = "v3_payout.failed"
	WebhookEventPayoutReversed          = "v3_payout.reversed"
	WebhookEventPayoutRejected          = "v3_payout.rejected"
	WebhookEventPayoutPendingCompliance = "v3_payout.pending_compliance"
)
