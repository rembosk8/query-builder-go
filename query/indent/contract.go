package indent

type Sanitizer interface {
	Sanitize(v string) string
}

type ValueSanitizer interface {
	Sanitize(v any) string
	IsDefault(v any) bool
}
