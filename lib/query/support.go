package query

import (
	"reflect"
	"strings"
)

func CollapseField(in interface{}, out interface{}) {
	inValue := reflect.ValueOf(in)
	outValue := reflect.ValueOf(out).Elem()

	var f func([]string, reflect.Value)
	f = func(level []string, cursor reflect.Value) {
		for _, k := range cursor.MapKeys() {
			v := cursor.MapIndex(k).Elem()
			switch v.Kind() {
			case reflect.Map:
				f(append(level, k.String()), v)
			default:
				key := strings.Join(append(level, k.String()), ".")
				outValue.SetMapIndex(reflect.ValueOf(key), v)
			}
		}
	}

	f([]string{}, inValue)
}

func ExpandField(doc RawExpr) RawExpr {
	expanded := RawExpr{}
	collapsed := make(RawExpr)
	CollapseField(doc, &collapsed)

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
