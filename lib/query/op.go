package query

type Operator interface {
	Visit(*Visitor, *GroupBy)
}

type Count struct{}

func (this *Count) Visit(v *Visitor, g *GroupBy) {
	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$sum": 1,
	}

	v.Collect(Stage{
		"$group": group,
	})
}

type CountUnique struct {
	target string
}

func (this *CountUnique) Visit(v *Visitor, g *GroupBy) {
	group := g.Group
	distinctGroup := &GroupBy{
		Group: append(group, this.target),
	}

	distinctGroup.Visit(v)

	countGroup := &GroupBy{
		Group: group,
		Op:    &Count{},
	}

	countGroup.Visit(v)
}
