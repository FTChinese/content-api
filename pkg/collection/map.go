package collection

// Map transforms the element of a slice.
func Map[T, U any](in []T, f func(T) U) []U {
	out := make([]U, len(in))

	for i, v := range in {
		out[i] = f(v)
	}

	return out
}
