# Testing — Payments API v3

Semua testing v3 ada di folder ini. **Production `core-api` pakai root `payment_v3*`** — folder `v3/` dipakai untuk develop & verify sebelum migrasi.

## 1. Unit test (tanpa network)

Dari root repo:

```bash
go test ./v3 -v
```

Cakupan:
- URL builder (`urls_test.go`)
- Webhook helpers (`gateway_test.go`)
- Customer JSON format (`customer_test.go`)
- Response helpers — masked card, error (`types_response_test.go`)
- HTTP mock — session, customer, payment request (`gateway_httptest_test.go`)

## 2. Integration test (Xendit sandbox)

Butuh secret key sandbox. Set env lalu jalankan:

```bash
export XENDIT_SECRET_KEY=xnd_development_...

# Wajib untuk test session + customer
go test ./v3 -tags=integration -v -run IntegrationCreateCustomerAndPaymentSession

# Opsional — saved card (isi setelah save card di sandbox)
export XENDIT_PAYMENT_TOKEN_ID=pt_...
export XENDIT_CVN=123
go test ./v3 -tags=integration -v -run IntegrationGetPaymentToken
go test ./v3 -tags=integration -v -run IntegrationCreatePaymentRequestWithSavedCard

# Semua integration sekaligus
go test ./v3 -tags=integration -v -run Integration
```

Env opsional:

| Env | Default |
|-----|---------|
| `XENDIT_API_VERSION` | `2024-11-11` |
| `XENDIT_BUYER_REFERENCE_ID` | auto timestamp |
| `XENDIT_PAYER_EMAIL` | `v3-integration@grosenia.co.id` |

Atau pakai script:

```bash
./v3/run-tests.sh              # unit only
./v3/run-tests.sh integration  # sandbox (baca config.props example)
```

## 3. Manual examples (satu folder)

```bash
cd v3/examples
go run . card-session
go run . payment-link
go run . reusable-va
go run . payment-request   # butuh PAYMENT_TOKEN_ID + CVN di config.props
go run . get-token
go run . batch-disburse
```

## 4. Matriks sukses vs gagal

### Unit (httptest — tanpa network)

| # | Skenario | Sukses | Gagal |
|---|----------|--------|-------|
| 1 | card-session | `TestScenario01_CardSession_Success` | `TestScenario01_CardSession_Fail_InvalidCustomer` |
| 2 | payment-link | `TestScenario02_PaymentLink_Success` | `TestScenario02_PaymentLink_Fail_MissingPaymentMethods` + `Fail_InvalidAmount` |
| 3 | reusable-va | `TestScenario03_ReusableVA_Success` | `TestScenario03_ReusableVA_Fail_InvalidBank` + build fail bank |
| 4 | payment-request | `TestScenario04_PaymentRequest_Success` | `TestScenario04_PaymentRequest_Fail_InvalidToken` |
| 5 | get-token | `TestScenario05_GetToken_Success` | `TestScenario05_GetToken_Fail_NotFound` |
| 6 | batch-disburse | `moneyout.TestCreateBatchDisbursementHTTP` | `moneyout.TestCreateBatchDisbursementHTTP_Fail_Unauthorized` |

```bash
go test ./v3 -v -run 'TestScenario'
go test ./v3/moneyout -v
```

### Integration sandbox

**Sukses** (`ErrorStatus=false`):

```bash
go test ./v3 -tags=integration -v -run 'Integration' -run 'Fail' --skip  # use -run Integration only
go test ./v3 -tags=integration -v -run Integration
```

**Gagal** (`ErrorStatus=true`, error_code dari Xendit):

```bash
go test ./v3 -tags=integration -v -run IntegrationFail
```

### Examples CLI smoke (gagal tanpa network)

```bash
cd v3/examples
go test -v -run 'TestCLI_|TestLogHelpers_'
```

Harapan: exit code `1`, output berisi `GAGAL` dan pesan validasi.

