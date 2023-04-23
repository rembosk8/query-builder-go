package stringer

import "bytes"

func SnakeCase(camel string) string {
	var buf bytes.Buffer
	l := len(camel)
	for i, c := range camel {
		if 'A' <= c && c <= 'Z' {
			// just convert [A-Z] to _[a-z]. Convert ABCd into ab_cd
			if buf.Len() > 0 && (i < l-2) && !('A' <= camel[i+1] && camel[i+1] <= 'Z') {
				buf.WriteRune('_')
			}
			buf.WriteRune(c - 'A' + 'a')
		} else {
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
