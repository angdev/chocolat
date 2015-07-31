package query

type Stage RawExpr

type Pipeline []Stage

func (this *Pipeline) Append(stages ...Stage) *Pipeline {
	*this = append(*this, stages...)
	return this
}
