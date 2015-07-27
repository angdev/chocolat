package service

import (
	"github.com/angdev/chocolat/support/repo"
	"time"
)

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
	stage := repo.Doc{
		"$match": repo.Doc{
			"chocolat.created_at": repo.Doc{"$gt": this.Start, "$lt": this.End},
		},
	}

	return stage
}
