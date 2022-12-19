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
	Email                  string         `gorm:"type:string"`
	MasterHash             []byte         `gorm:"type:bytes"`
	MasterHashSalt         []byte         `gorm:"type:bytes"`
	ProtectedDatabaseKey   []byte         `gorm:"type:bytes"`
	ProtectedDatabaseKeyIV []byte         `gorm:"type:bytes"`
	SecureEntries          []SecureEntry  `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SessionTokens          []SessionToken `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type SecureEntry struct {
	Base
	UserId           uuid.UUID         `gorm:"type:uuid"`
	Username         []byte            `gorm:"type:bytes"`
	Password         []byte            `gorm:"type:bytes"`
	AdditionalFields []AdditionalField `gorm:"foreignKey:CredentialId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type AdditionalField struct {
	Base
	CredentialId uuid.UUID `gorm:"type:uuid"`
	Field        string    `gorm:"type:string"`
	Value        string    `gorm:"type:string"`
}

type SessionToken struct {
	Base
	UserId uuid.UUID `gorm:"type:uuid"`
	N      []byte    `sql:"index" gorm:"type:bytes"`
	E      int       `gorm:"type:int"`
	Expiry time.Time `gorm:"type:time"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (baseError error) {
	base.Id = uuid.New()
	return nil
}

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&SecureEntry{})
	db.AutoMigrate(&AdditionalField{})
	db.AutoMigrate(&SessionToken{})
}
