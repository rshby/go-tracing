package httpresponse

import (
	"github.com/gin-gonic/gin"
	"go-tracing/internal/http/httpresponse/dto"
	"net/http"
)

func ResponseBadRequest(c *gin.Context) {
	statusCode := http.StatusBadRequest
	c.Status(statusCode)
	c.JSON(statusCode, NewApiResponse("", ErrorBadRequest, nil))
}

func ResponseError(c *gin.Context, e *HttpError) {
	statusCode := e.StatusCode
	c.Status(statusCode)
	c.JSON(statusCode, NewApiResponse("", e, nil))
}

func ResponseOK(c *gin.Context, msg string, data any) {
	statusCode := http.StatusOK
	c.Status(statusCode)
	c.JSON(statusCode, NewApiResponse(msg, nil, data))
}

func NewApiResponse(msg string, e *HttpError, data any) *dto.ApiResponseDTO {
	response := dto.ApiResponseDTO{
		StatusCode: http.StatusOK,
		Message:    msg,
		Data:       data,
	}

	if e != nil {
		response.StatusCode = e.StatusCode
		response.Message = e.Error()
		response.ErrorCode = e.ErrorCode
	}

	return &response
}
