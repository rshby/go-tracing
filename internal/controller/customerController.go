package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-tracing/internal/http/httpresponse"
	serviceDto "go-tracing/internal/service/dto"
	serviceInterfaces "go-tracing/internal/service/interfaces"
	"go-tracing/internal/utils/helper"
	"go-tracing/otel"
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

	//logger := logrus.WithContext(c)

	// get from params
	id := helper.ExpectNumber[uint](c.Param("id"))

	ctx, span := otel.OtelApp.Start(c, fmt.Sprintf("v1/customer/%d", id))
	defer span.End()

	customer, httpError := controller.customerService.GetByID(ctx, id)
	if httpError != nil {
		//logger.Error(httpError)
		span.RecordError(httpError)
		httpresponse.ResponseError(c, httpError)
		return
	}

	httpresponse.ResponseOK(c, "success get data customer", customer)
	return
}
