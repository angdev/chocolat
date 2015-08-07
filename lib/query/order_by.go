package query

const (
	ASC  = 1
	DESC = -1
)

func NewOrderBy() *OrderBy {
	return &OrderBy{}
}

type OrderBy struct {
	orders []*Order
}

type Order struct {
	Field string
	Order int
}

func (this *OrderBy) AddOrder(orders ...*Order) *OrderBy {
	for _, order := range orders {
		this.orders = append(this.orders, order)
	}

	return this
}

func (this *OrderBy) Visit(v *Visitor) {
	if len(this.orders) == 0 {
		return
	}

	orderMap := make(RawExpr)
	for _, order := range this.orders {
		orderMap[order.Field] = order.Order
	}

	v.Collect(Stage{
		"$sort": orderMap,
	})
}
