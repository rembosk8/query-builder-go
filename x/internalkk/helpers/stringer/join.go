package stringer

import (
	"fmt"
	"strings"
)

func Join[T fmt.Stringer](elems []T, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0].String()
	}

	var b strings.Builder

	b.WriteString(elems[0].String())
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s.String())
	}

	return b.String()
}
