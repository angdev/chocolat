package query

type Arel struct {
	*ArelNodes
}

type ArelNodes struct {
	Select  *Select
	Where   *Where
	GroupBy *GroupBy
}

func NewArel() *Arel {
	return &Arel{
		&ArelNodes{
			Select:  NewSelect(),
			Where:   NewWhere(),
			GroupBy: NewGroupBy(),
		},
	}
}

func (this *Arel) Pipeline() *Pipeline {
	v := NewVisitor(this)

	this.ArelNodes.Where.Visit(v)
	this.ArelNodes.GroupBy.Visit(v)

	return &v.Pipeline
}

func (this *Arel) GroupByGiven() bool {
	return len(this.ArelNodes.GroupBy.Group) != 0
}

func (this *Arel) Select(fields ...string) *Arel {
	return this
}

func (this *Arel) Where(conds ...*Condition) *Arel {
	this.ArelNodes.Where.Condition(conds...)
	return this
}

func (this *Arel) GroupBy(groups ...string) *Arel {
	this.ArelNodes.GroupBy.AddGroup(groups...)
	return this
}

func (this *Arel) Count() *Arel {
	op := &Count{}
	this.ArelNodes.GroupBy.SetOp(op)
	return this
}

func (this *Arel) CountUnique(target string) *Arel {
	op := &CountUnique{target: target}
	this.ArelNodes.GroupBy.SetOp(op)
	return this
}

func (this *Arel) Min(target string) *Arel {
	op := &Min{target: target}
	this.ArelNodes.GroupBy.SetOp(op)
	return this
}

func (this *Arel) Max(target string) *Arel {
	op := &Max{target: target}
	this.ArelNodes.GroupBy.SetOp(op)
	return this
}

func (this *Arel) Sum(target string) *Arel {
	op := &Sum{target: target}
	this.ArelNodes.GroupBy.SetOp(op)
	return this
}

func (this *Arel) Average(target string) *Arel {
	op := &Average{target: target}
	this.ArelNodes.GroupBy.SetOp(op)
	return this
}
