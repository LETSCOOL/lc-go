package lg

import "strings"

type StructTagAttrs struct {
	segments []string
	attrs    []StructTagAttr
}

func (a StructTagAttrs) Attrs() []StructTagAttr {
	return a.attrs
}

func (a StructTagAttrs) AttrsWithValOnly() []StructTagAttr {
	return Filter(a.attrs, func(attr StructTagAttr) bool {
		return attr.ValOnly
	})
}

func (a StructTagAttrs) FirstAttrWithValOnly() (attr StructTagAttr, exists bool) {
	return FilterFirst(a.attrs, func(attr StructTagAttr) bool {
		return attr.ValOnly
	})
}

func (a StructTagAttrs) AttrsWithKey() []StructTagAttr {
	return Filter(a.attrs, func(attr StructTagAttr) bool {
		return !attr.ValOnly
	})
}

func (a StructTagAttrs) FirstAttrsWithKey() (attr StructTagAttr, exists bool) {
	return FilterFirst(a.attrs, func(attr StructTagAttr) bool {
		return !attr.ValOnly
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
