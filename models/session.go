package models

import (
	"time"

	"github.com/anyingiit/UnifyAuthCenter/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	UUID      uuid.UUID `gorm:"primarykey;unique;not null;type:string"`
	RoleId    int       `gorm:"not null;type:int"`
	Role      Role      `gorm:"foreignKey:RoleId"`
	CreatedAt time.Time
	ExpiredAt time.Time
}

func (s *Session) Create() (result *gorm.DB) {
	return db.Db.Create(s)
}

func (s *Session) First() (result *gorm.DB) {
	return db.Db.First(s)
}

func (s *Session) Delete() (result *gorm.DB) {
	return db.Db.Delete(s)
}
