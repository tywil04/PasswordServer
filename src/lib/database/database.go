package database

import (
	"os"
	customErrors "passwordserver/src/lib/cerrors"

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
		panic(customErrors.ErrorEnvironmentEnvNotFound)
	}

	db, dbError := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	Database = db

	if dbError != nil {
		panic(customErrors.ErrorLoadingDatabase)
	}

	MigrateModels(db)
}
