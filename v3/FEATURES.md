# Xendit Go Client — Daftar Fitur

Library ini punya **dua lapisan**:

| Lapisan | Import | Keterangan |
|---------|--------|------------|
| **Legacy (existing)** | `github.com/grosenia/xendit-go-client` | API lama, **tidak diubah** — tetap dipakai production |
| **Payments API v3 (baru)** | `github.com/grosenia/xendit-go-client/v3` | Package terpisah untuk kartu kredit/debit v3 |

---

## Legacy — `package xenditgo` (root)

### Money In — Invoice & VA

| Fitur | File | Endpoint Xendit |
|-------|------|-----------------|
| Create Invoice | `invoice.go` | `POST /v2/invoices` |
| Expire Invoice | `invoice.go` | `POST /v2/invoices/{id}/expire` |
| Create Fixed VA | `invoice.go` | `POST /callback_virtual_accounts` |
| Update Fixed VA | `invoice.go` | `PATCH /callback_virtual_accounts/{id}` |
| Get Fixed VA | `invoice.go` | `GET /callback_virtual_accounts/{id}` |
| Payment Method (open VA Mandiri dll.) | `invoice.go` | `POST /v2/payment_methods` |

**Example:** `example-create-invoice/`, `example-create-fixedva/`

### Credit Card — Legacy v2

| Fitur | File | Endpoint Xendit |
|-------|------|-----------------|
| Credit Card Charge | `credit_card.go` | `POST /credit_card_charges` |

**Example:** `example-create-credit-card-charge/`

### QRIS

| Fitur | File | Endpoint |
|-------|------|----------|
| Create QR Code | `create_qr.go` | QR API |
| Pay QR Code | `create_qr.go` | QR API |

**Example:** `example-create-qrcode/`

### Money Out

| Fitur | File | Endpoint |
|-------|------|----------|
| Create Payout | `payout.go` | `POST /disbursements` |
| Get Payout | `payout.go` | `GET /disbursements/{id}` |
| Void Payout | `payout.go` | `POST /disbursements/{id}/void` |
| Batch Disbursement | `disbursement.go` | Batch API |

**Example:** `example-create-payout/`, `example-create-batch-disbursement/`

### Webhook helpers (legacy)

| File | Event |
|------|-------|
| `callback.go` | Invoice paid, Fixed VA, QR, credit_card.* |

---

## Payments API v3 — `package xenditv3` (folder `v3/`)

Package **baru**, tidak mengganti root package. **Isi `payment_v3*.go` di root sudah diduplikasi di sini** — ini sumber kebenaran untuk kode baru.

| Root (duplikat, jangan dipakai untuk kode baru) | Di `v3/` |
|-------------------------------------------------|----------|
| `payment_v3_constants.go` | `constants.go` |
| `payment_v3_request.go` | `types_request.go` |
| `payment_v3_response.go` | `types_response.go` |
| `payment_v3.go` | `gateway.go` + `client.go` |
| `payment_v3_test.go` | `gateway_test.go`, `urls_test.go` |
| `customer.go` (bagian Create/Get customer) | `customer.go` |
| `callback.go` (webhook payment.capture / payment_token) | `webhook.go` |
| URL di `envtype.go` (`CreatePaymentSessionURL`, dll.) | `urls.go` |

**Kenapa file root `payment_v3*` belum dihapus?**  
`core-api` masih import `github.com/grosenia/xendit-go-client` (`package xenditgo`) dan memanggil `gateway.CreatePaymentSession`, `xenditgo.SessionTypePay`, dll. Kalau file root dihapus sekarang, build `core-api` akan rusak.

**Rencana:** setelah `core-api` pindah ke `import xenditv3 "github.com/grosenia/xendit-go-client/v3"`, file `payment_v3*` di root bisa di-deprecate lalu dihapus.

| Fitur | Method v3 | Legacy |
|-------|-----------|--------|
| Payment link (invoice) | `CreatePaymentLinkSession` | `CreateInvoice` |
| Fixed / reusable VA | `CreateReusableVirtualAccount` | `CreateFixedVa` |
| Batch disbursement | `moneyout.CreateBatchDisbursement` | `CreateBatchDisbursement` |
| Kartu session / token | `CreatePaymentSession`, … | `payment_v3*` root |

| Endpoint | Path |
|----------|------|
| Payment Session | `/sessions` |
| Payment Request | `/v3/payment_requests` |
| Payment Token | `/v3/payment_tokens` |
| Customer | `/customers` |
| Batch disbursement | `/batch_disbursements` (moneyout) |
| Webhook | `payment.capture`, `payment_token.activated`, `payment_session.completed` |

**Docs:** [README.md](./README.md)

**Examples:** [`examples/`](./examples/) — satu `main.go`, pilih command:

```bash
cd v3/examples && go run . payment-link
```

**Tests:** `*_test.go` (package ini)

---

## Migrasi ke v3 (rencana)

**Roadmap lengkap:** [MIGRATION.md](./MIGRATION.md)

Ringkas:
1. **Kartu** — sudah di `v3/` ✅
2. **Invoice / VA checkout** → Payment Session `PAYMENT_LINK`
3. **Fixed VA** → Payment Request `REUSABLE_PAYMENT_CODE`
4. **QRIS** → Payment Request `QRIS`
5. **Payout / disbursement** → `v3/moneyout/` (API disbursement, bukan Payments v3)
6. Hapus root `xenditgo` legacy
