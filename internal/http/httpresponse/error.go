package httpresponse

import "net/http"

type HttpError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error_code"`
}

func (h *HttpError) Error() string {
	return h.Message
}

func (h *HttpError) WithStatusCode(v int) *HttpError {
	h.StatusCode = v
	return h
}

func (h *HttpError) WithMessage(v string) *HttpError {
	h.Message = v
	return h
}

func (h *HttpError) WithErrorCode(v string) *HttpError {
	h.ErrorCode = v
	return h
}

func NewHttError() *HttpError {
	return &HttpError{}
}

var (
	// 400
	ErrorBadRequest = &HttpError{StatusCode: http.StatusBadRequest, Message: "bad request", ErrorCode: "ERR400001"}

	// 404
	ErrorCustomerNotFound = &HttpError{StatusCode: http.StatusNotFound, Message: "customer not found", ErrorCode: "ERR404001"}

	// 500
	ErrorInternalServerError = &HttpError{StatusCode: http.StatusInternalServerError, Message: "internal server error", ErrorCode: "ERR500001"}
)
