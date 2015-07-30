package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/angdev/chocolat/support/repo"
	"github.com/jinzhu/now"
	"time"
)

type QueryParams struct {
	CollectionName string    `json:"event_collection"`
	TimeFrame      TimeFrame `json:"timeframe"`
	GroupBy        GroupBy   `json:"group_by"`
	Filters        Filters   `json:"filters"`
	Interval       Interval  `json:"interval"`
}

type TimeFrame struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	absolute bool
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
		t.absolute = true
		return nil
	} else if err = json.Unmarshal(data, &rel); err == nil {
		return errors.New("Relative timeframe is not implemented")
	} else {
		return TimeFrameError{}
	}
}

func (this *TimeFrame) IsGiven() bool {
	return !(this.Start.IsZero() && this.End.IsZero())
}

// Create a mongo aggregation pipeline format
func (this *TimeFrame) Pipe(...repo.Doc) Pipe {
	if !this.IsGiven() {
		return nil
	}

	return NewPipe(PipeStage{
		"$match": repo.Doc{
			"chocolat.created_at": repo.Doc{"$gt": this.Start, "$lt": this.End},
		},
	})
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

func (this GroupBy) Pipe(ops ...repo.Doc) Pipe {
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

	project := repo.Doc{}
	for _, field := range this {
		project[field] = variablize("_id", field)
	}
	project["_id"] = false
	project["result"] = variablize("result")

	return NewPipe(PipeStage{
		"$group": group,
	}, PipeStage{
		"$project": project,
	})
}

type Filter struct {
	PropertyName  string      `json:"property_name"`
	Operator      string      `json:"operator"`
	PropertyValue interface{} `json:"property_value"`
}

type FilterError struct{}

func (f *Filter) QueryOp() (string, interface{}) {
	var op string
	var value interface{}

	switch f.Operator {
	case "contains":
		op = "$regex"
		value = fmt.Sprintf("/%s/", f.PropertyValue.(string))
	case "not_contains":
		op = "$regex"
		value = fmt.Sprintf("/^(%s)/", f.PropertyValue.(string))
	default:
		op = fmt.Sprintf("$%s", f.Operator)
		value = f.PropertyValue
	}

	return op, value
}

func (FilterError) Error() string {
	return "Invalid filter"
}

type Filters []Filter

func (f Filters) Pipe(...repo.Doc) Pipe {
	if len(f) == 0 {
		return nil
	}

	match := map[string]interface{}{}

	for _, filter := range f {
		op, value := filter.QueryOp()
		deepAssign(match, value, filter.PropertyName, op)
	}

	return NewPipe(PipeStage{
		"$match": match,
	})
}

type Interval string

func (i *Interval) IsGiven() bool {
	return (*i) != ""
}

func (i *Interval) NextTime(t time.Time) time.Time {
	switch *i {
	case "daily":
		return now.New(t.AddDate(0, 0, 1)).BeginningOfDay()
	}
	// temp
	return t
}
