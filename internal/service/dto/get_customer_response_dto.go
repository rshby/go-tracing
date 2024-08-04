package dto

import (
	"go-tracing/database/model"
	"go-tracing/internal/utils/helper"
)

type GetCustomerResponseDTO struct {
	ID             uint   `json:"id,omitempty"`
	IdentityNumber string `json:"identity_number"`
	FullName       string `json:"full_name,omitempty"`
	Status         string `json:"status,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

func CustomerEntityToResponseDTO(input *model.Customer) *GetCustomerResponseDTO {
	response := GetCustomerResponseDTO{
		ID:             input.ID,
		IdentityNumber: input.IdentityNumber,
		FullName:       input.FullName,
		Status:         input.Status,
		CreatedAt:      helper.TimeToStringIndonesia(input.CreatedAt),
		UpdatedAt:      helper.TimeToStringIndonesia(input.UpdatedAt),
	}

	return &response
}
