# v3 — Checklist Testing (sebelum core-api)

Status audit `xendit-go-client/v3` — fokus **kartu kredit/debit Payments API v3** saja.

## Gateway methods — coverage

| Method | Unit | Integration | Example (`go run . …`) |
|--------|------|-------------|------------------------|
| `CreatePaymentLinkSession` | ✅ | ✅ | `payment-link` |
| `CreateReusableVirtualAccount` | ✅ | ✅ | `reusable-va` |
| `CreateCustomer` | ✅ | ✅ | `card-session` |
| `CreatePaymentSession` | ✅ | ✅ | `card-session` |
| `CreatePaymentRequest` | ✅ | ⏭ skip | `payment-request` |
| `GetPaymentToken` | ✅ | ⏭ skip | `get-token` |
| `moneyout.CreateBatchDisbursement` | ✅ | — | `batch-disburse` |

## Helpers — coverage

| Helper | Unit test |
|--------|-----------|
| URL builders | ✅ `urls_test.go`, `customer_test.go` |
| `NormalizeCustomerID` | ✅ |
| `ErrorResponse` | ✅ |
| `PaymentTokenResponse.MaskedCardNumber` | ✅ |
| `PaymentRequestResponse.AuthenticationURL` | ✅ |
| `PaymentCaptureCallback` | ✅ |
| `PaymentTokenCallback` | ✅ |
| `PaymentSessionCompletedCallback` | ✅ |

## Parity vs root `payment_v3*`

| Item | v3/ | Root | Catatan |
|------|-----|------|---------|
| Session CRUD | ✅ | ✅ | Parity OK |
| Payment request CRUD | ✅ | ✅ | v3 punya `channel_code` extra (fase migrasi VA/QR) |
| Payment token CRUD | ✅ | ✅ | Parity OK |
| Customer API | ✅ | ✅ | v3 fix: `given_names` di root body POST /customers |
| Webhook helpers | ✅ | partial | v3 tambah `PaymentSessionCompletedCallback` |
| CC legacy v2 | — | `credit_card.go` | Bukan scope v3 |

## Integration sandbox — hasil terakhir

```bash
./v3/run-tests.sh              # 22+ unit tests
./v3/run-tests.sh integration  # customer + session + get + cancel PASS
```

| Test | Status | Catatan |
|------|--------|---------|
| `TestIntegrationCreateCustomerAndPaymentSession` | ✅ PASS | create → get → cancel |
| `TestIntegrationGetPaymentToken` | ⏭ SKIP | isi `PAYMENT_TOKEN_ID` di config.props |
| `TestIntegrationCreatePaymentRequestWithSavedCard` | ⏭ SKIP | butuh token + CVN |

## Yang belum (OK untuk sekarang — belum ke core-api)

- [ ] Example manual saved card end-to-end (butuh `pt_xxx` dari sandbox)
- [ ] Integration `CancelPaymentToken` (destructive — jangan auto-run)
- [ ] VA / invoice / QR — fase migrasi berikutnya ([MIGRATION.md](./MIGRATION.md))

## Gate ke core-api

Lanjut ke `core-api` **setelah**:

1. ✅ Semua unit test `./v3` PASS
2. ✅ Integration session flow PASS (create → get → cancel)
3. ⏳ Integration saved card PASS — isi `PAYMENT_TOKEN_ID` + `CVN` di `examples/create-payment-request/config.props`, lalu `./v3/run-tests.sh integration`

```bash
go test ./v3 -v -count=1
./v3/run-tests.sh integration
cd v3/examples && go run . card-session
```
