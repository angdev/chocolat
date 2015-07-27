package model

import (
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type Project struct {
	gorm.Model
	UUID    string
	ApiKeys []ApiKey
}

func (this *Project) AfterCreate(tx *gorm.DB) error {
	scopes := []ApiScope{ApiReadKey, ApiWriteKey, ApiMasterKey}

	for _, scope := range scopes {
		err := tx.Create(&ApiKey{
			Project: *this,
			Value:   uuid.NewV4().String(),
			Scope:   scope,
		}).Error

		if err != nil {
			return err
		}
	}

	return nil
}

func ProjectByUUID(uuid string) *Project {
	var project Project
	if DB().First(&project, &Project{UUID: uuid}).RecordNotFound() {
		return nil
	} else {
		return &project
	}
}

func (this *Project) apiKeys(out []ApiKey) *gorm.DB {
	return DB().Model(this).Related(&out)
}

func (this *Project) ReadKey() *ApiKey {
	var key ApiKey
	if this.apiKeys([]ApiKey{}).First(&key, &ApiKey{Scope: ApiReadKey}).RecordNotFound() {
		return nil
	} else {
		return &key
	}
}

func (this *Project) WriteKey() *ApiKey {
	var key ApiKey
	if this.apiKeys([]ApiKey{}).First(&key, &ApiKey{Scope: ApiWriteKey}).RecordNotFound() {
		return nil
	} else {
		return &key
	}
}

func (this *Project) MasterKey() *ApiKey {
	var key ApiKey
	if this.apiKeys([]ApiKey{}).First(&key, &ApiKey{Scope: ApiMasterKey}).RecordNotFound() {
		return nil
	} else {
		return &key
	}
}
