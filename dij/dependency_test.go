// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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

type TestBase struct {
	lib1 *TestLib1 `di:"Lib1"`
	lib2 *TestLib2 `di:"Lib2"`
}

type TestExt1 struct {
	TestBase
}

type TestExt2 struct {
	TestBase
}

type TestComb struct {
	ext1 *TestExt1 `di:""`
	ext2 *TestExt2 `di:""`
}

// go test ./dij -v -run TestDI
func TestDI(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		appTyp := reflect.TypeOf(TestApp{})
		ref := DependencyReference{}
		config := map[string]any{
			"ip":   "192.168.0.1",
			"port": 3345,
		}
		ref["config"] = config
		//EnableLog() // uncomment for debug

		inst, err := CreateInstance(appTyp, &ref, "^")
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := inst.(*TestApp); ok {
			//fmt.Println("Inst type:", reflect.TypeOf(inst))
		} else {
			t.Fatal("didn't create a correct instance, ", reflect.TypeOf(inst))
		}

		if count := ref.StackCount(); count != 0 {
			t.Errorf("incorrect stack count: %d", count)
		}

		if stack := ref.StackHistory(); stack == nil {
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

	t.Run("share", func(t *testing.T) {
		comboType := reflect.TypeOf(TestComb{})
		ref := DependencyReference{}
		config := map[string]any{}
		ref["config"] = config
		//EnableLog() // uncomment for debug
		inst, err := CreateInstance(comboType, &ref, "^")
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := inst.(*TestComb); ok {
			//fmt.Println("Inst type:", reflect.TypeOf(inst))
		} else {
			t.Fatal("didn't create a correct instance, ", reflect.TypeOf(inst))
		}

		if count := ref.StackCount(); count != 0 {
			t.Errorf("incorrect stack count: %d", count)
		}

		if stack := ref.StackHistory(); stack == nil {
			t.Errorf("empty stack, why?")
		} else {
			checks := map[int]string{
				0: "*dij.TestComb",
				1: "*dij.TestExt1",
				2: "*dij.TestLib1",
				4: "*dij.TestLib2",
				6: "*dij.TestExt2",
			}
			for k, v := range checks {
				if name := stack.GetRecord(k).NameOfInstType(); name != v {
					t.Errorf("incorrect stack record: %v, type should be: %s != %s", stack.GetRecord(k), v, name)
				}
			}
		}
	})
}

type SampleApp struct {
	lib1    *SampleLib1 `di:"lib1"`
	lib2    *SampleLib2 `di:"lib2"`
	testInt int
}

type SampleLib1 struct {
	lib2 *SampleLib2 `di:"lib2"`
}

type SampleLib2 struct {
	val int `di:"val"`
}

// go test ./dij -v -run TestSample
func TestSample(t *testing.T) {
	t.Run("sample", func(t *testing.T) {
		appTyp := reflect.TypeOf(SampleApp{})
		ref := DependencyReference{"val": 123}
		inst, err := CreateInstance(appTyp, &ref, "^")
		if err != nil {
			t.Fatal(err)
		}
		if app, ok := inst.(*SampleApp); ok {
			if app.lib2 != app.lib1.lib2 {
				t.Errorf("incorrect injection, app.lib2(%v) != app.lib1.lib2(%v)\n", app.lib2, app.lib1.lib2)
			}
			if app.lib2.val != 123 {
				t.Errorf("incorrect injection, app.lib2.val(%d) != 123\n", app.lib2.val)
			}
		} else {
			t.Fatal("didn't create a correct instance, ", reflect.TypeOf(inst))
		}
	})
}

type DijParserSample struct {
	TestBase
	untitled  string
	empty     int            `di:""`
	underline []int          `di:"_"`
	upper     map[int]string `di:"^"`
	disabled  any            `di:"-"`
}

// go test ./dij -v -run TestDijParser
func TestDijParser(t *testing.T) {
	t.Run("fields", func(t *testing.T) {
		list := map[int]struct {
			name    string
			enabled bool
		}{
			0: {"", false},
			1: {"", false},
			2: {"empty", true},
			3: {"underline", true},
			4: {"/map[int]string", true},
			5: {"-", false},
		}
		typ := reflect.TypeOf(DijParserSample{})
		//pri := reflect.TypeOf(map[int]string{1: "234"})
		//t.Logf("====== %s, %v\n", FullnameOfType(pri), pri)
		for n, result := range list {
			if diTag, err := ParseDiTag(typ, n); err != nil {
				t.Errorf("%d [%v] error: %v\n", n, typ.Field(n), err)
			} else {
				if diTag.name != result.name || diTag.Enabled != result.enabled {
					t.Errorf("%d [%v] error, diTAg: %v, correct result: %v\n", n, typ.Field(n), diTag, result)
				} else {
					//t.Logf("%d [%v] correct, diTag: %v\n", n, typ.Field(n), diTag)
				}
			}
		}
	})
}

// go test ./dij -v -run TestDijBuild
func TestDijBuild(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		ref := DependencyReference{"val": 123}
		//var ptr *SampleApp
		ptr := &SampleApp{}
		ptr.testInt = 199
		t.Logf("ptr value: %p\n", ptr)
		ptr2, err := BuildInstance(ptr, &ref, "^")
		if err != nil {
			t.Error(err)
		}
		t.Logf("ptr, ptr2 value: %p, %p\n", ptr, ptr2)
		if ptr2.lib2 != ptr2.lib1.lib2 {
			t.Errorf("incorrect injection, app.lib2(%v) != app.lib1.lib2(%v)\n", ptr2.lib2, ptr2.lib1.lib2)
		}
		if ptr2.lib2.val != 123 {
			t.Errorf("incorrect injection, app.lib2.val(%d) != 123\n", ptr2.lib2.val)
		}
		if ptr2.testInt != 199 {
			t.Errorf("testInt value changed: %d\n", ptr2.testInt)
		}
	})

	t.Run("zero", func(t *testing.T) {
		ref := DependencyReference{"val": 123}
		var ptr *SampleApp
		t.Logf("ptr value: %p\n", ptr)
		ptr2, err := BuildInstance(ptr, &ref, "^")
		if err != nil {
			t.Error(err)
		}
		t.Logf("ptr, ptr2 value: %p, %p\n", ptr, ptr2)
		if ptr2.lib2 != ptr2.lib1.lib2 {
			t.Errorf("incorrect injection, app.lib2(%v) != app.lib1.lib2(%v)\n", ptr2.lib2, ptr2.lib1.lib2)
		}
		if ptr2.lib2.val != 123 {
			t.Errorf("incorrect injection, app.lib2.val(%d) != 123\n", ptr2.lib2.val)
		}
	})

}
