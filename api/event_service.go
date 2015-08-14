package api

import (
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/repo"
	"time"
)

type CreateSingleEventParams RawResult
type CreateMultipleEventParams map[string][]CreateSingleEventParams

type CreateEventParams struct {
	CollectionName string
	Events         map[string][]interface{}
}

func createEvent(p *model.Project, collName string, params CreateSingleEventParams) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	result := insertEvents(r, collName, params)
	if result[0]["success"] == true {
		return RawResult{"created": true}, nil
	} else {
		return RawResult{"created": false}, nil
	}
}

func createMultipleEvents(p *model.Project, params CreateMultipleEventParams) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	result := make(map[string][]RawResult)
	for event, docs := range params {
		result[event] = insertEvents(r, event, docs...)
	}

	return result, nil
}

func insertEvents(r *repo.Repository, event string, docs ...CreateSingleEventParams) []RawResult {
	result := []RawResult{}
	for _, doc := range docs {
		appendMetadata(doc)
		if err := r.Insert(event, &doc); err != nil {
			result = append(result, RawResult{"success": false})
		} else {
			result = append(result, RawResult{"success": true})
		}
	}

	return result
}

func appendMetadata(doc CreateSingleEventParams) {
	metadata := map[string]interface{}{}
	metadata["created_at"] = time.Now()

	doc["chocolat"] = metadata
}
