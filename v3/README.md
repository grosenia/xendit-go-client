# Payments API v3 — `github.com/grosenia/xendit-go-client/v3`

Package terpisah dari root `xenditgo`. **File legacy tidak diubah.**

Base URL: `https://api.xendit.co`  
Header: `api-version: 2024-11-11` (default)

## Install

```go
import xenditv3 "github.com/grosenia/xendit-go-client/v3"
```

## Quick start — Payment Session (kartu baru)

```go
client := xenditv3.NewClient("xnd_development_...")
gw := xenditv3.NewGateway(client)

session, err := gw.CreatePaymentSession(&xenditv3.CreatePaymentSessionRequest{
    ReferenceID:   "ORD-001",
    SessionType:   xenditv3.SessionTypePay,
    Mode:          xenditv3.SessionModeCardsSessionJS,
    Amount:        100000,
    Currency:      "IDR",
    Country:       "ID",
    CaptureMethod: xenditv3.CaptureMethodAutomatic,
    Customer: &xenditv3.SessionCustomer{
        ReferenceID: "buyer-123",
        Type:        "INDIVIDUAL",
        Email:       "buyer@example.com",
        IndividualDetail: &xenditv3.SessionIndividualDetail{
            GivenNames: "Buyer",
            Surname:    "Name",
        },
    },
    CardsSessionJS: &xenditv3.CardsSessionJS{
        SuccessReturnURL: "https://example.com/success",
        FailureReturnURL: "https://example.com/failure",
    },
})
// Frontend: cards-session.min.js + session.PaymentSessionID
```

## Quick start — Payment Request (kartu tersimpan + CVN)

```go
pr, err := gw.CreatePaymentRequest(&xenditv3.CreatePaymentRequestRequest{
    ReferenceID:    "ORD-002",
    PaymentTokenID: "pt_...",
    Type:           xenditv3.PaymentRequestTypePay,
    Country:        "ID",
    Currency:       "IDR",
    RequestAmount:  100000,
    CaptureMethod:  xenditv3.CaptureMethodAutomatic,
    ChannelProperties: &xenditv3.PaymentRequestChannelProperties{
        CardOnFileType:      xenditv3.CardOnFileCustomerUnscheduled,
        TransactionSequence: xenditv3.TransactionSequenceSubsequent,
        CardDetails:         &xenditv3.PaymentRequestCardDetails{Cvn: "123"},
    },
})
```

## Examples

Lihat [examples/README.md](./examples/README.md).

```bash
cd v3/examples
go run . card-session
./v3/run-tests.sh integration   # pakai v3/examples/config.props
```

## Tests

Lihat **[TESTING.md](./TESTING.md)** — unit test, integration sandbox, dan manual examples.

```bash
go test ./v3 -v
./v3/run-tests.sh integration   # pakai config.props dari examples
```
