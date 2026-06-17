package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	xenditv3 "github.com/grosenia/xendit-go-client/v3"
	"github.com/grosenia/xendit-go-client/v3/moneyout"
	viper "github.com/spf13/viper"
)

func init() {
	viperGetRaw = viper.GetString
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	loadConfig()

	var err error
	switch cmd {
	case "card-session":
		err = runCardSession()
	case "payment-link":
		err = runPaymentLink()
	case "reusable-va":
		err = runReusableVA()
	case "payment-request":
		err = runPaymentRequest()
	case "get-token":
		err = runGetToken()
	case "batch-disburse":
		err = runBatchDisburse()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		logFail(err.Error())
		printTestSummary(cmd, false)
		os.Exit(1)
	}
	printTestSummary(cmd, true)
}

func loadConfig() {
	viper.SetConfigType("props")
	name := strings.TrimSpace(os.Getenv("XENDIT_EXAMPLE_CONFIG"))
	if name == "" {
		name = "config"
	}
	viper.SetConfigName(name)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read %s.props: %w", name, err))
	}
}

func moneyInClient() xenditv3.Client {
	client := xenditv3.NewClient(viper.GetString("KEY_WRITE_MONEY_IN"))
	if v := viper.GetString("API_VERSION"); v != "" {
		client.APIVersion = v
	}
	return client
}

func moneyInGateway() *xenditv3.Gateway {
	return xenditv3.NewGateway(moneyInClient())
}

func refID(prefix string) string {
	if v := strings.TrimSpace(viper.GetString("ORDER_REFERENCE_ID")); v != "" {
		return v
	}
	return fmt.Sprintf("%s-%d", prefix, time.Now().Unix())
}

func resolveCustomerID(gw *xenditv3.Gateway) (string, error) {
	if id := strings.TrimSpace(viper.GetString("CUSTOMER_ID")); id != "" {
		id = xenditv3.NormalizeCustomerID(id)
		logOK("Pakai CUSTOMER_ID dari config: " + id)
		return id, nil
	}

	buyerRef := viper.GetString("BUYER_REFERENCE_ID")
	if buyerRef == "" {
		buyerRef = "buyer-" + strconv.FormatInt(time.Now().Unix(), 10)
	}
	logStep("GET /customers?reference_id=" + buyerRef)
	customer, err := gw.GetCustomerByReferenceID(buyerRef)
	if err != nil {
		return "", fmt.Errorf("get customer: %w", err)
	}
	if customer == nil {
		logStep("POST /customers — buat customer baru")
		customer, err = gw.CreateCustomer(&xenditv3.SessionCustomer{
			ReferenceID:  buyerRef,
			Type:         "INDIVIDUAL",
			Email:        viper.GetString("PAYER_EMAIL"),
			MobileNumber: viper.GetString("PAYER_MOBILE"),
			IndividualDetail: &xenditv3.SessionIndividualDetail{
				GivenNames: viper.GetString("PAYER_GIVEN_NAMES"),
				Surname:    viper.GetString("PAYER_SURNAME"),
			},
		})
		if err != nil {
			return "", fmt.Errorf("create customer: %w", err)
		}
		if customer.ErrorStatus {
			return "", apiErr(customer.Error())
		}
		logOK("Customer dibuat: " + customer.ID)
	} else {
		logOK("Customer reuse: " + customer.ID)
	}
	return customer.ID, nil
}

