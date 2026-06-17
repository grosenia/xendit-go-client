# V3 Examples — Manual Testing

Satu folder, satu `config.props`, log jelas per bagian.

```bash
cd v3/examples
go run . <command>
```

## Daftar testing

| # | Command | Bagian yang ditest | API |
|---|---------|-------------------|-----|
| 1 | `card-session` | Kartu baru + save card | `POST /sessions` · CARDS_SESSION_JS |
| 2 | `payment-link` | Checkout VA / payment link (ganti invoice) | `POST /sessions` · PAYMENT_LINK |
| 3 | `reusable-va` | Fixed VA seller (ganti fixed VA) | `POST /v3/payment_requests` · REUSABLE_PAYMENT_CODE |
| 4 | `payment-request` | Bayar kartu tersimpan + CVN | `POST /v3/payment_requests` · PAY |
| 5 | `get-token` | Cek token tersimpan | `GET /v3/payment_tokens/{id}` |
| 6 | `batch-disburse` | Cairkan uang (Money Out) | `POST /batch_disbursements` |

## Urutan disarankan

```
1. go run . card-session      → dapat payment_session_id
   (bayar + save card di frontend sandbox)

2. go run . get-token         → cek pt_xxx (isi PAYMENT_TOKEN_ID dulu)

3. go run . payment-request   → bayar lagi pakai token + CVN

--- VA / invoice (independen) ---

4. go run . payment-link      → dapat payment_url (ganti invoice)

5. go run . reusable-va       → dapat fixed VA seller

--- Money out ---

6. go run . batch-disburse    → butuh KEY_WRITE_MONEY_OUT
```

Setiap command menampilkan log: **TEST banner → config → step API → hasil → SUKSES/GAGAL**.

Automated test: `../run-tests.sh` dan `../TESTING.md`.
