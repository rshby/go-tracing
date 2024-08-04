package main

import (
	"github.com/sirupsen/logrus"
	"go-tracing/database"
	"go-tracing/internal/logger"
)

func init() {
	logger.SetupLogger()
}

func main() {
	mysqlDB, mysqlCloser := database.InitializeMysqlDatabase()
	defer mysqlCloser()
	logrus.Info(mysqlDB)
}
