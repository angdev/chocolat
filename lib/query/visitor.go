package query

func NewVisitor(arel *Arel) *Visitor {
	return &Visitor{
		Arel: arel,
	}
}

type Visitor struct {
	Arel     *Arel
	Pipeline Pipeline
}

func (this *Visitor) Collect(stages ...Stage) *Visitor {
	this.Pipeline.Append(stages...)
	return this
}
