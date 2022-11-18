package dij

// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import (
	"fmt"
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
	switch typ.Kind() {
	case reflect.Struct, reflect.Interface:
		return fmt.Sprintf("%s/%s%s", typ.PkgPath(), typ.Name(), name)
	default:
		return fmt.Sprintf("%s/%s", typ.PkgPath(), typ)
	}
}
