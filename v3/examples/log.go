package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func logBanner(title, subtitle string) {
	line := strings.Repeat("=", 60)
	fmt.Println()
	fmt.Println(line)
	fmt.Printf("  TEST: %s\n", title)
	if subtitle != "" {
		fmt.Printf("  %s\n", subtitle)
	}
	fmt.Println(line)
}

func logSection(n int, title string) {
	fmt.Printf("\n[%d] %s\n", n, title)
	fmt.Println(strings.Repeat("-", 40))
}

func logInfo(label, value string) {
	fmt.Printf("  %-22s %s\n", label+":", value)
}

func logStep(msg string) {
	fmt.Printf("  → %s\n", msg)
}

func logOK(msg string) {
	fmt.Printf("  ✓ %s\n", msg)
}

func logFail(msg string) {
	fmt.Fprintf(os.Stderr, "  ✗ %s\n", msg)
}

func logLegacy(label, legacyAPI, v3API string) {
	fmt.Printf("  %-22s legacy: %s\n", label+":", legacyAPI)
	fmt.Printf("  %-22s v3:     %s\n", "", v3API)
}

func apiErr(msg string) error {
	return fmt.Errorf("%s", msg)
}

func maskKey(key string) string {
	key = strings.TrimSpace(key)
	if len(key) <= 12 {
		return "(empty/invalid)"
	}
	return key[:12] + "..." + key[len(key)-4:]
}

func logConfigLoaded() {
	fmt.Println()
	fmt.Println("Config: config.props")
	logInfo("KEY_WRITE_MONEY_IN", maskKey(viperGet("KEY_WRITE_MONEY_IN")))
	logInfo("KEY_WRITE_MONEY_OUT", maskKey(viperGet("KEY_WRITE_MONEY_OUT")))
	logInfo("API_VERSION", viperGet("API_VERSION"))
}

func logMoneyInSetup() {
	logInfo("Base URL", "https://api.xendit.co")
	logInfo("Package", "github.com/grosenia/xendit-go-client/v3")
}

func logMoneyOutSetup() {
	logInfo("Base URL", "https://api.xendit.co")
	logInfo("Package", "github.com/grosenia/xendit-go-client/v3/moneyout")
}

func viperGet(key string) string {
	return strings.TrimSpace(viperGetRaw(key))
}

// viperGetRaw avoids import cycle — set by loadConfig wrapper
var viperGetRaw = func(string) string { return "" }

func printUsage() {
	fmt.Println(`Xendit v3 — manual testing (satu config.props)

Usage:
  go run . <command>

━━━ MONEY IN — Payments API v3 ━━━
  card-session      [1] Kartu baru — Cards Session JS
                    API: POST /sessions (mode CARDS_SESSION_JS)
                    Legacy: (baru, bukan credit_card_charges)

  payment-link      [2] Checkout VA / payment link
                    API: POST /sessions (mode PAYMENT_LINK)
                    Legacy: POST /v2/invoices (CreateInvoice)

  reusable-va       [3] Fixed VA seller
                    API: POST /v3/payment_requests (REUSABLE_PAYMENT_CODE)
                    Legacy: POST /callback_virtual_accounts (CreateFixedVa)

  payment-request   [4] Bayar kartu tersimpan + CVN
                    API: POST /v3/payment_requests (type PAY)
                    Butuh: PAYMENT_TOKEN_ID + CVN di config.props

  get-token         [5] Cek token tersimpan
                    API: GET /v3/payment_tokens/{id}
                    Butuh: PAYMENT_TOKEN_ID di config.props

━━━ MONEY OUT — bukan Payments v3 ━━━
  batch-disburse    [6] Batch cairkan uang
                    API: POST /batch_disbursements
                    Legacy: CreateBatchDisbursement (root package)
                    Butuh: KEY_WRITE_MONEY_OUT

Contoh:
  cd v3/examples
  go run . card-session
  go run . payment-link
  go run . reusable-va`)
}

func printTestSummary(cmd string, passed bool) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	if passed {
		fmt.Printf("  HASIL: %s — SUKSES\n", cmd)
	} else {
		fmt.Printf("  HASIL: %s — GAGAL\n", cmd)
	}
	fmt.Printf("  Waktu: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("=", 60))
}
