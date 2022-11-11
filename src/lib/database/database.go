package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
)

func DatabaseConnect() {
	db, dbError := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	Database = db

	if dbError != nil {
		panic("Couldn't open database.")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Credential{})
	db.AutoMigrate(&SessionToken{})
}
