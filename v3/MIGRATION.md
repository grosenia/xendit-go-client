# Migrasi ke `v3/` — Roadmap Grosenia

Tujuan: **satu package** `github.com/grosenia/xendit-go-client/v3` untuk semua integrasi Xendit, dengan Money In lewat **Payments API v3** dan Money Out tetap API disbursement (bukan Payments v3).

## Yang dimaksud "migrasi ke v3"

| Dimensi | Artinya |
|---------|---------|
| **Library** | Semua kode client pindah ke folder `v3/` — root `xenditgo` di-deprecate |
| **API Xendit (Money In)** | Invoice, VA, QR, kartu → Payments API v3 (`/sessions`, `/v3/payment_requests`, …) |
| **API Xendit (Money Out)** | Payout / batch disbursement — **API berbeda**, bukan Payments v3, tapi tetap dipindah ke `v3/moneyout/` |

---

## Mapping legacy → Payments API v3

| Grosenia sekarang | Legacy API | Target v3 | Routing type / mode |
|-------------------|------------|-----------|---------------------|
| Checkout VA (invoice link) | `POST /v2/invoices` | `POST /sessions` | `mode: PAYMENT_LINK` |
| VA per bank di invoice | `payment_methods: ["BCA", …]` | `allowed_payment_channels` di session | Payment Session |
| Fixed VA seller | `POST /callback_virtual_accounts` | `POST /v3/payment_requests` | `type: REUSABLE_PAYMENT_CODE` + `channel_code: *_VIRTUAL_ACCOUNT` |
| Mandiri open VA | `POST /v2/payment_methods` | `POST /v3/payment_requests` | channel Mandiri |
| QRIS | `POST /qr_codes` | `POST /v3/payment_requests` | `channel_code: QRIS` (PAY atau REUSABLE) |
| Kartu kredit v3 | `/sessions` + `/v3/payment_requests` | ✅ **sudah ada** | `CARDS_SESSION_JS` |
| Kartu legacy v2 | `POST /credit_card_charges` | hapus setelah kartu v3 stabil | — |
| Payout admin | `POST /payouts` | `v3/moneyout/` (API sama) | bukan Payments v3 |
| Batch disbursement | `POST /batch_disbursements` | `v3/moneyout/` (API sama) | bukan Payments v3 |

**Webhook lama → baru:**

| Legacy | v3 |
|--------|-----|
| Invoice paid callback | `payment.capture` |
| Fixed VA paid | `payment.capture` |
| QR callback | `payment.capture` |
| `credit_card.*` | `payment.capture` |
| — | `payment_token.activated` |
| — | `payment_session.completed` |

