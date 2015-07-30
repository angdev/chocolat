package service

import (
	"github.com/angdev/chocolat/support/repo"
	"gopkg.in/fatih/set.v0"
)

type Arel struct {
	sel   *SelectNode
	where *WhereNode
	group *GroupByNode
}

func NewArel() *Arel {
	var arel Arel
	arel.sel = NewSelectNode()
	arel.where = NewWhereNode()
	arel.group = NewGroupByNode()

	return &arel
}

func (this *Arel) Clone() *Arel {
	clone := Arel{
		sel:   this.sel.Clone(),
		where: this.where.Clone(),
		group: this.group.Clone(),
	}

	return &clone
}

func (this *Arel) Pipeline() *Pipeline {
	var p Pipeline

	clone := this.Clone()
	p.Append(clone.where.Visit(clone))
	p.Append(clone.group.Visit(clone))
	p.Append(clone.sel.Visit(clone))

	return &p
}

func (this *Arel) Select(fields ...string) *Arel {
	return this
}

func (this *Arel) Where(filters ...Filter) *Arel {
	this.where.Append(filters...)
	return this
}

func (this *Arel) GroupBy(groups ...string) *Arel {
	this.group.Append(groups...)
	return this
}

func (this *Arel) TimeFrame(t TimeFrame) *Arel {
	return this.Where(Filter{
		PropertyName:  "chocolat.created_at",
		Operator:      "gt",
		PropertyValue: t.Start,
	}, Filter{
		PropertyName:  "chocolat.created_at",
		Operator:      "lt",
		PropertyValue: t.End,
	})
}

func (this *Arel) Count() *Arel {
	this.group.op = NewCountNode()
	return this
}

func (this *Arel) CountUnique(target string) *Arel {
	this.group.op = NewCountUniqueNode(target)
	return this
}

type ArelNode interface {
	Clone() ArelNode
	Visit(*Arel) Pipe
}

type SelectNode struct {
	fields *set.Set
}

func NewSelectNode() *SelectNode {
	var node SelectNode
	node.fields = set.New()
	return &node
}

func (this *SelectNode) Clone() *SelectNode {
	return &SelectNode{fields: this.fields.Copy().(*set.Set)}
}

func (this *SelectNode) AddFields(fields ...string) {
	args := make([]interface{}, len(fields))
	this.fields.Add(args...)
}

func (this *SelectNode) Visit(area *Arel) Pipe {
	return nil
}

type WhereNode struct {
	filters Filters
}

func NewWhereNode() *WhereNode {
	return &WhereNode{}
}

func (this *WhereNode) Clone() *WhereNode {
	clonedFilter := make([]Filter, len(this.filters))
	copy(clonedFilter, this.filters)
	return &WhereNode{filters: clonedFilter}
}

func (this *WhereNode) Append(filters ...Filter) {
	this.filters = append(this.filters, filters...)
}

func (this *WhereNode) Visit(arel *Arel) Pipe {
	return this.filters.Pipe()
}

type GroupByNode struct {
	groups GroupBy
	op     ArelNode
}

func NewGroupByNode() *GroupByNode {
	node := &GroupByNode{}
	node.op = nil
	return node
}

func (this *GroupByNode) Clone() *GroupByNode {
	clonedGroups := make([]string, len(this.groups))
	copy(clonedGroups, this.groups)
	return &GroupByNode{groups: clonedGroups, op: this.op.Clone()}
}

func (this *GroupByNode) Append(fields ...string) {
	this.groups = append(this.groups, fields...)
}

func (this *GroupByNode) Visit(arel *Arel) Pipe {
	if this.op != nil {
		return this.op.Visit(arel)
	}
	return nil
}

func (this *GroupByNode) IsGiven() bool {
	return (this.groups != nil)
}

type CountNode struct{}

func NewCountNode() *CountNode {
	return &CountNode{}
}

func (this *CountNode) Clone() ArelNode {
	return &(*this)
}

func (this *CountNode) Visit(arel *Arel) Pipe {
	arel.sel.AddFields("result")
	return arel.group.groups.Pipe(repo.Doc{
		"result": repo.Doc{"$sum": 1},
	})
}

type CountUniqueNode struct {
	target string
}

func NewCountUniqueNode(target string) *CountUniqueNode {
	return &CountUniqueNode{target: target}
}

func (this *CountUniqueNode) Clone() ArelNode {
	return &(*this)
}

func (this *CountUniqueNode) Visit(arel *Arel) Pipe {
	arel.sel.AddFields("result")

	groups := arel.group.groups
	targetGroupPipe := append(groups, this.target).Pipe()
	countPipe := groups.Pipe(repo.Doc{
		"result": repo.Doc{"$sum": 1},
	})

	return targetGroupPipe.Join(countPipe)
}