func runCardSession() error {
	logBanner("card-session", "[1/6] Kartu kredit/debit — Cards Session JS")
	logLegacy("Migrasi", "— (flow baru)", "POST /sessions · mode CARDS_SESSION_JS")
	logConfigLoaded()
	logMoneyInSetup()

	logSection(1, "Siapkan customer")
	gw := moneyInGateway()
	customerID, err := resolveCustomerID(gw)
	if err != nil {
		return err
	}

	amount := viper.GetFloat64("AMOUNT")
	if amount <= 0 {
		amount = 50000
	}

	ref := refID("ord")
	logSection(2, "Buat payment session (kartu baru)")
	logStep("POST /sessions")
	logInfo("reference_id", ref)
	logInfo("amount", fmt.Sprintf("%.0f IDR", amount))
	logInfo("mode", xenditv3.SessionModeCardsSessionJS)

	session, err := gw.CreatePaymentSession(&xenditv3.CreatePaymentSessionRequest{
		ReferenceID:            ref,
		SessionType:            xenditv3.SessionTypePay,
		Mode:                   xenditv3.SessionModeCardsSessionJS,
		Amount:                 amount,
		Currency:               "IDR",
		Country:                "ID",
		CustomerID:             customerID,
		CaptureMethod:          xenditv3.CaptureMethodAutomatic,
		AllowSavePaymentMethod: xenditv3.AllowSavePaymentMethodOptional,
		CardsSessionJS: &xenditv3.CardsSessionJS{
			SuccessReturnURL: viper.GetString("SUCCESS_RETURN_URL"),
			FailureReturnURL: viper.GetString("FAILURE_RETURN_URL"),
		},
	})
	if err != nil {
		return err
	}
	if session.ErrorStatus {
		return apiErr(session.Error())
	}

	logSection(3, "Hasil — lanjut di frontend")
	logOK("Session ACTIVE")
	logInfo("payment_session_id", session.PaymentSessionID)
	logInfo("status", session.Status)
	logInfo("customer_id", session.CustomerID)
	logInfo("PUBLIC_API_KEY (frontend)", maskKey(viper.GetString("PUBLIC_API_KEY")))
	fmt.Println()
	fmt.Println("  Langkah berikutnya:")
	fmt.Println("  1. Pakai cards-session.min.js + payment_session_id di frontend")
	fmt.Println("  2. Setelah save card → isi PAYMENT_TOKEN_ID di config.props")
	fmt.Println("  3. Test: go run . get-token  dan  go run . payment-request")
	return nil
}

func runPaymentLink() error {
	logBanner("payment-link", "[2/6] Payment Link — ganti legacy invoice")
	logLegacy("Migrasi", "POST /v2/invoices", "POST /sessions · mode PAYMENT_LINK")
	logConfigLoaded()
	logMoneyInSetup()

	methods := splitCSV(viper.GetString("PAYMENT_METHODS"))
	if len(methods) == 0 {
		methods = []string{"BCA", "BNI", "MANDIRI", "BRI", "PERMATA"}
	}

	logSection(1, "Buat payment link session")
	logStep("POST /sessions (CreatePaymentLinkSession)")
	logInfo("reference_id", refID("inv"))
	logInfo("amount", fmt.Sprintf("%.0f IDR", viper.GetFloat64("AMOUNT")))
	logInfo("payment_methods", strings.Join(methods, ", "))
	logInfo("invoice_duration", fmt.Sprintf("%d detik", viper.GetInt("INVOICE_DURATION")))

	session, err := moneyInGateway().CreatePaymentLinkSession(&xenditv3.PaymentLinkSessionRequest{
		ReferenceID:        refID("inv"),
		Amount:             viper.GetFloat64("AMOUNT"),
		PayerEmail:         viper.GetString("PAYER_EMAIL"),
		Description:        viper.GetString("DESCRIPTION"),
		InvoiceDurationSec: viper.GetInt("INVOICE_DURATION"),
		PaymentMethods:     methods,
		SuccessReturnURL:   viper.GetString("SUCCESS_RETURN_URL"),
		FailureReturnURL:   viper.GetString("FAILURE_RETURN_URL"),
	})
	if err != nil {
		return err
	}
	if session.ErrorStatus {
		return apiErr(session.Error())
	}

	logSection(2, "Hasil — kirim link ke buyer")
	logOK("Payment link dibuat")
	logInfo("payment_url (= invoice_url)", session.PaymentURL())
	logInfo("payment_session_id", session.PaymentSessionID)
	logInfo("status", session.Status)
	logInfo("expires_at", session.ExpiresAt)
	fmt.Println()
	fmt.Println("  Webhook setelah bayar: payment.capture (bukan invoice callback legacy)")
	return nil
}

