package xenditv3

import (
	"fmt"
	"strings"
)

// ErrorResponse is the standard Xendit error body. Errors holds the
// field-level validation detail Xendit's API_VALIDATION_ERROR responses
// include alongside the generic top-level message ("Inputs are failing
// validation. The errors field contains details...") — without this,
// callers only ever see that generic sentence and never learn which field
// actually failed. Confirmed missing 2026-08-10 against a real prod
// v3/payouts validation failure.
type ErrorResponse struct {
	ErrorCode    string             `json:"error_code"`
	ErrorMessage string             `json:"message"`
	Errors       []ValidationDetail `json:"errors,omitempty"`
	ErrorStatus  bool               `json:"-"`
}

// ValidationDetail is one entry in ErrorResponse.Errors — Xendit's field path
// + message shape is not consistently documented across endpoints, so both
// fields are read loosely (path is sometimes "field" in older responses).
type ValidationDetail struct {
	Path    string `json:"path"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (d ValidationDetail) String() string {
	field := d.Path
	if field == "" {
		field = d.Field
	}
	if field == "" {
		return d.Message
	}
	return field + ": " + d.Message
}

func (e ErrorResponse) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("[%s] %s", e.ErrorCode, e.ErrorMessage)
	}
	details := make([]string, len(e.Errors))
	for i, d := range e.Errors {
		details[i] = d.String()
	}
	return fmt.Sprintf("[%s] %s (%s)", e.ErrorCode, e.ErrorMessage, strings.Join(details, "; "))
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
