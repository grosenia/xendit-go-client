# Example: Create Credit Card Charge

Sandbox example for `InvoiceGateway.CreateCreditCardCharge()` (`POST /credit_card_charges`).

Card data must **not** go through this Go client. Tokenize in the browser with Xendit.js, then charge using `token_id` (+ `authentication_id` when 3DS is required).

## 1. Prepare config

```bash
cp config-sample.props config.props
```

Edit `config.props`:

```properties
KEY_WRITE_MONEY_IN=xnd_development_xxx
PUBLIC_API_KEY=xnd_public_development_xxx
CHARGE_AMOUNT=100000
```

## 2. Generate token_id (sandbox)

1. Open `tokenize.html` in a browser (double-click or serve locally).
2. Paste `PUBLIC_API_KEY` from Xendit Dashboard.
3. Set **Amount** equal to `CHARGE_AMOUNT` in config.
4. Click **Create Token**.
5. Copy `token_id` and `authentication_id` (if present) into `config.props`:

```properties
TOKEN_ID=653f...
AUTHENTICATION_ID=auth-...
```

Sandbox 3DS test card (default in HTML): `4000000000001091`, CVN `123`, exp any future date.

## 3. Run charge

From repo root:

```bash
go run ./example-create-credit-card-charge/main.go
```

Or from this folder:

```bash
go run main.go
```

Optional CLI override:

```bash
go run main.go <token_id> [authentication_id] [amount] [external_id]
```

## 4. Handle 3DS

If the response includes `payer_authentication_url`, open that URL in a browser to finish authentication.

## Flow

```
tokenize.html (Xendit.js) → token_id
        ↓
main.go (secret key) → POST /credit_card_charges
        ↓
3DS URL (if needed) → webhook credit_card.*
```

## References

- [Xendit Tokenization](https://docs.xendit.co/credit-cards/integrations/tokenization)
- [Xendit Create Charge](https://docs.xendit.co/credit-cards/integrations/charges)
