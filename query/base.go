package query

import "github.com/rembosk8/query-builder-go/query/indent"

const defaultTag = "db"
const all = "*"

type BaseBuilder struct {
	bq baseQuery
}

type Option func(b *baseQuery)

func WithStructAnnotationTag(tag string) Option {
	return func(b *baseQuery) {
		b.tag = tag
	}
}

func WithIndentBuilder(ib *indent.Builder) Option {
	return func(b *baseQuery) {
		b.indentBuilder = ib
	}
}

func New(opts ...Option) BaseBuilder {
	b := baseQuery{indentBuilder: indent.NewBuilder(), tag: defaultTag}

	for _, o := range opts {
		o(&b)
	}

	return BaseBuilder{
		bq: b,
	}
}

func (b BaseBuilder) Select(fields ...string) Select {
	s := Select{
		baseQuery: b.bq,
	}

	for _, f := range fields {
		s.fields = append(s.fields, s.indentBuilder.Indent(f))
	}

	return s
}

func (b BaseBuilder) SelectV2(model any) Select {
	s := Select{
		baseQuery: b.bq,
	}
	s.addFieldsFromModel(model)

	return s
}

func (b BaseBuilder) Update(tableName string) Update {
	u := Update{
		baseQuery: b.bq,
	}
	u.setTable(tableName)
	return u
}

func (b BaseBuilder) DeleteFrom(tableName string) Delete {
	u := Delete{
		baseQuery: b.bq,
	}
	u.setTable(tableName)
	return u
}
