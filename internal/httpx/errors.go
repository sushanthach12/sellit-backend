package httpx

import (
	"encoding/json"
	"net/http"
)

/*
Status Code		Custom Code
400				invalid_id
404				not_found
500				internal_error
400 			malformed_json
422 			validation_failed
401 			unauthenticated
403 			forbidden
409 			conflict
429 			rate_limited
*/

// ErrorCode is a sealed enum: the underlying string field is unexported,
// so code outside this package cannot construct an ErrorCode directly
// (httpx.ErrorCode{value: "whatever"} won't compile from another package).
// The only valid values are the exported Code* constants below.
type ErrorCode struct {
	value string
}

// String lets ErrorCode be printed/logged normally.
func (c ErrorCode) String() string {
	return c.value
}

// MarshalJSON controls how it's serialized in the response body.
func (c ErrorCode) MarshalJSON() ([]byte, error) { // required: otherwise the value is not serialized inside of json encoder
	return json.Marshal(c.value)
}

var (
	CodeInvalidId        = ErrorCode{"invalid_id"}
	CodeNotFound         = ErrorCode{"not_found"}
	CodeInternalError    = ErrorCode{"internal_error"}
	CodeMalformedJson    = ErrorCode{"malformed_json"}
	CodeValidationFailed = ErrorCode{"validation_failed"}
	CodeUnAuthenticated  = ErrorCode{"unauthenticated"}
	CodeForbidden        = ErrorCode{"forbidden"}
	CodeConflict         = ErrorCode{"conflict"}
	CodeRateLimited      = ErrorCode{"rate_limited"}
)

type errorPayload struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type errorResponse struct {
	Error errorPayload `json:"error"`
}

func Error(w http.ResponseWriter, status_code int, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: errorPayload{
			Code:    code,
			Message: message,
		},
	})
}

// WriteJSON encodes v as JSON and writes it to w with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return // 204 must not have a body
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Nothing more we can do once headers are written, but at least log it
		// if you have a logger accessible here — otherwise this is a safe no-op.
		_ = err
	}
}
