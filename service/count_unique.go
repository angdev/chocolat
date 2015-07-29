package service

import (
	"github.com/angdev/chocolat/model"
)

type CountUniqueParams struct {
	QueryParams
	TargetProperty string `json:"target_property"`
}

func CountUnique(p *model.Project, params *CountUniqueParams) (interface{}, error) {
	a := NewAggregator(p, &params.QueryParams)
	if result, err := a.CountUnique(params.TargetProperty); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
