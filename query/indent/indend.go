package indent

import (
	"fmt"
)

type Indent string

func (s Indent) String() string {
	return string(s)
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
