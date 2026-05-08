package utils

// Ptr returns a pointer to any value passed to it.
func Ptr[T any](v T) *T {
	return &v
}
