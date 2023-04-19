package pointer

// To returns a pointer to the passed value.
func To[T any](t T) *T {
	return &t
}

// Get returns the value from the passed pointer or the zero value if the pointer is nil.
func Get[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}
