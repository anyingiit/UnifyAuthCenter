package models

import (
	"github.com/anyingiit/UnifyAuthCenter/db"
	"gorm.io/gorm"
)

// 注意: 目前角色都是手动创建的
const (
	RoleAdminId = iota + 1
	RoleInternalId
	RoleUserId
)

type Role struct {
	ID          int    `gorm:"primarykey;unique;not null;type:int"`
	Description string `gorm:"unique;not null;type:string"`
}

func (r *Role) Create() (result *gorm.DB) {
	return db.Db.Create(r)
}

func (r *Role) First() (result *gorm.DB) {
	return db.Db.First(r)
}

func (r *Role) Delete() (result *gorm.DB) {
	return db.Db.Delete(r)
}
