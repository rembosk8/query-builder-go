package indent

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

func (b *Builder) Indent(name string) Indent {
	if b.indentSanitizer != nil {
		name = b.indentSanitizer.Sanitize(name)
	}
	return Indent(name)
}

func (b *Builder) Value(val any) Value {
	return Value{
		Value:     val,
		sanitizer: b.valSanitizer,
	}
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
