package xenditv3

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
