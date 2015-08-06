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

type Min struct {
	target string
}

func (this *Min) Visit(v *Visitor, g *GroupBy) {
	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$min": variablize(this.target),
	}

	v.Collect(Stage{
		"$group": group,
	})
}

type Max struct {
	target string
}

func (this *Max) Visit(v *Visitor, g *GroupBy) {
	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$max": variablize(this.target),
	}

	v.Collect(Stage{
		"$group": group,
	})
}

type Sum struct {
	target string
}

func (this *Sum) Visit(v *Visitor, g *GroupBy) {
	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$sum": variablize(this.target),
	}

	v.Collect(Stage{
		"$group": group,
	})
}

type Average struct {
	target string
}

func (this *Average) Visit(v *Visitor, g *GroupBy) {
	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$avg": variablize(this.target),
	}

	v.Collect(Stage{
		"$group": group,
	})
}
