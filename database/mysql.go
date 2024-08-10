package database

import (
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go-tracing/database/migration"
	"go-tracing/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	MysqlDB *gorm.DB
)

func InitializeMysqlDatabase() (*gorm.DB, func()) {
	db, err := gorm.Open(mysql.Open(config.DatabaseDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		logrus.Fatalf("failed to connect database : %s", err.Error())
	}

	MysqlDB = db
	MysqlDB.Use(otelgorm.NewPlugin())

	migration.Migration(MysqlDB)
	return MysqlDB, func() {
		s, _ := db.DB()
		_ = s.Close()
		logrus.Info("close database mysql")
	}
}

func CloseDB(mysql *gorm.DB) {
	if mysql != nil {
		db, err := mysql.DB()
		if err != nil {
			logrus.Error(err)
		}

		_ = db.Close()
	}
}
