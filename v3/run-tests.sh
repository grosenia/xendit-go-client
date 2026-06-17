#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

MODE="${1:-unit}"

echo "==> unit tests"
go test ./v3 -v -count=1

echo "==> examples smoke tests (fail scenarios)"
go test ./v3/examples -v -count=1 -run 'TestCLI_|TestLogHelpers_'

if [[ "$MODE" == "integration" ]]; then
  read_prop() {
    local file="$1" key="$2"
    [[ -f "$file" ]] || return 0
    grep -E "^${key}=" "$file" | head -1 | cut -d= -f2- | sed 's/\r$//'
  }

  CONFIG="$ROOT/v3/examples/config.props"

  export XENDIT_SECRET_KEY="$(read_prop "$CONFIG" KEY_WRITE_MONEY_IN)"
  export XENDIT_API_VERSION="$(read_prop "$CONFIG" API_VERSION)"
  export XENDIT_BUYER_REFERENCE_ID="$(read_prop "$CONFIG" BUYER_REFERENCE_ID)"
  export XENDIT_PAYER_EMAIL="$(read_prop "$CONFIG" PAYER_EMAIL)"
  export XENDIT_PAYMENT_TOKEN_ID="$(read_prop "$CONFIG" PAYMENT_TOKEN_ID)"
  export XENDIT_CVN="$(read_prop "$CONFIG" CVN)"

  echo "==> integration tests (sandbox)"
  go test ./v3 -tags=integration -v -count=1 -run Integration

  echo "==> integration fail tests (sandbox — harapan ErrorStatus=true)"
  go test ./v3 -tags=integration -v -count=1 -run IntegrationFail
fi
