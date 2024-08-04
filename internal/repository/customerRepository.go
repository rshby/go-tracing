package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-tracing/database"
	"go-tracing/database/model"
	"go-tracing/internal/repository/interfaces"
	"go-tracing/internal/utils/helper"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository is method to create customerRepository
func NewCustomerRepository() interfaces.CustomerRepository {
	return &customerRepository{
		db: database.MysqlDB,
	}
}

func (c *customerRepository) Create(ctx context.Context, tx *gorm.DB, input *model.Customer) (*model.Customer, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"input": helper.Dump(input),
	})

	if tx == nil {
		tx = c.db
	}

	if err := tx.WithContext(ctx).Model(&model.Customer{}).Create(input).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return input, nil
}

func (c *customerRepository) GetByID(ctx context.Context, ID uint) (*model.Customer, error) {
	logger := logrus.WithContext(ctx).WithField("id", ID)

	var customer model.Customer
	err := c.db.WithContext(ctx).Model(&model.Customer{}).Take(&customer, "id = ?", ID).Error

	switch err {
	case nil:
		return &customer, nil
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}
}
