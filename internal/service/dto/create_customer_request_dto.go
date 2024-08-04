package dto

import (
	"go-tracing/database/model"
)

type CreateCustomerRequestDTO struct {
	IdentityNumber string `json:"identity_number" validate:"omitempty"`
	FullName       string `json:"full_name" validate:"omitempty"`
}

func (c *CreateCustomerRequestDTO) ToEntity() *model.Customer {
	input := model.Customer{
		IdentityNumber: c.IdentityNumber,
		FullName:       c.FullName,
		Status:         "created",
	}

	return &input
}
