package lib

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
)

type Base struct {
	Id        uuid.UUID  `gorm:"type:uuid"`
	CreatedAt time.Time  `gorm:"type:time"`
	UpdatedAt time.Time  `gorm:"type:time"`
	DeletedAt *time.Time `sql:"index"`
}

type User struct {
	Base
	MasterHash           []byte       `gorm:"type:bytes"`
	MasterHashSalt       []byte       `gorm:"type:bytes"`
	ProtectedDatabaseKey []byte       `gorm:"type:bytes"`
	Credentials          []Credential `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tokens               [][]byte     `gorm:"type:bytes"`
}

type Credential struct {
	Base
	UserId   uuid.UUID `gorm:"type:uuid"`
	Username []byte    `gorm:"type:bytes"`
	Password []byte    `gorm:"type:bytes"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.Id = uuid.New()
	return
}

func DatabaseConnect() {
	db, dbError := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	Database = db

	if dbError != nil {
		panic("Couldn't open database.")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Credential{})
}
