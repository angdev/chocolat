package query

type Arel struct {
	*ArelNodes
}

type ArelNodes struct {
	Select  *Select
	Where   *Where
	GroupBy *GroupBy
	OrderBy *OrderBy
	Limit   *Limit
}

func NewArel() *Arel {
	return &Arel{
		&ArelNodes{
			Select:  NewSelect(),
			Where:   NewWhere(),
			GroupBy: NewGroupBy(),
			OrderBy: NewOrderBy(),
			Limit:   NewLimit(),
		},
	}
}

func (this *Arel) Pipeline() *Pipeline {
	v := NewVisitor(this)

	this.ArelNodes.Select.Visit(v, this)
	this.ArelNodes.Where.Visit(v, this)
	this.ArelNodes.OrderBy.Visit(v, this)
	this.ArelNodes.Limit.Visit(v, this)
	this.ArelNodes.GroupBy.Visit(v, this)

	return &v.Pipeline
}

func (this *Arel) GroupByGiven() bool {
	return len(this.ArelNodes.GroupBy.Group) != 0
}

func (this *Arel) Select(fields ...string) *Arel {
	this.ArelNodes.Select.AddField(fields...)
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

func (this *Arel) OrderBy(orders ...*Order) *Arel {
	this.ArelNodes.OrderBy.AddOrder(orders...)
	return this
}

func (this *Arel) Limit(n int) *Arel {
	this.ArelNodes.Limit.N = n
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

func (this *Arel) Collect(target string) *Arel {
	op := &Collect{target: target}
	this.ArelNodes.GroupBy.SetOp(op)
	return this
}

func (this *Arel) SelectUnique(target string) *Arel {
	op := &SelectUnique{target: target}
	this.ArelNodes.GroupBy.SetOp(op)
	return this
}
