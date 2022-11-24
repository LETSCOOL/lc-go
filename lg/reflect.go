// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lg

import (
	"reflect"
	"strings"
)

var TypeOfType reflect.Type

func init() {
	TypeOfType = reflect.TypeOf(reflect.TypeOf(struct{}{}))
}

func IsTypeInstance(inst any) bool {
	return reflect.TypeOf(inst) == TypeOfType
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
