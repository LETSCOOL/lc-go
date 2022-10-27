package di

import (
	"fmt"
	"reflect"
	"testing"
)

type TestI struct {
	v int
}

type TestI2 struct {
	TestI
}

func (t *TestI2) DeclareInjection() map[DependencyKey]any {
	return map[DependencyKey]any{}
}

func (t *TestI2) InjectDependency(p struct {
	str1 string
	_    *string `di:"str2"`
}) {
	fmt.Println("Run initialization", reflect.TypeOf(p))
}

// go test lc-go/di -run TestDependency
func TestDependency(t *testing.T) {
	// root := MakeRootDependency()
	var inf interface{}
	var a = 42
	var x *TestI
	//var x2 *TestI2
	println(inf == nil)
	println(reflect.TypeOf(inf))
	println(reflect.TypeOf(a))
	inf = 42
	println(reflect.TypeOf(inf) == reflect.TypeOf(a))
	inf = a
	println(reflect.TypeOf(inf) == reflect.TypeOf(a))
	inf = x
	println(reflect.TypeOf(x).Name())
	println(reflect.TypeOf(inf).Name())
	println(inf == nil)
	println(x == nil)
}

// go test lc-go/di -v -run TestInitialization
func TestInitialization(t *testing.T) {
	typ := reflect.TypeOf(TestI2{})
	str2 := "12345"
	ref := map[DependencyKey]any{
		"str1": "1234",
		"str2": &str2,
	}
	fmt.Println("type of type", reflect.TypeOf(typ) == reflect.TypeOf(reflect.TypeOf(struct{}{})))
	dep, err := ToInstance(typ)
	if err != nil {
		t.Fatal(err)
	}
	err = InitialDependency(typ, dep, &ref)
	if err != nil {
		t.Fatal(err)
	}
}

type TestApp struct {
}

func (a *TestApp) DeclareInjection() map[DependencyKey]any {
	//config := map[string]any{
	//	"ip":   "192.168.0.1",
	//	"port": 3345,
	//}
	return map[DependencyKey]any{
		//"config": config,
		//"Lib1": reflect.TypeOf(TestLib1{}),
		//"Lib2": reflect.TypeOf(&TestLib2{}),
	}
}

func (a *TestApp) InjectDependency(p struct {
	cfg  map[string]any `di:"config"`
	lib1 *TestLib1      `di:"Lib1"`
	lib2 *TestLib2      `di:"Lib2"`
}) {
	fmt.Printf("TestApp\n\tcfg: %v\n\tlib1: %p %v\n\tlib2: %p %v\n", p.cfg, p.lib1, reflect.TypeOf(p.lib1), p.lib2, reflect.TypeOf(p.lib2))
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

// go test lc-go/di -v -run TestDI
func TestDI(t *testing.T) {
	appTyp := reflect.TypeOf(TestApp{})
	ref := map[DependencyKey]any{}
	config := map[string]any{
		"ip":   "192.168.0.1",
		"port": 3345,
	}
	ref["config"] = config

	inst, err := CreateDependency(appTyp, &ref)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Inst type:", reflect.TypeOf(inst))
}
