package query

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"regexp"
)

func NewCondition(field string, op string, value interface{}) *Condition {
	return &Condition{
		Field: field,
		Op:    op,
		Value: value,
	}
}

type Condition struct {
	Field string
	Op    string
	Value interface{}
}

func (this *Condition) RawExpr() RawExpr {
	condExpr := make(RawExpr)
	condExpr[this.Field] = this.OpExpr()
	return condExpr
}

func (this *Condition) OpExpr() RawExpr {
	var (
		op    string
		value interface{}
	)
	expr := make(RawExpr)

	switch this.Op {
	case "contains":
		op = "$regex"
		value = bson.RegEx{regexp.QuoteMeta(this.Value.(string)), ""}
	case "not_contains":
		op = "$not"
		value = bson.RegEx{regexp.QuoteMeta(this.Value.(string)), ""}
	default:
		op = fmt.Sprintf("$%s", this.Op)
		value = this.Value
	}

	expr[op] = value
	return expr
}
