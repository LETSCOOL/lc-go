package dij

import (
	"reflect"
	"unsafe"
)

// GetMethodForReceiverType Get method for a receiver whose type is typ.
func GetMethodForReceiverType(typ reflect.Type, fnName string) (reflect.Method, bool) {
	method, ok := typ.MethodByName(fnName)
	if ok {
		return method, ok
	}
	ptrTyp := reflect.PointerTo(typ)
	return ptrTyp.MethodByName(fnName)
}

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

func FullnameOfType(typ reflect.Type) string {
	name := ""
	for typ.Kind() == reflect.Pointer {
		name += "*"
		typ = typ.Elem()
	}
	return typ.PkgPath() + "/" + typ.Name() + name
}
