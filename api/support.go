package api

import (
	"encoding/base64"
	"encoding/json"
)

func decodeData(data string, out interface{}) error {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(decoded, &out); err != nil {
		return err
	}

	return nil
}
