package dto

type ApiResponseDTO struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
	ErrorCode  string `json:"error_code,omitempty"`
	Data       any    `json:"data,omitempty"`
}
