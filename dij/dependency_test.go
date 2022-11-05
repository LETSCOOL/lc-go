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

func (a *TestApp) DidDependencyInjection() {
	fmt.Printf("TestApp - DidDependencyInjection\n")
}

func (a *TestApp) DidDependencyInitialization() {
	fmt.Printf("TestApp - DidDependencyInitialization\n")
}

type TestLib1 struct {
	v int
}

func (l *TestLib1) InjectDependency(p struct {
	cfg map[string]any `di:"config"`
}) {
	fmt.Printf("TestLib1\n\taddress: %p\n\tcfg: %v\n", l, p.cfg)
}

func (l *TestLib1) DidDependencyInjection() {
	fmt.Printf("TestLib1 - DidDependencyInjection\n")
}

func (l *TestLib1) DidDependencyInitialization() {
	fmt.Printf("TestLib1 - DidDependencyInitialization\n")
}

type TestLib2 struct {
	v int
}

func (l *TestLib2) InjectDependency(p struct {
	lib1 *TestLib1 `di:"Lib1"`
}) {
	fmt.Printf("TestLib2\n\taddress: %p\n\tlib1: %p\n", l, p.lib1)
}

func (l *TestLib2) DidDependencyInjection() {
	fmt.Printf("TestLib2 - DidDependencyInjection\n")
}

func (l *TestLib2) DidDependencyInitialization() {
	fmt.Printf("TestLib2 - DidDependencyInitialization\n")
}

type TestLib3 struct {
	v   int
	tmp *TestApp `di:"^"`
}

func (l *TestLib3) DidDependencyInjection() {
	fmt.Printf("TestLib3 - DidDependencyInjection\n")
}

func (l *TestLib3) DidDependencyInitialization() {
	fmt.Printf("TestLib3 - DidDependencyInitialization\n")
}

// go test ./dij -v -run TestDI
func TestDI(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		appTyp := reflect.TypeOf(TestApp{})
		ref := map[DependencyKey]any{}
		config := map[string]any{
			"ip":   "192.168.0.1",
			"port": 3345,
		}
		ref["config"] = config
		EnableLog()

		inst, err := CreateInstance(appTyp, &ref, "^")
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := inst.(*TestApp); ok {
			//fmt.Println("Inst type:", reflect.TypeOf(inst))
		} else {
			t.Fatal("didn't create a create instance, ", reflect.TypeOf(inst))
		}

		if count := GetCountOfDependencyStack(&ref); count != 0 {
			t.Errorf("incorrect stack count: %d", count)
		}

		if stack := GetHistoryOfDependencyStack(&ref); stack == nil {
			t.Errorf("empty stack, why?")
		} else {
			checks := map[int]string{
				0: "*dij.TestApp",
				1: "*dij.TestLib3",
				3: "*dij.TestLib1",
				5: "*dij.TestLib2",
			}
			for k, v := range checks {
				if name := stack.GetRecord(k).NameOfInstType(); name != v {
					t.Errorf("incorrect stack record: %v, type should be: %s != %s", stack.GetRecord(k), v, name)
				}
			}
		}
	})
}
