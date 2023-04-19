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

var _ fmt.Stringer = &Value{}

func (v Value) String() string {
	if v.sanitizer != nil {
		return v.sanitizer.Sanitize(v.Value)
	}

	return fmt.Sprintf("%v", v.Value)
}
