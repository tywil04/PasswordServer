package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	Id        uuid.UUID      `gorm:"type:uuid"`
	CreatedAt time.Time      `gorm:"type:time"`
	UpdatedAt time.Time      `gorm:"type:time"`
	DeletedAt gorm.DeletedAt `sql:"index"`
}

type User struct {
	Base
	Email                string         `gorm:"type:string"`
	MasterHash           []byte         `gorm:"type:bytes"`
	MasterHashSalt       []byte         `gorm:"type:bytes"`
	ProtectedDatabaseKey []byte         `gorm:"type:bytes"`
	Credentials          []Credential   `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SessionTokens        []SessionToken `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Credential struct {
	Base
	UserId   uuid.UUID `gorm:"type:uuid"`
	Username []byte    `gorm:"type:bytes"`
	Password []byte    `gorm:"type:bytes"`
}

type SessionToken struct {
	Base
	UserId uuid.UUID `gorm:"type:uuid"`
	N      []byte    `sql:"index" gorm:"type:bytes"`
	E      int       `gorm:"type:int"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (baseError error) {
	base.Id = uuid.New()
	return
}

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Credential{})
	db.AutoMigrate(&SessionToken{})
}
