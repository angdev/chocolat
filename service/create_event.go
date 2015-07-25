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
	r := repo.NewRepository(dbName)
	defer r.Close()

	var event *repo.Doc
	var err error

	if params.EncodedData != "" {
		event, err = support.DecodeData(params.EncodedData)
		if err != nil {
			return nil, err
		}
	}

	if err = r.Insert(params.CollectionName, event); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}
