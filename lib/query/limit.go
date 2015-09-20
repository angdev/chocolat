package query

func NewLimit() *Limit {
	return &Limit{-1}
}

type Limit struct {
	N int
}

func (this *Limit) Visit(v *Visitor, arel *Arel) {
	if this.N == -1 {
		return
	}

	v.Collect(Stage{
		"$limit": this.N,
	})
}
