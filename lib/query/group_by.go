package query

func NewGroupBy() *GroupBy {
	return &GroupBy{}
}

type GroupBy struct {
	Group Group
	Op    Operator
}

func (this *GroupBy) AddGroup(groups ...string) *GroupBy {
	this.Group = append(this.Group, groups...)
	return this
}

func (this *GroupBy) SetOp(op Operator) *GroupBy {
	this.Op = op
	return this
}

func (this *GroupBy) Visit(v *Visitor) {
	if this.Op != nil {
		this.Op.Visit(v, this)
	} else {
		group := make(RawExpr)
		group["_id"] = this.Group.RawExpr()

		v.Collect(Stage{
			"$group": group,
		})
	}

	// Need to capsule
	project := make(RawExpr)
	groups := make(RawExpr)
	for _, field := range this.Group {
		groups[field] = variablize("_id", field)
	}
	project["_id"] = false
	project["groups"] = groups
	project["result"] = variablize("result")

	v.Collect(Stage{
		"$project": project,
	})
}

type Group []string

func (this Group) RawExpr() RawExpr {
	if len(this) == 0 {
		return nil
	}

	expr := make(RawExpr)
	for _, field := range this {
		expr[field] = variablize(field)
	}

	return expandField(expr)
}
