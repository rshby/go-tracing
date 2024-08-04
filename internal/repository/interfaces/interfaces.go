package interfaces

import (
	"context"
	"go-tracing/database/model"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(ctx context.Context, tx *gorm.DB, input *model.Customer) (*model.Customer, error)
	GetByID(ctx context.Context, ID uint) (*model.Customer, error)
}
