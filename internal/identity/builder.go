package identity

import "fmt"

type Builder struct {
	indentSanitizer Sanitizer
	valSanitizer    ValueSanitizer
}

func NewBuilder(opts ...BuilderOption) *Builder {
	b := &Builder{}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *Builder) Ident(name string) Identity {
	if b.indentSanitizer != nil {
		name = b.indentSanitizer.Sanitize(name)
	}
	return Identity(name)
}

func (b *Builder) Idents(names ...string) []Identity {
	// todo: think about sanitizing on place, not post hoc
	if b.valSanitizer != nil {
		res := make([]Identity, len(names))
		for i := range names {
			res[i] = b.indentSanitizer.Sanitize(names[i])
		}

		return res
	}

	return names
}

func (b *Builder) Value(val any) string {
	// todo: think about sanitizing on place, not post hoc
	if b.valSanitizer != nil {
		return b.valSanitizer.Sanitize(val)
	}

	return fmt.Sprintf("%v", val)
}

func (b *Builder) Values(vals []any) []string {
	// todo: think about sanitizing on place, not post hoc
	res := make([]string, len(vals))
	if b.valSanitizer != nil {
		for i := range vals {
			res[i] = b.valSanitizer.Sanitize(vals[i])
		}

		return res
	}

	for i := range vals {
		res[i] = fmt.Sprintf("%v", vals[i])
	}

	return res
}

func (b *Builder) IsStandard(v Value) bool {
	if b.valSanitizer != nil {
		return b.valSanitizer.IsStandard(v)
	}

	return false
}

type BuilderOption func(builder *Builder)

func WithIndentSerializer(s Sanitizer) BuilderOption {
	return func(b *Builder) {
		b.indentSanitizer = s
	}
}

func WithValueSerializer(s ValueSanitizer) BuilderOption {
	return func(b *Builder) {
		b.valSanitizer = s
	}
}
