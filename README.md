# xendit-go-client

This is the package in Go language to communicate with Xendit

[![Build Status](https://travis-ci.com/grosenia/xendit-go-client.svg?branch=master)](https://travis-ci.com/grosenia/xendit-go-client) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Supported features:
- Creating Invoice
- Creating Fixed VA
- Credit Card Charge — legacy v2 (`POST /credit_card_charges`)
- Payments API v3 — Cards + Saved Card
  - `POST /sessions` — create payment session (PAY, SAVE, CARDS_SESSION_JS)
  - `GET /sessions/{id}`, `POST /sessions/{id}/cancel`
  - `POST /v3/payment_requests` — pay / pay-and-save / pay with token
  - `GET /v3/payment_requests/{id}`
  - `GET /v3/payment_tokens/{id}`, `POST /v3/payment_tokens/{id}/cancel`
  - Webhook helpers: `payment.capture`, `payment_token.activated`
- Payout
  - Create Payout
  - Void Payout
  - Get Payout
- Create Batch Disbursement

## Examples:
- `example-create-invoice`
- `example-create-fixedva`
- `example-create-credit-card-charge` — legacy v2: tokenize via `tokenize.html`, then charge with `main.go`

### Payments API v3 (package terpisah)

Semua kode, constant, test, docs, dan examples v3 ada di folder **[`v3/`](./v3/)**:

```go
import xenditv3 "github.com/grosenia/xendit-go-client/v3"
```

- Docs: [`v3/README.md`](./v3/README.md)
- Examples: [`v3/examples/`](./v3/examples/)
- Tests: `go test ./v3`

Root package `xenditgo` juga masih punya v3 methods (legacy integrasi core-api). Untuk project baru, pakai `v3/` saja.

### Payments API v3 (root `xenditgo` — existing)

Set `client.ApiVersion = xenditgo.PaymentsAPIVersion` (default `2024-11-11` if empty).

```go
session, err := gateway.CreatePaymentSession(&xenditgo.XenditCreatePaymentSessionReq{
    ReferenceID: orderNo,
    SessionType: xenditgo.SessionTypePay,
    Mode:        xenditgo.SessionModeCardsSessionJS,
    Amount:      total,
    Currency:    "IDR",
    Country:     "ID",
    PaymentTokenID: savedTokenID, // one-click with saved card
    ChannelProperties: &xenditgo.XenditPaymentSessionChannelProperties{
        Cards: &xenditgo.XenditPaymentSessionCardsChannelProperties{
            CardOnFileType: xenditgo.CardOnFileCustomerUnscheduled,
        },
    },
    CardsSessionJS: &xenditgo.XenditPaymentSessionCardsSessionJS{
        SuccessReturnURL: successURL,
        FailureReturnURL: failureURL,
    },
})
// Frontend: cards-session.min.js + CVN → payment_session_id from session.PaymentSessionID
```

Docs: https://docs.xendit.co/docs/cards-one-click-with-cvn
- `example-create-payout`
- `example-create-qrcode`
- `example-create-batch-disbursement`
 
## Code Snippets:
- Sample code to handling payment notifications (coming soon)


Want to help to contribute ? Reach me, can help on documentation, more feature implementation, etc.


