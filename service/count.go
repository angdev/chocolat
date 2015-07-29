package service

import (
	"github.com/angdev/chocolat/model"
	"github.com/k0kubun/pp"
)

type CountParams struct {
	QueryParams
}

func Count(p *model.Project, params *CountParams) (interface{}, error) {
	a := NewAggregator(p, &params.QueryParams)
	if result, err := a.Count(); err != nil {
		return nil, err
	} else {
		pp.Println(result)
		return result, nil
	}
}
