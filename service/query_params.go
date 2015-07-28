package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/angdev/chocolat/support/repo"
	"strings"
	"time"
)

type QueryParams struct {
	CollectionName string    `json:"event_collection"`
	TimeFrame      TimeFrame `json:"timeframe"`
	GroupBy        GroupBy   `json:"group_by"`
	Filters        Filters   `json:"filters"`
}

type TimeFrame struct {
	Start    time.Time
	End      time.Time
	Absolute bool
}

type TimeFrameError struct{}

func (TimeFrameError) Error() string {
	return "Invalid timeframe"
}

func (t *TimeFrame) UnmarshalJSON(data []byte) error {
	var abs struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	var rel string

	if err := json.Unmarshal(data, &abs); err == nil {
		t.Start = abs.Start
		t.End = abs.End
		t.Absolute = true
		return nil
	} else if err = json.Unmarshal(data, &rel); err == nil {
		return errors.New("Relative timeframe is not implemented")
	} else {
		return TimeFrameError{}
	}
}

// Create a mongo aggregation pipeline format
func (this *TimeFrame) Pipe(...repo.Doc) repo.Doc {
	if this.Start.IsZero() && this.End.IsZero() {
		return nil
	}

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

func (g *GroupBy) UnmarshalJSON(data []byte) error {
	var multi []string
	var single string

	if err := json.Unmarshal(data, &multi); err == nil {
		*g = multi
	} else if err = json.Unmarshal(data, &single); err == nil {
		*g = []string{single}
	} else {
		return GroupByError{}
	}

	return nil
}

func (this GroupBy) Pipe(ops ...repo.Doc) repo.Doc {
	group := repo.Doc{}

	if len(this) == 0 {
		group["_id"] = nil
	} else {
		criteria := repo.Doc{}
		for _, field := range this {
			criteria[field] = variablize(field)
		}
		group["_id"] = expandField(criteria)
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

type Filter struct {
	PropertyName  string      `json:"property_name"`
	Operator      string      `json:"operator"`
	PropertyValue interface{} `json:"property_value"`
}

type FilterError struct{}

func (f *Filter) QueryOp() repo.Doc {
	op := repo.Doc{}

	switch f.Operator {
	case "contains":
		op["$regex"] = fmt.Sprintf("/%s/", f.PropertyValue.(string))
	case "not_contains":
		op["$regex"] = fmt.Sprintf("/^(%s)/", f.PropertyValue.(string))
	default:
		op[fmt.Sprintf("$%s", f.Operator)] = f.PropertyValue
	}

	return op
}

func (FilterError) Error() string {
	return "Invalid filter"
}

type Filters []Filter

func (f Filters) Pipe(...repo.Doc) repo.Doc {
	if len(f) == 0 {
		return nil
	}

	match := map[string]interface{}{}

	for _, filter := range f {
		op := filter.QueryOp()
		match[filter.PropertyName] = op
	}

	return repo.Doc{
		"$match": expandField(match),
	}
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

func expandField(doc repo.Doc) repo.Doc {
	expanded := repo.Doc{}
	collapsed := collapseField(doc)

	for k, v := range collapsed {
		keys := strings.Split(k, ".")
		deepAssign(expanded, v, keys...)
	}

	return expanded
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
