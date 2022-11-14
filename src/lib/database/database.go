package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
)

func DatabaseConnect() {
	var databasePath string
	if os.Getenv("ENVIRONMENT") == "testing" {
		databasePath = os.Getenv("TEST_DB_PATH")
	} else if os.Getenv("ENVIRONMENT") == "development" {
		databasePath = os.Getenv("DEV_DB_PATH")
	} else if os.Getenv("ENVIRONMENT") == "production" {
		databasePath = os.Getenv("DB_PATH")
	} else {
		panic("Invalid 'ENVIRONMENT' set.")
	}

	db, dbError := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	Database = db

	if dbError != nil {
		panic("Couldn't open database.")
	}

	MigrateModels(db)
}
