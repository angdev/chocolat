package service

import (
	"errors"
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/support/repo"
	"labix.org/v2/mgo"
)

type AggregateResult struct {
	Result interface{} `json:"result"`
}

type IntervalResult struct {
	Value     interface{} `json:"value"`
	TimeFrame *TimeFrame  `json:"timeframe"`
}

type Aggregator struct {
	project *model.Project
	params  *QueryParams
}

func NewAggregator(p *model.Project, params *QueryParams) *Aggregator {
	return &Aggregator{project: p, params: params}
}

func (this *Aggregator) Count() (interface{}, error) {
	r := repo.NewRepository(this.project.RepoName())
	defer r.Close()

	c := r.C(this.params.CollectionName)
	pipeline := countPipeline(c, this.params)

	if this.params.Interval.IsGiven() {
		result, err := intervalAggregate(pipeline, this.params)
		return result, err
	}

	return aggregate(pipeline, this.params)
}

func aggregate(p *Pipeline, params *QueryParams) (*AggregateResult, error) {
	result, err := p.Result()
	return &AggregateResult{Result: result}, err
}

func intervalAggregate(p *Pipeline, params *QueryParams) (*AggregateResult, error) {
	if !params.TimeFrame.IsGiven() {
		return nil, errors.New("timeframe is not given")
	}

	var results []IntervalResult

	start := params.TimeFrame.Start
	frameEnd := params.TimeFrame.End

	for start.Before(frameEnd) {
		end := params.Interval.NextTime(start)
		t := TimeFrame{Start: start, End: end}
		if result, err := p.Copy().Prepend(t.Pipe()).Result(); err != nil {
			return nil, err
		} else {
			results = append(results, IntervalResult{Value: result, TimeFrame: &t})
			start = end
		}
	}

	return &AggregateResult{Result: results}, nil
}

func countPipeline(c *mgo.Collection, params *QueryParams) *Pipeline {
	pipeline := NewPipeline(c)
	pipeline.Append(params.Filters.Pipe())
	pipeline.Append(params.GroupBy.Pipe(repo.Doc{
		"result": repo.Doc{"$sum": 1},
	}))

	return pipeline
}
