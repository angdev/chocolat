package api

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
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

func decodeUrlEncodedJSON(data string, out interface{}) error {
	decoded, err := url.QueryUnescape(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(decoded), &out); err != nil {
		return err
	}

	return nil
}
