package models

import (
	"github.com/anyingiit/UnifyAuthCenter/db"
	"gorm.io/gorm"
)

type Sessions []Session

func (s *Sessions) Find() (result *gorm.DB) {
	return db.Db.Find(s)
}
