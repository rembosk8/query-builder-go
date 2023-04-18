package indent

import (
	"fmt"
)

type Indent struct {
	Name      string
	sanitizer Sanitizer
}

var _ fmt.Stringer = &Indent{}

func (s Indent) String() string {
	if s.sanitizer != nil {
		return s.sanitizer.Sanitize(s.Name)
	}

	return s.Name
}

type Value struct {
	Value     any
	sanitizer ValueSanitizer
}

func (v Value) String() string {
	if v.sanitizer != nil {
		return v.sanitizer.Sanitize(v.Value)
	}

	return fmt.Sprintf("%v", v.Value)
}

var _ fmt.Stringer = &Value{}

type Builder struct {
	indentSanitizer Sanitizer
	valSanitizer    ValueSanitizer
}

func NewBuilder(sanitizer Sanitizer) Builder {
	return Builder{indentSanitizer: sanitizer}
}

func (b Builder) Indent(name string) Indent {
	return Indent{
		Name:      name,
		sanitizer: b.indentSanitizer,
	}
}

func (b Builder) Value(val any) Value {
	return Value{
		Value:     val,
		sanitizer: b.valSanitizer,
	}
}
