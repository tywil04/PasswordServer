package database

import (
	"os"

	psErrors "passwordserver/src/lib/errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
)

func DatabaseConnect() {
	var databasePath string
	if os.Getenv("ENVIRONMENT") == "testing" {
		databasePath = os.Getenv("TESTING_DB_PATH")
	} else if os.Getenv("ENVIRONMENT") == "development" {
		databasePath = os.Getenv("DEVELOPMENT_DB_PATH")
	} else if os.Getenv("ENVIRONMENT") == "production" {
		databasePath = os.Getenv("DB_PATH")
	} else {
		panic(psErrors.ErrorEnvironmentEnvNotFound)
	}

	db, dbError := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if dbError != nil {
		panic(psErrors.ErrorLoadingDatabase)
	}

	Database = db

	MigrateModels(db)
}
