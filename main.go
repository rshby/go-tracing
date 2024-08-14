package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go-tracing/database"
	"go-tracing/internal/config"
	"go-tracing/internal/http/middleware"
	"go-tracing/internal/http/router"
	"go-tracing/internal/logger"
	"go-tracing/otel"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	logger.SetupLogger()

}

func main() {
	_, closerTracer := otel.InitTracerApp(context.Background(), "go-tracing")
	//otel.NewMetrixPrometheus(context.Background(), "go-tracing")
	premetheusCloser := otel.InitiaizeMetricWithOtelPremetheus(context.Background(), "go-tracing")
	defer premetheusCloser()

	mysqlDB, mysqlCloser := database.InitializeMysqlDatabase()
	defer mysqlCloser()
	logrus.Info(mysqlDB)

	app := gin.Default()
	app.Use(middleware.TraceMiddleware())

	app.GET("/metrics", PromHandler())

	// router
	router.NewRouter(&app.RouterGroup)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppPort()),
		Handler: app,
	}

	var (
		chanExit   = make(chan bool)
		chanSignal = make(chan os.Signal)
		chanErr    = make(chan error)
	)

	ctx, span := otel.OtelApp.Start(context.Background(), "")
	defer span.End()

	logrus.Info(ctx)

	signal.Notify(chanSignal, os.Interrupt)

	go func() {
		for {
			select {
			case <-chanSignal:
				logrus.WithContext(ctx).Info("chan signal receive")
				timeoutCtx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)

				_ = server.Shutdown(timeoutCtx)
				closerTracer()
				// close sql
				database.CloseDB(database.MysqlDB)
				cancelFunc()
				chanExit <- true
				return
			case e := <-chanErr:
				logrus.WithContext(ctx).Info("error server receive")
				logrus.Error(e)
				_ = server.Close()

				timeoutCtx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
				_ = server.Shutdown(timeoutCtx)
				closerTracer()

				// close sql
				database.CloseDB(database.MysqlDB)

				cancelFunc()
				chanExit <- true
				return
			}
		}
	}()

	go func() {
		logrus.Infof("running on port %d", config.AppPort())
		if err := server.ListenAndServe(); err != nil {
			chanErr <- err
			return
		}
	}()

	<-chanExit
	logrus.Info("server exit🔴")
}

func PromHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
