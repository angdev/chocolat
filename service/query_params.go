package service

import (
	"encoding/json"
	"errors"
	"github.com/angdev/chocolat/support/repo"
	"strings"
	"time"
)

type QueryParams struct {
	CollectionName string     `json:"event_collection"`
	TimeFrame      *TimeFrame `json:"timeframe"`
	GroupBy        *GroupBy   `json:"group_by"`
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
