// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lg

import (
	"reflect"
	"strings"
	"time"
)

var TypeOfType reflect.Type
var TypeOfErrorIf reflect.Type
var TypeOfTime reflect.Type

func init() {
	TypeOfType = reflect.TypeOf(reflect.TypeOf(struct{}{}))
	type Er struct {
		First error
	}
	TypeOfErrorIf = reflect.TypeOf(Er{}).Field(0).Type
	TypeOfTime = reflect.TypeOf(time.Now())
}

func IsTypeInstance(inst any) bool {
	return reflect.TypeOf(inst) == TypeOfType
}

func IsError(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Pointer:
		return IsError(t.Elem())
	case reflect.Struct:
		if method, ok := t.MethodByName("Error"); ok {
			methodTyp := method.Type
			if methodTyp.NumOut() == 1 && methodTyp.Out(0).Kind() == reflect.String {
				return true
			}
		}
	case reflect.Interface:
		return t == TypeOfErrorIf
	default:
	}
	return false
}

type StructTagAttrs struct {
	segments []string
	attrs    []StructTagAttr // in order, same order as segments
}

func (a *StructTagAttrs) Attrs() []StructTagAttr {
	return a.attrs
}

func (a *StructTagAttrs) AttrsWithValOnly() []StructTagAttr {
	return Filter(a.attrs, func(attr StructTagAttr) bool {
		return attr.ValOnly
	})
}

func (a *StructTagAttrs) ContainsAttrWithValOnly(val string) bool {
	_, exists := FilterFirst(a.attrs, func(attr StructTagAttr) bool {
		return attr.ValOnly && attr.Val == val
	})
	return exists
}

func (a *StructTagAttrs) FilterAttrWithValOnly(f func(attr StructTagAttr) bool) []StructTagAttr {
	return Filter(a.attrs, func(attr StructTagAttr) bool {
		return attr.ValOnly && f(attr)
	})
}

func (a *StructTagAttrs) FilterValOnly(f func(val string) bool) []string {
	return FilterAndMap(a.attrs, func(attr StructTagAttr) (d string, is bool) {
		if is = attr.ValOnly && f(attr.Val); is {
			d = attr.Val
		}
		return
	})
}

func (a *StructTagAttrs) FirstAttrWithValOnly() (attr StructTagAttr, exists bool) {
	if len(a.attrs) > 0 {
		if attr := a.attrs[0]; attr.ValOnly {
			return attr, true
		}
	}
	return
}

func (a *StructTagAttrs) AttrsWithKey() []StructTagAttr {
	return Filter(a.attrs, func(attr StructTagAttr) bool {
		return !attr.ValOnly
	})
}

func (a *StructTagAttrs) FirstAttrsWithKey(key string) (attr StructTagAttr, exists bool) {
	return FilterFirst(a.attrs, func(attr StructTagAttr) bool {
		return !attr.ValOnly && attr.Key == key
	})
}

type StructTagAttr struct {
	Orig    string
	Key     string
	Val     string
	ValOnly bool // Orig is not "xxx=yyy" pattern
}

func ParseStructTag(tagValue string) StructTagAttrs {
	segments := Map(strings.Split(tagValue, ","), func(v string) string {
		return strings.TrimSpace(v)
	})

	attrs := make([]StructTagAttr, 0, len(segments))

	for _, text := range segments {
		i := strings.Index(text, "=")
		var attr StructTagAttr
		attr.Orig = text
		if i <= 0 {
			attr.Val = text
			attr.ValOnly = true
		} else {
			attr.Key = text[0:i]
			attr.Val = text[i+1:]
			attr.ValOnly = false
		}
		attrs = append(attrs, attr)
	}

	return StructTagAttrs{
		segments: segments,
		attrs:    attrs,
	}
}