func runReusableVA() error {
	logBanner("reusable-va", "[3/6] Fixed VA seller — ganti legacy fixed VA")
	logLegacy("Migrasi", "POST /callback_virtual_accounts", "POST /v3/payment_requests · REUSABLE_PAYMENT_CODE")
	logConfigLoaded()
	logMoneyInSetup()

	bank := viper.GetString("BANK_CODE")
	channel, err := xenditv3.MapLegacyBankCode(bank)
	if err != nil {
		return err
	}

	logSection(1, "Buat reusable VA")
	logStep("POST /v3/payment_requests (CreateReusableVirtualAccount)")
	logInfo("reference_id", refID("va"))
	logInfo("bank_code (config)", bank)
	logInfo("channel_code (v3)", channel)
	logInfo("display_name", viper.GetString("DISPLAY_NAME"))

	pr, err := moneyInGateway().CreateReusableVirtualAccount(&xenditv3.ReusableVARequest{
		ReferenceID:    refID("va"),
		BankCode:       bank,
		DisplayName:    viper.GetString("DISPLAY_NAME"),
		ExpiresAt:      viper.GetString("EXPIRES_AT"),
		ExpectedAmount: viper.GetFloat64("EXPECTED_AMOUNT"),
	})
	if err != nil {
		return err
	}
	if pr.ErrorStatus {
		return apiErr(pr.Error())
	}

	logSection(2, "Hasil")
	logOK("Reusable VA dibuat")
	logInfo("payment_request_id", pr.PaymentRequestID)
	logInfo("channel_code", pr.ChannelCode)
	logInfo("status", pr.Status)
	if v := pr.PresentToCustomerValue(); v != "" {
		logInfo("present_to_customer", v)
	}
	fmt.Println()
	fmt.Println("  Webhook setelah bayar: payment.capture")
	return nil
}

func runPaymentRequest() error {
	logBanner("payment-request", "[4/6] Bayar pakai kartu tersimpan + CVN")
	logLegacy("Migrasi", "—", "POST /v3/payment_requests · type PAY + payment_token_id")
	logConfigLoaded()
	logMoneyInSetup()

	tokenID := strings.TrimSpace(viper.GetString("PAYMENT_TOKEN_ID"))
	cvn := strings.TrimSpace(viper.GetString("CVN"))
	if tokenID == "" || cvn == "" || strings.Contains(tokenID, "REPLACE") {
		return fmt.Errorf("isi PAYMENT_TOKEN_ID dan CVN di config.props (dapat setelah card-session + save card)")
	}

	amount := viper.GetFloat64("AMOUNT")
	if amount <= 0 {
		amount = 50000
	}

	logSection(1, "Buat payment request (saved card)")
	logStep("POST /v3/payment_requests")
	logInfo("reference_id", refID("ord"))
	logInfo("payment_token_id", tokenID)
	logInfo("amount", fmt.Sprintf("%.0f IDR", amount))
	logInfo("cvn", "*** (len="+fmt.Sprint(len(cvn))+")")

	pr, err := moneyInGateway().CreatePaymentRequest(&xenditv3.CreatePaymentRequestRequest{
		ReferenceID:    refID("ord"),
		PaymentTokenID: tokenID,
		Type:           xenditv3.PaymentRequestTypePay,
		Country:        "ID",
		Currency:       "IDR",
		RequestAmount:  amount,
		CaptureMethod:  xenditv3.CaptureMethodAutomatic,
		ChannelProperties: &xenditv3.PaymentRequestChannelProperties{
			CardOnFileType:      xenditv3.CardOnFileCustomerUnscheduled,
			TransactionSequence: xenditv3.TransactionSequenceSubsequent,
			CardDetails:         &xenditv3.PaymentRequestCardDetails{Cvn: cvn},
		},
	})
	if err != nil {
		return err
	}
	if pr.ErrorStatus {
		return apiErr(pr.Error())
	}

	logSection(2, "Hasil")
	logOK("Payment request dibuat")
	logInfo("payment_request_id", pr.PaymentRequestID)
	logInfo("status", pr.Status)
	if u := pr.AuthenticationURL(); u != "" {
		logInfo("3ds_url", u)
		fmt.Println("  → Buyer perlu selesaikan 3DS di URL di atas")
	}
	return nil
}

