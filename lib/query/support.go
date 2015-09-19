package query

import (
	"strings"
)

func CollapseField(doc RawExpr) RawExpr {
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

func ExpandField(doc RawExpr) RawExpr {
	expanded := RawExpr{}
	collapsed := CollapseField(doc)

	for k, v := range collapsed {
		keys := strings.Split(k, ".")
		DeepAssign(expanded, v, keys...)
	}

	return expanded
}

func Variablize(fields ...string) string {
	if len(fields) == 0 {
		return ""
	} else {
		return "$" + strings.Join(fields, ".")
	}
}

func DeepAssign(d RawExpr, value interface{}, keys ...string) {
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
