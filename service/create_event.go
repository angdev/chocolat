package service

import (
	"github.com/angdev/chocolat/support"
	"github.com/angdev/chocolat/support/repo"
)

type CreateEventParams struct {
	CollectionName string
	EncodedData    string
}

func CreateEvent(dbName string, params *CreateEventParams) (*repo.Doc, error) {
	r, err := repo.NewRepository(dbName)
	if err != nil {
		return nil, err
	}

	var event *repo.Doc
	if params.EncodedData != "" {
		event, err = support.DecodeData(params.EncodedData)
		if err != nil {
			return nil, err
		}
	}

	if err := r.Insert(params.CollectionName, event); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}
