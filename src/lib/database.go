package lib

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Base struct {
	Id        uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type User struct {
	Base
	MasterHash           []byte
	ProtectedDatabaseKey []byte
	Credentials          []Credential `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tokens               [][]byte
}

type Credential struct {
	Base
	UserId   uuid.UUID
	Username []byte
	Password []byte
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.Id = uuid.New()
	return
}

func DatabaseConnect() {
	db, dbError := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	if dbError != nil {
		panic("Couldn't open database.")
	}

	db.AutoMigrate(&Base{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Credential{})
}
