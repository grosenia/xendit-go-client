package moneyout

import "fmt"

// ErrorResponse is the standard Xendit error body.
type ErrorResponse struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"message"`
	ErrorStatus  bool   `json:"-"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("[%s] %s", e.ErrorCode, e.ErrorMessage)
}

func (e *ErrorResponse) markHTTPError(httpStatus int) {
	if e == nil {
		return
	}
	if httpStatus != 200 && httpStatus != 201 {
		e.ErrorStatus = true
		if e.ErrorCode == "" && e.ErrorMessage == "" {
			e.ErrorMessage = fmt.Sprintf("unexpected http status %d", httpStatus)
		}
	} else {
		e.ErrorStatus = false
	}
}

// DisbursementItem is one row in a batch disbursement.
type DisbursementItem struct {
	Amount            float64 `json:"amount"`
	ExternalID        string  `json:"external_id"`
	BankCode          string  `json:"bank_code"`
	BankAccountName   string  `json:"bank_account_name"`
	BankAccountNumber string  `json:"bank_account_number"`
	Description       string  `json:"description"`
}

// CreateBatchDisbursementRequest is JSON for POST /batch_disbursements.
type CreateBatchDisbursementRequest struct {
	Reference     string             `json:"reference"`
	Disbursements []DisbursementItem `json:"disbursements"`
}

// CreateBatchDisbursementResponse is JSON from POST /batch_disbursements.
type CreateBatchDisbursementResponse struct {
	ID                  string `json:"id"`
	Reference           string `json:"reference"`
	Status              string `json:"status"`
	Created             string `json:"created"`
	TotalUploadedCount  int64  `json:"total_uploaded_count"`
	TotalUploadedAmount int64  `json:"total_uploaded_amount"`
	ErrorResponse
}

// CreatePayoutRequest is JSON for POST /payouts.
type CreatePayoutRequest struct {
	ExternalID string  `json:"external_id"`
	Amount     float64 `json:"amount"`
}

// PayoutResponse is JSON from payout APIs.
type PayoutResponse struct {
	ID               string  `json:"id"`
	ExternalID       string  `json:"external_id"`
	Amount           float64 `json:"amount"`
	Status           string  `json:"status"`
	PayoutURL        string  `json:"payout_url,omitempty"`
	CreatedTimestamp string  `json:"created"`
	ErrorResponse
}
