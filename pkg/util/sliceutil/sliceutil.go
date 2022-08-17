package sliceutil

// Repeat returns a slice with value repeated n times.
func Repeat[T any](value T, n int) []T {
	arr := make([]T, n)
	for i := 0; i < n; i++ {
		arr[i] = value
	}
	return arr
}

func Generate[T any](n int, f func(i int) T) []T {
	arr := make([]T, n)
	for i := 0; i < n; i++ {
		arr[i] = f(i)
	}
	return arr
}

func Transform[T, R any](ts []T, f func(i int, t T) R) []R {
	rs := make([]R, len(ts))
	for i, t := range ts {
		rs[i] = f(i, t)
	}
	return rs
}
