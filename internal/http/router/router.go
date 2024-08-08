package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-tracing/internal/controller"
	"go-tracing/internal/repository"
	"go-tracing/internal/service"
	"go-tracing/otel"
	"net/http"
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
			customerGroup.GET("/ok", func(c *gin.Context) {
				ctx, span := otel.OtelApp.Start(c.Request.Context(), "controller ok")
				defer span.End()

				logrus.Info(ctx)

				c.Status(http.StatusOK)
				c.JSON(http.StatusOK, gin.H{
					"status_code": http.StatusOK,
					"message":     "success cek",
				})
			})
		}
	}
}
