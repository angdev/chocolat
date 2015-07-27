package service

import (
	"errors"
	"github.com/angdev/chocolat/model"
	"time"

	"github.com/angdev/chocolat/support/repo"
)

func NewCountParams(collName string, params repo.Doc) (*CountParams, error) {
	ok := false
	var err error
	var out interface{}

	var countParams CountParams
	countParams.CollectionName = collName
	if out, ok = params["timeframe"]; ok {
		if countParams.TimeFrame, err = NewTimeFrame(out); err != nil {
			return nil, err
		}
	}

	return &countParams, nil
}

type CountParams struct {
	CollectionName string
	TimeFrame      *TimeFrame
}

type TimeFrame struct {
	Start time.Time
	End   time.Time
}

func NewTimeFrame(t interface{}) (*TimeFrame, error) {
	switch t.(type) {
	case string:
		return nil, errors.New("Not implemented")
	case map[string]interface{}:
		v, err := absoluteTimeFrame(t.(map[string]interface{}))
		return v, err
	default:
		return nil, errors.New("Invalid timeframe")
	}
}

func absoluteTimeFrame(v map[string]interface{}) (*TimeFrame, error) {
	start, err := time.Parse(time.RFC3339, v["start"].(string))
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, v["end"].(string))
	if err != nil {
		return nil, err
	}

	return &TimeFrame{Start: start, End: end}, nil
}

func Count(p *model.Project, params *CountParams) (repo.Doc, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	if count, err := r.C(params.CollectionName).Count(); err != nil {
		return nil, err
	} else {
		return repo.Doc{"result": count}, nil
	}
}
