package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-tracing/database"
	"go-tracing/internal/config"
	"go-tracing/internal/http/router"
	"go-tracing/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

func init() {
	logger.SetupLogger()
}

func main() {
	mysqlDB, mysqlCloser := database.InitializeMysqlDatabase()
	defer mysqlCloser()
	logrus.Info(mysqlDB)

	app := gin.Default()

	// router
	router.NewRouter(&app.RouterGroup)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppPort()),
		Handler: app,
	}

	var (
		wg         = &sync.WaitGroup{}
		chanSignal = make(chan os.Signal)
	)

	signal.Notify(chanSignal, os.Interrupt)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-chanSignal:
				_ = server.Close()
				return
			}
		}
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		logrus.Infof("running on port %d", config.AppPort())
		if err := server.ListenAndServe(); err != nil {
			logrus.Error(err)
			return
		}
	}(wg)

	wg.Wait()
}
