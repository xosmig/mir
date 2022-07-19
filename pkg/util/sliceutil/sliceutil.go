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

func ToSet[T comparable](slice []T) map[T]bool {
	m := make(map[T]bool, len(slice))
	for _, x := range slice {
		m[x] = true
	}
	return m
}
