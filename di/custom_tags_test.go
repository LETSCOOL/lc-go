package di

import (
	"reflect"
	"testing"
)

type IsUser interface {
	GetName() string
	SetName(n string)
}

type User struct {
	Id    int    `validate:"-"`
	Name  string `validate:"presence,min=2,max=32"`
	Email string `validate:"email,required"`
}

func (u User) GetName() string {
	return u.Name
}

func (u User) SetName(n string) {
	u.Name = n
}

var tagName = "validate"

// go test lc-go/di -v -run TestTags
func TestTags(t *testing.T) {
	user := User{
		Id:    1,
		Name:  "John Doe",
		Email: "john@example",
	}
	var inf interface{} = user

	isUser := inf.(IsUser)

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	typ := reflect.TypeOf(isUser)

	// Get the type and kind of our user variable
	t.Log("Type:", typ.Name())
	t.Log("Kind:", typ.Kind(), typ.Kind() == reflect.Struct)
	t.Logf("Num of Method: %d, of Field: %d\n", typ.NumMethod(), typ.NumField())

	t.Logf("")
	t.Logf("Methods:")
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		methodType := method.Type

		t.Logf("%d. %v (%v is kind of %v)\n", i+1, method.Name, methodType, methodType.Kind())

		for j := 0; j < methodType.NumIn(); j++ {
			paramTyp := methodType.In(j)
			t.Log("\t", j, ". ", paramTyp)
		}
	}

	t.Logf("")
	t.Logf("Fields:")
	// Iterate over all available fields and read the tag value
	for i := 0; i < typ.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := typ.Field(i)

		// Get the field tag value
		// field.Tag.Lookup()
		tag := field.Tag.Get(tagName)

		t.Logf("%d. %v (%v), tag: '%v'\n", i+1, field.Name, field.Type.Name(), tag)
	}
}
