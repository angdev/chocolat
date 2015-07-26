package model

import (
	"github.com/jinzhu/gorm"
)

type Project struct {
	gorm.Model
	UUID string
}
