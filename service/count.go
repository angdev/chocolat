package service

import (
	"github.com/angdev/chocolat/support/repo"
)

type CountParams struct {
	CollectionName string
}

func Count(dbName string, params *CountParams) (*repo.Doc, error) {
	r, err := repo.NewRepository(dbName)
	if err != nil {
		return nil, err
	}

	if count, err := r.C(params.CollectionName).Count(); err != nil {
		return nil, err
	} else {
		return &repo.Doc{"result": count}, nil
	}
}
