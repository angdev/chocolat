package api

import (
	"encoding/json"
	"strconv"
)

type StringInt int

func (this *StringInt) UnmarshalJSON(data []byte) error {
	var parsed interface{}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	switch parsed.(type) {
	case int:
		*this = StringInt(parsed.(int))
	case string:
		if i, err := strconv.Atoi(parsed.(string)); err != nil {
			return err
		} else {
			*this = StringInt(i)
		}
	default:
		*this = 0
	}

	return nil
}
