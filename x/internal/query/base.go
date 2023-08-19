package query

import (
	"github.com/rembosk8/query-builder-go/x/internal/identity"
)

const (
	defaultTag = "db"
	all        = "*"
)

type BaseBuilder struct {
	bq baseQuery
}

type Option func(b *baseQuery)

func WithStructAnnotationTag(tag string) Option {
	return func(b *baseQuery) {
		b.tag = tag
	}
}

func WithIdentityBuilder(ib *identity.Builder) Option {
	return func(b *baseQuery) {
		b.indentBuilder = ib
	}
}

func New(opts ...Option) BaseBuilder {
	b := baseQuery{indentBuilder: identity.NewBuilder(), tag: defaultTag}

	for _, o := range opts {
		o(&b)
	}

	return BaseBuilder{
		bq: b,
	}
}

func (b BaseBuilder) Select(fields ...string) *SelectCore {
	s := SelectCore{
		core: core{indentBuilder: b.bq.indentBuilder},
	}

	s.fields = append(s.fields, fields...)

	return &s
}

//func (b BaseBuilder) SelectV2(model any) Select {
//	s := Select{
//		baseQuery: b.bq,
//	}
//	s.addFieldsFromModel(model)
//
//	return s
//}

//func (b BaseBuilder) Update(tableName string) Update {
//	u := Update{
//		baseQuery: b.bq,
//	}
//	u.setTable(tableName)
//	return u
//}
//
//func (b BaseBuilder) DeleteFrom(tableName string) Delete {
//	u := Delete{
//		baseQuery: b.bq,
//	}
//	u.setTable(tableName)
//	return u
//}
//
//func (b BaseBuilder) InsertInto(tableName string) Insert {
//	i := Insert{
//		baseQuery: b.bq,
//	}
//	i.setTable(tableName)
//	return i
//}
