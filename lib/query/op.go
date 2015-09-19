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
	distinctGroup := make(RawExpr)
	distinctGroup["_id"] = Group(append(g.Group, this.target)).RawExpr()

	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$sum": 1,
	}

	v.Collect(Stage{
		"$group": group,
	})
}

type Min struct {
	target string
}

func (this *Min) Visit(v *Visitor, g *GroupBy) {
	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$min": Variablize(this.target),
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
		"$max": Variablize(this.target),
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
		"$sum": Variablize(this.target),
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
		"$avg": Variablize(this.target),
	}

	v.Collect(Stage{
		"$group": group,
	})
}

type Collect struct {
	target string
}

func (this *Collect) Visit(v *Visitor, g *GroupBy) {
	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$push": Variablize(this.target),
	}

	v.Collect(Stage{
		"$group": group,
	})
}

type SelectUnique struct {
	target string
}

func (this *SelectUnique) Visit(v *Visitor, g *GroupBy) {
	group := make(RawExpr)
	group["_id"] = g.Group.RawExpr()
	group["result"] = RawExpr{
		"$addToSet": Variablize(this.target),
	}

	v.Collect(Stage{
		"$group": group,
	})
}
