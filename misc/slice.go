package misc

func IndexOf[T comparable](slice []T, target T) int {
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

func Find[T any](slice []T, fn func(T) bool) (int, T) {
	for i, v := range slice {
		if fn(v) {
			return i, v
		}
	}
	var zero T
	return -1, zero
}

func MapSlice[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}
