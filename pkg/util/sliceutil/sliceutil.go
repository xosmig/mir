package sliceutil

import "golang.org/x/exp/constraints"

// Repeat returns a slice with value repeated n times.
func Repeat[T any, N constraints.Integer](value T, n N) []T {
	slice := make([]T, n)
	for i := N(0); i < n; i++ {
		slice[i] = value
	}
	return slice
}

func Generate[T any, N constraints.Integer](n N, g func(i N) T) []T {
	slice := make([]T, n)
	for i := N(0); i < n; i++ {
		slice[i] = g(i)
	}
	return slice
}

func Transform[T, R any](ts []T, f func(t T) R) []R {
	rs := make([]R, len(ts))
	for i := range ts {
		rs[i] = f(ts[i])
	}
	return rs
}
