package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-tracing/database"
	"go-tracing/internal/http/httpresponse"
	repositoryInterfaces "go-tracing/internal/repository/interfaces"
	"go-tracing/internal/service/dto"
	"go-tracing/internal/service/interfaces"
	"go-tracing/internal/utils/helper"
	"gorm.io/gorm"
)

type customerServices struct {
	db                 *gorm.DB
	customerRepository repositoryInterfaces.CustomerRepository
}

func NewCustomerService(customerRepository repositoryInterfaces.CustomerRepository) interfaces.CustomerService {
	return &customerServices{
		db:                 database.MysqlDB,
		customerRepository: customerRepository,
	}
}

func (c *customerServices) Create(ctx context.Context, request *dto.CreateCustomerRequestDTO) *httpresponse.HttpError {
	logger := logrus.WithContext(ctx).WithField("request", helper.Dump(request))

	// create transaction
	tx := c.db.Begin()
	defer tx.Rollback()

	// create
	_, err := c.customerRepository.Create(ctx, tx, request.ToEntity())
	if err != nil {
		logger.Error(err)
		return httpresponse.ErrorInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		logger.Error(err)
		return httpresponse.ErrorInternalServerError
	}

	return nil
}

func (c *customerServices) GetByID(ctx context.Context, ID uint) (*dto.GetCustomerResponseDTO, *httpresponse.HttpError) {
	logger := logrus.WithContext(ctx).WithField("id", ID)

	// get from db
	customer, err := c.customerRepository.GetByID(ctx, ID)
	if err != nil {
		logger.Error(err)
		return nil, httpresponse.ErrorInternalServerError
	}

	if customer == nil {
		logger.Error(httpresponse.ErrorCustomerNotFound)
		return nil, httpresponse.ErrorCustomerNotFound
	}

	return dto.CustomerEntityToResponseDTO(customer), nil
}
