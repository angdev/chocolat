package query

import (
	"strings"
)

func collapseField(doc RawExpr) RawExpr {
	collapsed := RawExpr{}

	var f func([]string, RawExpr)
	f = func(level []string, cursor RawExpr) {
		for k, v := range cursor {
			switch v.(type) {
			case RawExpr:
				f(append(level, k), v.(RawExpr))
			default:
				collapsed[strings.Join(append(level, k), ".")] = v
			}
		}
	}

	f([]string{}, doc)

	return collapsed
}

func expandField(doc RawExpr) RawExpr {
	expanded := RawExpr{}
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

func deepAssign(d RawExpr, value interface{}, keys ...string) {
	cursor := d
	midKeys, lastKey := keys[:len(keys)-1], keys[len(keys)-1]
	for _, key := range midKeys {
		if _, ok := cursor[key]; !ok {
			cursor[key] = RawExpr{}
		}
		cursor = cursor[key].(RawExpr)
	}
	cursor[lastKey] = value
}
