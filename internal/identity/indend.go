package identity

import (
	"fmt"
)

type Identity string

func (s Identity) String() string {
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

func (v Value) IsStandard() bool {
	if v.sanitizer != nil {
		return v.sanitizer.IsStandard(v.Value)
	}

	return false
}
