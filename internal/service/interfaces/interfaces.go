package interfaces

import (
	"context"
	"go-tracing/internal/http/httpresponse"
	"go-tracing/internal/service/dto"
)

type CustomerService interface {
	Create(ctx context.Context, request *dto.CreateCustomerRequestDTO) *httpresponse.HttpError
	GetByID(ctx context.Context, ID uint) (*dto.GetCustomerResponseDTO, *httpresponse.HttpError)
}
