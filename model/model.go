package model

import (
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func Init(x *gorm.DB) {
	db = x
	db.AutoMigrate(&Project{}, &ApiKey{})
}

func DB() *gorm.DB {
	return db
}
