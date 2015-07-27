package service

import (
	"github.com/angdev/chocolat/model"
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

func Count(p *model.Project, params *CountParams) (repo.Doc, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	pipes := []repo.Doc{}
	if params.TimeFrame != nil {
		pipes = append(pipes, params.TimeFrame.Pipe())
	}
	pipes = append(pipes, countPipe())

	pipe := r.C(params.CollectionName).Pipe(pipes)
	iter := pipe.Iter()

	result := []repo.Doc{}
	if err := iter.All(&result); err != nil {
		return nil, err
	}

	if len(result) != 0 {
		return repo.Doc{"result": result[0]["count"]}, nil
	} else {
		return repo.Doc{"result": 0}, nil
	}
}

func countPipe() repo.Doc {
	return repo.Doc{
		"$group": repo.Doc{
			"_id":   nil,
			"count": repo.Doc{"$sum": 1},
		},
	}
}
