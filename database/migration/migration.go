package migration

import (
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Migration is function to migrate database mysql
func Migration(db *gorm.DB) {
	migrations := migrate.FileMigrationSource{
		Dir: "database/migrationSchema",
	}

	migrate.SetTable("migrations")

	dbMySql, err := db.DB()
	if err != nil {
		logrus.Error(err)
		return
	}

	n, err := migrate.Exec(dbMySql, "mysql", migrations, migrate.Up)
	if err != nil {
		logrus.Errorf("failed to migrate database mysql : %v", err)
		return
	}

	logrus.Infof("success migration %d up!", n)
}
