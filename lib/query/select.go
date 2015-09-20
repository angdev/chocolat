package query

func NewSelect() *Select {
	return &Select{Fields: []string{}}
}

type Select struct {
	Fields []string
}

func (this *Select) AddField(fields ...string) *Select {
	this.Fields = append(this.Fields, fields...)
	return this
}

func (this *Select) Visit(v *Visitor, arel *Arel) {
	if len(this.Fields) == 0 {
		return
	}

	op := &NoOp{this.Fields}
	arel.ArelNodes.GroupBy.SetOp(op)
}
