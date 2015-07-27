package service

import (
	"github.com/angdev/chocolat/support/repo"
	"strings"
	"time"
)

type QueryParams struct {
	CollectionName string
	TimeFrame      *TimeFrame
	GroupBy        GroupBy
}

func NewQueryParams(collName string, params repo.Doc) (*QueryParams, error) {
	ok := false
	var err error
	var out interface{}

	var qp QueryParams
	qp.CollectionName = collName

	if out, ok = params["timeframe"]; ok {
		if qp.TimeFrame, err = NewTimeFrame(out); err != nil {
			return nil, err
		}
	}

	if out, ok = params["group_by"]; ok {
		if qp.GroupBy, err = NewGroupBy(out); err != nil {
			return nil, err
		}
	}

	return &qp, nil
}

type TimeFrame struct {
	Start time.Time
	End   time.Time
}

type TimeFrameError struct{}

func (TimeFrameError) Error() string {
	return "Invalid timeframe"
}

func NewTimeFrame(t interface{}) (*TimeFrame, error) {
	switch t.(type) {
	case string:
		return nil, TimeFrameError{}
	case map[string]interface{}:
		v, err := absoluteTimeFrame(t.(map[string]interface{}))
		return v, err
	default:
		return nil, TimeFrameError{}
	}
}

func absoluteTimeFrame(v map[string]interface{}) (*TimeFrame, error) {
	s, ok := v["start"].(string)
	if !ok {
		return nil, TimeFrameError{}
	}

	e, ok := v["end"].(string)
	if !ok {
		return nil, TimeFrameError{}
	}

	start, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, e)
	if err != nil {
		return nil, err
	}

	return &TimeFrame{Start: start, End: end}, nil
}

// Create a mongo aggregation pipeline format
func (this *TimeFrame) Pipe() repo.Doc {
	return repo.Doc{
		"$match": repo.Doc{
			"chocolat.created_at": repo.Doc{"$gt": this.Start, "$lt": this.End},
		},
	}
}

type GroupBy []string

type GroupByError struct{}

func (GroupByError) Error() string {
	return "Invalid group_by"
}

func NewGroupBy(g interface{}) (GroupBy, error) {
	switch g.(type) {
	case string:
		return GroupBy{g.(string)}, nil
	case []interface{}:
		gs := g.([]interface{})
		groups := make([]string, len(gs))
		for i, v := range gs {
			groups[i] = v.(string)
		}
		return GroupBy(groups), nil
	default:
		return nil, GroupByError{}
	}
}

func (this GroupBy) Pipe(ops ...repo.Doc) repo.Doc {
	group := repo.Doc{
		"_id": expandFields(this),
	}

	for _, op := range ops {
		for k, v := range op {
			group[k] = v
		}
	}

	return repo.Doc{
		"$group": group,
	}
}

func expandFields(fields []string) repo.Doc {
	if len(fields) == 0 {
		return nil
	}

	expanded := repo.Doc{}

	for _, field := range fields {
		keys := strings.Split(field, ".")
		deepAssign(expanded, variablize(field), keys...)
	}

	return expanded
}

func collapseField(doc repo.Doc) repo.Doc {
	collapsed := repo.Doc{}

	var f func([]string, repo.Doc)
	f = func(level []string, cursor repo.Doc) {
		for k, v := range cursor {
			switch v.(type) {
			case repo.Doc:
				f(append(level, k), v.(repo.Doc))
			default:
				collapsed[strings.Join(append(level, k), ".")] = v
			}
		}
	}

	f([]string{}, doc)

	return collapsed
}

func variablize(fields ...string) string {
	if len(fields) == 0 {
		return ""
	} else {
		return "$" + strings.Join(fields, ".")
	}
}

func deepAssign(d repo.Doc, value interface{}, keys ...string) {
	cursor := d
	midKeys, lastKey := keys[:len(keys)-1], keys[len(keys)-1]
	for _, key := range midKeys {
		if _, ok := cursor[key]; !ok {
			cursor[key] = repo.Doc{}
		}
		cursor = cursor[key].(repo.Doc)
	}
	cursor[lastKey] = value
}
