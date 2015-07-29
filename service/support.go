package service

import (
	"github.com/angdev/chocolat/support/repo"
	"strings"
)

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

func expandField(doc repo.Doc) repo.Doc {
	expanded := repo.Doc{}
	collapsed := collapseField(doc)

	for k, v := range collapsed {
		keys := strings.Split(k, ".")
		deepAssign(expanded, v, keys...)
	}

	return expanded
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
