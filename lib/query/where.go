package query

import (
	"github.com/imdario/mergo"
)

func NewWhere() *Where {
	return &Where{Conditions: make(map[string]*Condition)}
}

type Where struct {
	Conditions map[string]*Condition
}

func (this *Where) Condition(conds ...*Condition) *Where {
	for _, cond := range conds {
		hash := cond.Field + "," + cond.Op
		this.Conditions[hash] = cond
	}
	return this
}

func (this *Where) Visit(v *Visitor) {
	match := make(RawExpr)

	for _, cond := range this.Conditions {
		mergo.MapWithOverwrite(&match, cond.RawExpr())
	}

	v.Collect(Stage{
		"$match": match,
	})
}
