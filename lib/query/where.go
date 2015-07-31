package query

import (
	"github.com/imdario/mergo"
	"github.com/k0kubun/pp"
)

func NewWhere() *Where {
	return &Where{Conditions: make(map[string]*Condition)}
}

type Where struct {
	Conditions map[string]*Condition
}

func (this *Where) Condition(conds ...*Condition) *Where {
	for _, cond := range conds {
		this.Conditions[cond.Field] = cond
	}
	return this
}

func (this *Where) Visit(v *Visitor) {
	match := make(RawExpr)

	pp.Println(match)

	for _, cond := range this.Conditions {
		mergo.Merge(&match, cond.RawExpr())
	}

	v.Collect(Stage{
		"$match": match,
	})
}
