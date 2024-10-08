package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-tracing/internal/http/httpresponse"
	serviceDto "go-tracing/internal/service/dto"
	serviceInterfaces "go-tracing/internal/service/interfaces"
	"go-tracing/internal/utils/helper"
	"go-tracing/otel"
	"go.opentelemetry.io/otel/codes"
)

type CustomerController struct {
	customerService serviceInterfaces.CustomerService
}

func NewCustomerController(customerService serviceInterfaces.CustomerService) *CustomerController {
	return &CustomerController{customerService: customerService}
}

func (controller *CustomerController) Create(c *gin.Context) {
	logger := logrus.WithContext(c)

	// decode
	var request serviceDto.CreateCustomerRequestDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error(err)
		httpresponse.ResponseBadRequest(c)
		return
	}

	if httpErr := controller.customerService.Create(c.Request.Context(), &request); httpErr != nil {
		logger.Error(httpErr)
		httpresponse.ResponseError(c, httpErr)
		return
	}

	// ok
	httpresponse.ResponseOK(c, "success create customer", nil)
	return
}

func (controller *CustomerController) GetByID(c *gin.Context) {
	ctx, span := otel.OtelApp.Start(c.Request.Context(), helper.MyCaller(1))
	defer span.End()
	//logger := logrus.WithContext(c)

	// get from params
	id := helper.ExpectNumber[uint](c.Param("id"))

	customer, httpError := controller.customerService.GetByID(ctx, id)
	if httpError != nil {
		span.SetStatus(codes.Error, httpError.Message)
		span.RecordError(httpError)
		httpresponse.ResponseError(c, httpError)
		return
	}

	httpresponse.ResponseOK(c, "success get data customer", customer)
	return
}
