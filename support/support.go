package support

import (
	"encoding/base64"
	"encoding/json"
	"github.com/angdev/chocolat/support/repo"
)

func DecodeData(data string) (repo.Doc, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	var unmarshalled repo.Doc
	if err := json.Unmarshal(decoded, &unmarshalled); err != nil {
		return nil, err
	}

	return unmarshalled, nil
}
