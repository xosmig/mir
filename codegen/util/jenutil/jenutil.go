package jenutil

import (
	"fmt"
	"reflect"

	"github.com/dave/jennifer/jen"
)

func QualFromType(tp reflect.Type) *jen.Statement {
	// Defined type with a package path
	if tp.PkgPath() != "" && tp.Name() != "" {
		return jen.Qual(tp.PkgPath(), tp.Name())
	}

	// Built-in defined type (e.g., int, string, error, etc)
	if tp.Name() != "" {
		return jen.Id(tp.Name())
	}

	// Non-defined types (e.g., arrays, pointers, etc)
	switch tp.Kind() {

	case reflect.Array:
		return jen.Index(jen.Lit(tp.Len())).Add(QualFromType(tp.Elem()))

	case reflect.Pointer:
		return jen.Op("*").Add(QualFromType(tp.Elem()))

	case reflect.Slice:
		return jen.Index().Add(QualFromType(tp.Elem()))

	case reflect.Map:
		return jen.Map(QualFromType(tp.Key())).Add(QualFromType(tp.Elem()))

	case reflect.Interface:
		return jen.InterfaceFunc(func(g *jen.Group) {
			for methodIdx := 0; methodIdx < tp.NumMethod(); methodIdx++ {
				m := tp.Method(methodIdx)
				g.Id(m.Name).ParamsFunc(func(g *jen.Group) {
					for inIdx := 0; inIdx < m.Type.NumIn(); inIdx++ {
						g.Add(QualFromType(m.Type.In(inIdx)))
					}
				}).ParamsFunc(func(group *jen.Group) {
					for outIdx := 0; outIdx < m.Type.NumOut(); outIdx++ {
						group.Add(QualFromType(m.Type.Out(outIdx)))
					}
				})
			}
		})

	case reflect.Struct:
		return jen.StructFunc(func(g *jen.Group) {
			for fieldIdx := 0; fieldIdx < tp.NumField(); fieldIdx++ {
				f := tp.Field(fieldIdx)
				if f.Anonymous {
					g.Add(QualFromType(f.Type))
				} else {
					g.Id(f.Name).Add(QualFromType(f.Type))
				}
			}
		})

	case reflect.Func:
		return jen.Func().ParamsFunc(func(g *jen.Group) {
			for inIdx := 0; inIdx < tp.NumIn(); inIdx++ {
				g.Add(QualFromType(tp.In(inIdx)))
			}
		}).ParamsFunc(func(g *jen.Group) {
			for outIdx := 0; outIdx < tp.NumOut(); outIdx++ {
				g.Add(QualFromType(tp.Out(outIdx)))
			}
		})

	case reflect.Chan:
		switch tp.ChanDir() {
		case reflect.RecvDir:
			return jen.Op("<-").Chan().Add(QualFromType(tp.Elem()))
		case reflect.SendDir:
			return jen.Chan().Op("<-").Add(QualFromType(tp.Elem()))
		case reflect.BothDir:
			return jen.Chan().Add(QualFromType(tp.Elem()))
		default:
			panic(fmt.Errorf("unexpected ChanDir: %v", tp.ChanDir()))
		}

	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32,
		reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String, reflect.UnsafePointer:
		panic(fmt.Errorf("type of kind %v cannot be non-defined", tp.Kind()))

	default:
		panic(fmt.Errorf("unknown go type kind: %v", tp.Kind()))
	}
}