Ref Xendit:
- [Payments API overview](https://docs.xendit.co/docs/payments-via-api-overview)
- [Migrate Payment Link → Payment Session](https://docs.xendit.co/docs/migrate-to-payment-session)
- [Create payment request](https://docs.xendit.co/apidocs/create-payment-request)

---

## Fase implementasi

### Fase 0 — Kartu (DONE)

- [x] `v3/` package: session, payment request, token, customer, webhook
- [x] Examples + unit + integration test
- [ ] `core-api` import `xenditv3` (masih pakai root `payment_v3*`)

### Fase 1 — Switch `core-api` kartu ke `v3/`

- [ ] `go.mod` core-api → `xendit-go-client/v3`
- [ ] Ganti type `xenditgo.XenditCreatePaymentSessionReq` → `xenditv3.CreatePaymentSessionRequest`
- [ ] Webhook handlers pakai `xenditv3.PaymentCaptureCallback`
- [ ] Hapus duplikat root `payment_v3*` setelah verified

**Impact:** card payment only. VA/invoice/QR tidak berubah.

### Fase 2 — Invoice / VA checkout (order payment)

**core-api:** `xendit_service.go` — `CreateInvoice`, `xenditCreateInvoiceNonFixedVA`, checkout invoice

- [ ] Tambah di `v3/`: Payment Session `PAYMENT_LINK` + `allowed_payment_channels`
- [ ] Map `external_id` → `reference_id`, `invoice_duration` → `expires_at`
- [ ] Example: `v3/examples/create-payment-link-session/`
- [ ] Spike sandbox: BCA, Mandiri, BNI VA via session
- [ ] Dual-write / feature flag: legacy invoice vs session (staging)
- [ ] Webhook: handle `payment.capture` untuk order VA (parallel dengan invoice callback)

**Bank yang dipakai Grosenia:** BCA, Mandiri, BNI, Permata, BRI (cek `xendit_service.go` switch bank)

### Fase 3 — Fixed VA seller

**core-api:** `CreateFixedVa`, `UpdateFixedVa`, business fixed VA

- [ ] `v3/`: Payment Request `REUSABLE_PAYMENT_CODE` + channel properties (`expires_at`, `display_name`, verification data BRI)
- [ ] Example: `v3/examples/create-reusable-va/`
- [ ] Migrasi FVA expiration webhook ke `payment.capture`

### Fase 4 — QRIS

- [ ] `v3/`: Payment Request `QRIS` channel
- [ ] Example: `v3/examples/create-qris-payment-request/`
- [ ] Ganti `create_qr.go` flow di core-api

### Fase 5 — Bersihkan legacy Money In

- [ ] Hapus `invoice.go` (bagian money in), `create_qr.go`, `credit_card.go` dari root
- [ ] Hapus `example-create-invoice`, `example-create-fixedva`, `example-create-qrcode`, `example-create-credit-card-charge` di root → pindah ke `v3/examples/legacy/` atau hapus

### Fase 6 — Money Out ke `v3/moneyout/`

Payout **bukan** Payments API v3 — hanya reorganisasi library:

- [ ] `v3/moneyout/payout.go` — copy dari root `payout.go`
- [ ] `v3/moneyout/disbursement.go` — copy dari root `disbursement.go`
- [ ] Examples: `v3/examples/create-payout/`, `create-batch-disbursement/`
- [ ] Migrasi `core-api` payout/disbursement import

### Fase 7 — Deprecate root package

- [ ] Root `xenditgo` hanya re-export `v3` (optional compatibility shim) atau hapus total
- [ ] Tag major version `v2.0.0`

---

## Struktur folder target

```
v3/
├── client.go, gateway.go, constants.go, errors.go, urls.go
├── types_request.go, types_response.go
├── customer.go, webhook.go
├── moneyout/              # fase 6 — API disbursement (bukan Payments v3)
│   ├── payout.go
│   └── disbursement.go
├── examples/
│   ├── main.go              # go run . <command>
│   ├── config.props         # satu config semua fitur
│   └── README.md
├── TESTING.md
├── MIGRATION.md                     # file ini
└── *_test.go
```

---

## Urutan disarankan untuk Grosenia

1. **Fase 1** dulu (kartu) — risiko kecil, sudah production
2. **Fase 2** invoice VA — impact terbesar (checkout Android/web)
3. **Fase 3** fixed VA seller
4. **Fase 4** QRIS
5. **Fase 6** moneyout reorganisasi (bisa parallel)
6. **Fase 7** cleanup

Setiap fase: sandbox test → staging feature flag → production → hapus legacy path.

---

## Yang perlu diingat

- **BRI VA (Aug 2025+):** wajib verification data di `channel_properties` — rencanakan sebelum fase 3.
- **Installment di Payment Session:** belum support Xendit (Q1 2026) — cek kalau Grosenia pakai cicilan.
- **Payout tidak masuk Payments v3** — tetap API `/payouts` dan `/batch_disbursements`.

Mulai dari fase mana? Rekomendasi: **Fase 1** (core-api kartu → `v3/`) lalu **Fase 2 spike** invoice/VA di sandbox.
