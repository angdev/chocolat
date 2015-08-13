package model

import (
	"database/sql/driver"
	"errors"
	"github.com/jinzhu/gorm"
)

type ApiKey struct {
	gorm.Model
	Project   Project
	ProjectID int64
	Value     string
	Scope     ApiScope
	Revoked   bool
}

func (this *ApiKey) TableName() string {
	return "api_keys"
}

func ApiKeyByValue(value string) *ApiKey {
	var key ApiKey
	if db.First(&key, &ApiKey{Value: value}).RecordNotFound() {
		return nil
	} else {
		return &key
	}
}

type ApiScope int64

const (
	ApiReadKey ApiScope = iota
	ApiWriteKey
	ApiMasterKey
)

func (this *ApiScope) Scan(value interface{}) error {
	if v, ok := value.(int64); ok {
		*this = ApiScope(v)
		return nil
	} else {
		return errors.New("Cannot convert to ApiScope")
	}
}

func (this ApiScope) Value() (driver.Value, error) {
	return int64(this), nil
}
