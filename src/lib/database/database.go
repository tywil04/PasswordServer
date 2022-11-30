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

func DatabaseConnect() *gorm.DB {
	var databasePath string
	switch os.Getenv("ENVIRONMENT") {
	case "testing":
		databasePath = os.Getenv("TESTING_DB_PATH")
	case "development":
		databasePath = os.Getenv("DEVELOPMENT_DB_PATH")
	case "production":
		databasePath = os.Getenv("DB_PATH")
	default:
		panic(psErrors.ErrorEnvironmentEnvNotFound)
	}

	db, dbError := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if dbError != nil {
		panic(psErrors.ErrorLoadingDatabase)
	}

	Database = db

	MigrateModels(db)

	return Database
}
