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

func countPipeline(c *mgo.Collection, params *QueryParams) *Pipeline {
	pipeline := NewPipeline(c)
	pipeline.Append(params.Filters.Pipe())
	pipeline.Append(params.GroupBy.Pipe(repo.Doc{
		"result": repo.Doc{"$sum": 1},
	}))

	return pipeline
}

func (this *Aggregator) CountUnique(target string) (interface{}, error) {
	r := repo.NewRepository(this.project.RepoName())
	defer r.Close()

	c := r.C(this.params.CollectionName)
	pipeline := countUniquePipeline(c, this.params, target)

	if this.params.Interval.IsGiven() {
		result, err := intervalAggregate(pipeline, this.params)
		return result, err
	}

	return aggregate(pipeline, this.params)
}

func countUniquePipeline(c *mgo.Collection, params *QueryParams, target string) *Pipeline {
	pipeline := NewPipeline(c)
	pipeline.Append(params.Filters.Pipe())

	targetGroup := append(params.GroupBy, target)
	pipeline.Append(targetGroup.Pipe())
	pipeline.Append(params.GroupBy.Pipe(repo.Doc{
		"result": repo.Doc{"$sum": 1},
	}))

	return pipeline
}

func (this *Aggregator) Min(target string) (interface{}, error) {
	// $group Op = $min
}

func (this *Aggregator) Max(target string) (interface{}, error) {
	// $group Op = $max
}

func (this *Aggregator) Sum(target string) (interface{}, error) {
	// $group Op = $sum
}

func (this *Aggregator) Average(target string) (interface{}, error) {
	// $group Op = $avg
}

func (this *Aggregator) Percentile(target string, percent int) (interface{}, error) {
	// 0. Get Count results
	// 1. For each result, Timeframe, GroupBy, Filter -> to one filter
	// 2. Match {1} -> Skip Percentile(Count, #percent) -> Limit 1
}

func (this *Aggregator) Median(target string) (interface{}, error) {
	// return Percentile(target, 50)
}

func (this *Aggregator) SelectUnique(target string) (interface{}, error) {
	// Unique w/ no counting
}
