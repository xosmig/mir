package cmputil

import (
	"github.com/google/go-cmp/cmp"
	"unicode"
	"unicode/utf8"
)

// IgnoreAllUnexported returns a cmp.Ignore option that ignores all unexported
// struct fields during the comparison.
func IgnoreAllUnexported() cmp.Option {
	filter := func(path cmp.Path) bool {
		sf, ok := path.Last().(cmp.StructField)
		return ok && !isExported(sf.Name())
	}

	return cmp.FilterPath(filter, cmp.Ignore())
}

// isExported reports whether the identifier is exported.
func isExported(id string) bool {
	r, _ := utf8.DecodeRuneInString(id)
	return unicode.IsUpper(r)
}
