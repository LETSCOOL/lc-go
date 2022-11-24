// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package lg

import (
	"encoding/json"
	"reflect"
	"testing"
)

// go test ./lg -v -run TestParseStructTag
func TestParseStructTag(t *testing.T) {
	t.Run("ValOnly", func(t *testing.T) {
		attrsObj := ParseStructTag("a, b, c")
		attrs := attrsObj.Attrs()
		//log.Println(attrs)
		if len(attrs) != 3 {
			t.Errorf("incorrect number attributes: %d", len(attrs))
		}
		attrsValOnly := attrsObj.AttrsWithValOnly()
		if len(attrsValOnly) != 3 {
			t.Errorf("incorrect number attributes (val only): %d", len(attrsValOnly))
		}
		if !reflect.DeepEqual(attrs, attrsValOnly) {
			t.Errorf("not equal: %v != %v", attrs, attrsValOnly)
		}
	})
	t.Run("WithKey", func(t *testing.T) {
		attrsObj := ParseStructTag("a, b=bb , c=ccc")
		attrs := attrsObj.Attrs()
		attrsValOnly := attrsObj.AttrsWithValOnly()
		attrsWithKey := attrsObj.AttrsWithKey()
		//log.Println(attrs)
		if len(attrs) != 3 {
			t.Errorf("incorrect number attributes: %d", len(attrs))
		}
		if len(attrsValOnly) != 1 {
			t.Errorf("incorrect number attributes (val only): %d", len(attrsValOnly))
		}
		//t.Log(attrsValOnly)
		if len(attrsWithKey) != 2 {
			t.Errorf("incorrect number attributes (with key): %d", len(attrsWithKey))
		}
		//t.Log(attrsWithKey)
		if reflect.DeepEqual(attrs, attrsValOnly) {
			t.Errorf("should not equal: %v != %v", attrs, attrsValOnly)
		}
		if reflect.DeepEqual(attrs, attrsWithKey) {
			t.Errorf("should not equal: %v != %v", attrs, attrsWithKey)
		}
		if reflect.DeepEqual(attrsValOnly, attrsWithKey) {
			t.Errorf("should not equal: %v != %v", attrsValOnly, attrsWithKey)
		}
	})
	t.Run("FirstOnly", func(t *testing.T) {
		attrsObj := ParseStructTag("a, b=bb , c=ccc, d")
		if first, exists := attrsObj.FirstAttrWithValOnly(); !exists || first.Orig != "a" {
			t.Errorf("error: %v, %v", first, exists)
		}

		if first, exists := attrsObj.FirstAttrsWithKey("c"); !exists || first.Key != "c" || first.Val != "ccc" {
			t.Errorf("error: %v, %v", first, exists)
		}
	})
	t.Run("NoFirst", func(t *testing.T) {
		attrsObj := ParseStructTag("b=bb,c=ccc, d")
		if _, exists := attrsObj.FirstAttrWithValOnly(); exists {
			t.Errorf("First attr should not exists")
		}
	})
	t.Run("EmptyFirst", func(t *testing.T) {
		attrsObj := ParseStructTag(",b=bb,c=ccc, d")
		if attr, exists := attrsObj.FirstAttrWithValOnly(); !exists || attr.Val != "" {
			t.Errorf("First attr should be emtpy string(%v, %v)", exists, attr)
		}
	})

}

// go test ./lg -v -run TestUnmarshal
func TestUnmarshal(t *testing.T) {
	t.Run("all types", func(t *testing.T) {
		var i int16
		if err := json.Unmarshal([]byte("123"), &i); err != nil {
			t.Error(err)
		} else {
			if i != 123 {
				t.Errorf("incorrect value: %d", i)
			}
			//t.Log(i)
		}
	})
}
