package api

import (
	"fmt"
	"github.com/angdev/chocolat/lib/query"
	"github.com/angdev/chocolat/model"
	"github.com/angdev/chocolat/repo"
	"reflect"
	"time"
)

type CreateSingleEventParams RawResult
type CreateMultipleEventParams map[string][]CreateSingleEventParams

func docToSchema(doc *query.RawExpr) interface{} {
	delete(*doc, "_id")
	delete(*doc, "chocolat")

	collapsed := make(query.RawExpr)
	query.CollapseField(*doc, &collapsed)

	for k, v := range collapsed {
		switch v.(type) {
		case string:
			collapsed[k] = "string"
		case bool:
			collapsed[k] = "bool"
		case int, int16, int32, int64, int8, float32, float64, uint, uint16, uint32, uint64, uint8:
			collapsed[k] = "number"
		default:
			collapsed[k] = reflect.ValueOf(v).Type().String()
		}
	}

	return collapsed
}

func inspectCollection(p *model.Project, name string) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	c := r.C(name)
	var result query.RawExpr

	if err := c.Find(nil).One(&result); err != nil {
		return nil, err
	} else {
		return inspectResult{
			Name:       name,
			Properties: docToSchema(&result),
			Url:        fmt.Sprintf("/3.0/projects/<PROJECT_ID>/events/%s", name),
		}, nil
	}
}

func inspectAllCollections(p *model.Project) (interface{}, error) {
	r := repo.NewRepository(p.RepoName())
	defer r.Close()

	collectionNames, err := r.CollectionNames()
	if err != nil {
		return nil, err
	}

	results := []interface{}{}
	for _, name := range collectionNames {
		if singleResult, err := inspectCollection(p, name); err != nil {
			return nil, err
		} else {
			results = append(results, singleResult)
		}
	}

	return results, nil
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
