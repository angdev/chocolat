package model

import (
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func Init(x *gorm.DB) {
	db = x
}

func DB() *gorm.DB {
	return db
}
