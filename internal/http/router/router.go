package router

import (
	"github.com/gin-gonic/gin"
	"go-tracing/internal/controller"
	"go-tracing/internal/repository"
	"go-tracing/internal/service"
)

type router struct {
	app *gin.RouterGroup
}

func NewRouter(app *gin.RouterGroup) {
	r := router{app: app}
	r.SetEnpoint()
}

func (r *router) SetEnpoint() {
	customerController := controller.NewCustomerController(service.NewCustomerService(repository.NewCustomerRepository()))

	apiV1 := r.app.Group("/v1")
	{
		customerGroup := apiV1.Group("/customer")
		{
			customerGroup.GET("/:id", customerController.GetByID)
		}
	}
}
