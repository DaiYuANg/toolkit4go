package dbx

import (
	"reflect"
	"slices"
)

var (
	signedIntKinds   = []reflect.Kind{reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32}
	unsignedIntKinds = []reflect.Kind{reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32}
)

func isSignedIntKind(kind reflect.Kind) bool {
	return slices.Contains(signedIntKinds, kind)
}

func isUnsignedIntKind(kind reflect.Kind) bool {
	return slices.Contains(unsignedIntKinds, kind)
}

func isByteSliceType(typ reflect.Type) bool {
	return typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Uint8
}