func runGetToken() error {
	logBanner("get-token", "[5/6] Cek payment token tersimpan")
	logLegacy("Migrasi", "—", "GET /v3/payment_tokens/{id}")
	logConfigLoaded()
	logMoneyInSetup()

	tokenID := strings.TrimSpace(viper.GetString("PAYMENT_TOKEN_ID"))
	if tokenID == "" || strings.Contains(tokenID, "REPLACE") {
		return fmt.Errorf("isi PAYMENT_TOKEN_ID di config.props")
	}

	logSection(1, "Ambil detail token")
	logStep("GET /v3/payment_tokens/" + tokenID)

	token, err := moneyInGateway().GetPaymentToken(tokenID)
	if err != nil {
		return err
	}
	if token.ErrorStatus {
		return apiErr(token.Error())
	}

	logSection(2, "Hasil")
	logOK("Token ditemukan")
	logInfo("payment_token_id", token.PaymentTokenID)
	logInfo("status", token.Status)
	logInfo("masked_card", token.MaskedCardNumber())
	logInfo("customer_id", token.CustomerID)
	return nil
}

func runBatchDisburse() error {
	logBanner("batch-disburse", "[6/6] Batch disbursement — Money Out")
	logLegacy("Migrasi", "POST /batch_disbursements (root)", "v3/moneyout · API sama")
	logConfigLoaded()
	logMoneyOutSetup()

	headerID := fmt.Sprintf("v3-batch-%d", time.Now().Unix())

	logSection(1, "Upload batch disbursement")
	logStep("POST /batch_disbursements")
	logInfo("reference", headerID)
	logInfo("idempotency_key", headerID)
	logInfo("bank_code", viper.GetString("BANK_CODE"))
	logInfo("amount", fmt.Sprintf("%.0f IDR", viper.GetFloat64("DISBURSEMENT_AMOUNT")))
	logInfo("account", viper.GetString("BANK_ACCOUNT_NAME")+" / "+viper.GetString("BANK_ACCOUNT_NUMBER"))

	client := moneyout.NewClient(viper.GetString("KEY_WRITE_MONEY_OUT"))
	resp, err := moneyout.NewGateway(client).CreateBatchDisbursement(headerID, &moneyout.CreateBatchDisbursementRequest{
		Reference: headerID,
		Disbursements: []moneyout.DisbursementItem{{
			ExternalID:        headerID + "-1",
			Amount:            viper.GetFloat64("DISBURSEMENT_AMOUNT"),
			BankCode:          viper.GetString("BANK_CODE"),
			BankAccountName:   viper.GetString("BANK_ACCOUNT_NAME"),
			BankAccountNumber: viper.GetString("BANK_ACCOUNT_NUMBER"),
			Description:       viper.GetString("DISBURSEMENT_DESCRIPTION"),
		}},
	})
	if err != nil {
		return err
	}
	if resp.ErrorStatus {
		return apiErr(resp.Error())
	}

	logSection(2, "Hasil")
	logOK("Batch terupload")
	logInfo("id", resp.ID)
	logInfo("reference", resp.Reference)
	logInfo("status", resp.Status)
	return nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
