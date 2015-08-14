package api

import (
	"encoding/json"
	"errors"
	"github.com/angdev/chocolat/lib/query"
	"github.com/jinzhu/now"
	"strconv"
	"strings"
	"time"
)

var (
	InvalidTimeFrameError = errors.New("Invalid timeframe")
	InvalidGroupByError   = errors.New("Invalid group_by")
	InvalidFilterError    = errors.New("Invalid filter")
	InvalidIntervalError  = errors.New("Invalid timeframe")
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
		return InvalidTimeFrameError
	}
}

func (this *TimeFrame) IsGiven() bool {
	return !(this.Start.IsZero() && this.End.IsZero())
}

type GroupBy []string

func (g *GroupBy) UnmarshalJSON(data []byte) error {
	var multi []string
	var single string

	if err := json.Unmarshal(data, &multi); err == nil {
		*g = multi
	} else if err = json.Unmarshal(data, &single); err == nil {
		*g = []string{single}
	} else {
		return InvalidGroupByError
	}

	return nil
}

type Filter struct {
	PropertyName  string      `json:"property_name"`
	Operator      string      `json:"operator"`
	PropertyValue interface{} `json:"property_value"`
}

type Filters []Filter

type Interval string

func (i Interval) IsGiven() bool {
	return i != ""
}

func (i Interval) NextTime(t time.Time) (time.Time, error) {
	if i.isCustomInterval() {
		return i.nextCustomTime(t)
	}
	return i.nextSupportedTime(t)
}

func (i Interval) isCustomInterval() bool {
	return len(strings.Split(string(i), "_")) == 3
}

func (i Interval) nextSupportedTime(t time.Time) (time.Time, error) {
	switch i {
	case "minutely":
		return i.nextTime(t, 1, "minutes")
	case "hourly":
		return i.nextTime(t, 1, "hours")
	case "daily":
		return i.nextTime(t, 1, "days")
	case "weekly":
		return i.nextTime(t, 1, "weeks")
	case "monthly":
		return i.nextTime(t, 1, "months")
	case "yearly":
		return i.nextTime(t, 1, "years")
	}
	return time.Time{}, InvalidIntervalError
}

func (i Interval) nextCustomTime(t time.Time) (time.Time, error) {
	params := strings.Split(string(i), "_")

	if params[0] != "every" {
		return time.Time{}, InvalidIntervalError
	}

	units := params[2]

	if n, err := strconv.Atoi(params[1]); err != nil {
		return time.Time{}, err
	} else {
		return i.nextTime(t, n, units)
	}
}

func (i Interval) nextTime(t time.Time, n int, units string) (time.Time, error) {
	switch units {
	case "minutes":
		return now.New(t.Add(time.Minute * time.Duration(n))).BeginningOfMinute(), nil
	case "hours":
		return now.New(t.Add(time.Hour * time.Duration(n))).BeginningOfHour(), nil
	case "days":
		return now.New(t.AddDate(0, 0, n)).BeginningOfDay(), nil
	case "weeks":
		return now.New(t.AddDate(0, 0, 7*n)).BeginningOfDay(), nil
	case "months":
		return now.New(t.AddDate(0, n, 0)).BeginningOfDay(), nil
	case "years":
		return now.New(t.AddDate(n, 0, 0)).BeginningOfDay(), nil
	}
	return time.Time{}, InvalidIntervalError
}
