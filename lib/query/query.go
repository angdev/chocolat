package query

type RawExpr map[string]interface{}

func New() *Query {
	return &Query{}
}

type Query struct {
}
