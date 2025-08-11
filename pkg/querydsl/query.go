package querydsl

// Imported for ctx in repo, but not used here

type Operator string

const (
	OpEqual            Operator = "eq"
	OpContains         Operator = "contains"
	OpGreaterThan      Operator = "gt"
	OpGreaterThanEqual Operator = "gte"
	OpLessThan         Operator = "lt"
	OpLessThanEqual    Operator = "lte"
	// Add more ops as needed (e.g., "in" for arrays)
)

type Filter struct {
	Field string // e.g., "name" or "posts.title"
	Op    Operator
	Value any
}

type Join struct {
	Relation string
	Filters  []Filter
	Joins    []Join // For nested joins
}

type Order struct {
	Field string
	Desc  bool
}

type Query struct {
	Select  []string
	Filters []Filter // Top-level filters
	Joins   []Join
	Orders  []Order
	Limit   int
	Offset  int
}

// Builder for fluent Query construction
type Builder struct {
	q Query
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Select(fields ...string) *Builder {
	b.q.Select = fields
	return b
}

func (b *Builder) Where(field string, op Operator, value any) *Builder {
	b.q.Filters = append(b.q.Filters, Filter{Field: field, Op: op, Value: value})
	return b
}

func (b *Builder) Join(relation string, buildFn func(*JoinBuilder)) *Builder {
	jb := newJoinBuilder(relation)
	buildFn(jb)
	b.q.Joins = append(b.q.Joins, jb.build())
	return b
}

func (b *Builder) OrderBy(field string, desc bool) *Builder {
	b.q.Orders = append(b.q.Orders, Order{Field: field, Desc: desc})
	return b
}

func (b *Builder) Limit(l int) *Builder {
	b.q.Limit = l
	return b
}

func (b *Builder) Offset(o int) *Builder {
	b.q.Offset = o
	return b
}

func (b *Builder) Build() Query {
	return b.q
}

// JoinBuilder for nested joins and filters
type JoinBuilder struct {
	relation string
	filters  []Filter
	joins    []*JoinBuilder // Pointers for nesting
}

func newJoinBuilder(relation string) *JoinBuilder {
	return &JoinBuilder{relation: relation}
}

func (jb *JoinBuilder) Where(field string, op Operator, value any) *JoinBuilder {
	jb.filters = append(jb.filters, Filter{Field: field, Op: op, Value: value})
	return jb
}

func (jb *JoinBuilder) Join(relation string, buildFn func(*JoinBuilder)) *JoinBuilder {
	subjb := newJoinBuilder(relation)
	buildFn(subjb)
	jb.joins = append(jb.joins, subjb)
	return jb
}

func (jb *JoinBuilder) build() Join {
	j := Join{
		Relation: jb.relation,
		Filters:  jb.filters,
	}
	for _, sub := range jb.joins {
		j.Joins = append(j.Joins, sub.build())
	}
	return j
}
