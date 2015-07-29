package service

import (
	"strings"
)

func collapseField(doc map[string]interface{}) map[string]interface{} {
	collapsed := map[string]interface{}{}

	var f func([]string, map[string]interface{})
	f = func(level []string, cursor map[string]interface{}) {
		for k, v := range cursor {
			switch v.(type) {
			case map[string]interface{}:
				f(append(level, k), v.(map[string]interface{}))
			default:
				collapsed[strings.Join(append(level, k), ".")] = v
			}
		}
	}

	f([]string{}, doc)

	return collapsed
}

func expandField(doc map[string]interface{}) map[string]interface{} {
	expanded := map[string]interface{}{}
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

func deepAssign(d map[string]interface{}, value interface{}, keys ...string) {
	cursor := d
	midKeys, lastKey := keys[:len(keys)-1], keys[len(keys)-1]
	for _, key := range midKeys {
		if _, ok := cursor[key]; !ok {
			cursor[key] = map[string]interface{}{}
		}
		cursor = cursor[key].(map[string]interface{})
	}
	cursor[lastKey] = value
}
