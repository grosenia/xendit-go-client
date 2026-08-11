package xenditv3

import (
	"encoding/json"
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

// ValidationDetail is one entry in ErrorResponse.Errors. Xendit's own shape
// for this is inconsistent across endpoints AND was confirmed 2026-08-11
// against a real prod v3/payouts response to sometimes be a plain JSON
// string per entry (e.g. "errors": ["recipient.address.postal_code is
// required"]), not the {path, message} object this type originally assumed
// — that mismatch made json.Unmarshal fail the ENTIRE response, wiping out
// even error_code/message and making the diagnostic strictly worse than
// before this field existed. UnmarshalJSON below accepts either shape.
type ValidationDetail struct {
	Path    string `json:"path"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (d *ValidationDetail) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		d.Message = s
		return nil
	}
	type alias ValidationDetail // avoid infinite recursion into this UnmarshalJSON
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*d = ValidationDetail(a)
	return nil
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
