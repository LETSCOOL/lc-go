package dij

import (
	"fmt"
	"reflect"
	"testing"
)

type TestApp struct {
	test string    `di:"-"`
	lib3 *TestLib3 `di:"^"`
}

func (a *TestApp) InjectDependency(p struct {
	cfg  map[string]any `di:"config"`
	lib1 *TestLib1      `di:"Lib1"`
	lib2 *TestLib2      `di:"Lib2"`
}) {
	fmt.Printf("TestApp\n\tcfg: %v\n\tlib1: %p %v\n\tlib2: %p %v\n\tlib3: %p %v\n",
		p.cfg,
		p.lib1, reflect.TypeOf(p.lib1),
		p.lib2, reflect.TypeOf(p.lib2),
		a.lib3, reflect.TypeOf(a.lib3))
}

type TestLib1 struct {
	v int
}

func (l *TestLib1) InjectDependency(p struct {
	cfg map[string]any `di:"config"`
}) {
	fmt.Printf("TestLib1\n\taddress: %p\n\tcfg: %v\n", l, p.cfg)
}

type TestLib2 struct {
	v int
}

func (l *TestLib2) InjectDependency(p struct {
	lib1 *TestLib1 `di:"Lib1"`
}) {
	fmt.Printf("TestLib2\n\taddress: %p\n\tlib1: %p\n", l, p.lib1)
}

type TestLib3 struct {
	v   int
	tmp *TestApp `di:"^"`
}

// go test lc-go/dij -v -run TestDI
func TestDI(t *testing.T) {
	appTyp := reflect.TypeOf(TestApp{})
	ref := map[DependencyKey]any{}
	config := map[string]any{
		"ip":   "192.168.0.1",
		"port": 3345,
	}
	ref["config"] = config

	inst, err := CreateInstance(appTyp, &ref, "^")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Inst type:", reflect.TypeOf(inst))
}
