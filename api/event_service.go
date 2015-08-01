package api

import (
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/support/repo"
	"time"
)

type CreateEventParams struct {
	CollectionName string
	Events         map[string][]repo.Doc
}

func createEvent(p *model.Project, params *CreateEventParams) (repo.Doc, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	docs := params.Events[params.CollectionName]

	result := insertEvents(r, params.CollectionName, docs...)
	if result[0]["success"] == true {
		return repo.Doc{"created": true}, nil
	} else {
		return repo.Doc{"created": false}, nil
	}
}

func createMultipleEvents(p *model.Project, params *CreateEventParams) (map[string][]repo.Doc, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	result := map[string][]repo.Doc{}
	for event, docs := range params.Events {
		result[event] = insertEvents(r, event, docs...)
	}

	return result, nil
}

func insertEvents(r *repo.Repository, event string, docs ...repo.Doc) []repo.Doc {
	result := []repo.Doc{}
	for _, doc := range docs {
		appendMetadata(doc)
		if err := r.Insert(event, &doc); err != nil {
			result = append(result, repo.Doc{"success": false})
		} else {
			result = append(result, repo.Doc{"success": true})
		}
	}

	return result
}

func appendMetadata(doc repo.Doc) {
	metadata := map[string]interface{}{}
	metadata["created_at"] = time.Now()

	doc["chocolat"] = metadata
}
