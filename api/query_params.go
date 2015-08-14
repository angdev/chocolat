package api

import (
	"encoding/json"
	"errors"
	"github.com/angdev/chocolat/lib/query"
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

func (this *QueryParams) ToQuery() *query.Arel {
	var conds []*query.Condition
	for _, f := range this.Filters {
		conds = append(conds, query.NewCondition(f.PropertyName, f.Operator, f.PropertyValue))
	}

	t := this.TimeFrame
	if t.IsGiven() {
		conds = append(conds,
			query.NewCondition("chocolat.created_at", "gt", t.Start),
			query.NewCondition("chocolat.created_at", "lt", t.End))
	}

	return query.NewArel().Where(conds...).GroupBy(this.GroupBy...)
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

type Filter struct {
	PropertyName  string      `json:"property_name"`
	Operator      string      `json:"operator"`
	PropertyValue interface{} `json:"property_value"`
}

type FilterError struct{}

func (FilterError) Error() string {
	return "Invalid filter"
}

type Filters []Filter

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
